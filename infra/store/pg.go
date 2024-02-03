package store

import (
	"context"
	"fmt"
	"os"

	"github.com/frisk038/swipe_dungeon/business/models"
	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	conn *pgxpool.Pool
}

const (
	insertUser = `WITH e AS(
		            INSERT INTO users(name, player_id, power_type, chara_id) 
		            values($1, $2, $3, $4) 
		            ON CONFLICT DO NOTHING
		            RETURNING user_id
	            ) SELECT * FROM e 
				UNION SELECT user_id FROM users WHERE name=$1;`
	selectUserID       = "SELECT user_id from users WHERE name=$1;"
	updateUserType     = "UPDATE users SET power_type=$1, chara_id=$2 WHERE user_id=$3;"
	insertUserLocation = "INSERT INTO userlocation(user_id, coord, city) values($1, ST_MakePoint($2, $3), $4);"
	selectNearbyUser   = `SELECT DISTINCT ON (name) name, power_type, chara_id, city, users.user_id
							FROM userlocation
							LEFT JOIN users ON users.user_id = userlocation.user_id
							WHERE ST_DWithin(coord, ST_MakePoint($1, $2)::geography, $3)
							AND users.user_id NOT IN (SELECT seen_user FROM seen WHERE user_id = $4)
							AND users.user_id != $4 
							ORDER BY name, created_at
							LIMIT $5;`
	insertUserScore = `INSERT INTO userscore (user_id, score, user_level)
						VALUES ($1, $2, $3)
						ON CONFLICT (user_id) DO UPDATE
						SET score = GREATEST(EXCLUDED.score, userscore.score),
							user_level = CASE WHEN EXCLUDED.score > userscore.score 
							THEN EXCLUDED.user_level 
							ELSE userscore.user_level END;`
	topPlayer = `SELECT users.name, userscore.score 
					FROM userscore LEFT JOIN users 
					ON userscore.user_id=users.user_id 
					ORDER BY score DESC LIMIT 1`
	insertSeenUser = "INSERT INTO seen(user_id, seen_user) values($1, $2);"

	MaxNearbyUserLimit    = 20
	MaxNearbyUserDistance = 1000
)

func New() (*Client, error) {
	conn, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

func (c *Client) InsertUser(ctx context.Context, user models.User) (uuid.UUID, error) {
	row := c.conn.QueryRow(ctx, insertUser, user.Name, user.PlayerID, user.PowerType, user.CharaID)
	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}

func (c *Client) SelectUserID(ctx context.Context, playerID string) (uuid.UUID, error) {
	row := c.conn.QueryRow(ctx, selectUserID, playerID)
	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}

func (c *Client) UpdateUserInfo(ctx context.Context, user models.User) error {
	_, err := c.conn.Exec(ctx, updateUserType, user.PowerType, user.CharaID, user.UserID)
	return err
}

func (c *Client) InsertUserLocation(ctx context.Context, user_id uuid.UUID, coord models.Location) error {
	_, err := c.conn.Exec(ctx, insertUserLocation, user_id, coord.Latitude, coord.Longitude, coord.City)
	return err
}

func (c *Client) SelectNearbyUser(ctx context.Context, user_id uuid.UUID, coord models.Location) ([]models.User, error) {
	rows, err := c.conn.Query(ctx, selectNearbyUser,
		coord.Latitude, coord.Longitude, MaxNearbyUserDistance, user_id, MaxNearbyUserLimit)
	if err != nil {
		return nil, err
	}
	users := make([]models.User, 0)
	var (
		user        models.User
		ptype       *string
		charaID     *int64
		defaultType = string(models.FIRE)
	)

	for rows.Next() {
		err = rows.Scan(&user.Name, &ptype, &charaID, &user.Loc.City, &user.UserID)
		if err != nil {
			return nil, err
		}
		if ptype == nil {
			ptype = &defaultType
		}
		user.PowerType = models.PowerType(*ptype)

		if charaID != nil {
			user.CharaID = *charaID
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (c *Client) InsertUserScore(ctx context.Context, userID uuid.UUID, score models.Score) error {
	if _, err := c.conn.Exec(ctx, insertUserScore, userID, score.Floor, score.Level); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetLeaderboard(ctx context.Context) (models.LeaderBoard, error) {
	row := c.conn.QueryRow(ctx, topPlayer)
	var lb models.LeaderBoard

	if err := row.Scan(&lb.UserName, &lb.Score); err != nil {
		return models.LeaderBoard{}, err
	}
	return lb, nil
}

func (c *Client) InsertSeenUser(ctx context.Context, userID uuid.UUID, seen []models.User) error {

	rows := [][]any{}
	for _, s := range seen {
		rows = append(rows, []any{userID, s.UserID})
	}

	fmt.Println(rows)

	_, err := c.conn.CopyFrom(
		ctx,
		pgx.Identifier{"seen"},
		[]string{"user_id", "seen_user"},
		pgx.CopyFromRows(rows))

	return err
}
