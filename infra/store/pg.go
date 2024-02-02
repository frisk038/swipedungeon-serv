package store

import (
	"context"
	"os"

	"github.com/frisk038/swipe_dungeon/business/models"
	"github.com/google/uuid"

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
	insertUserLocation = "INSERT INTO userlocation(user_id, coord) values($1, ST_MakePoint($2, $3));"
	selectNearbyUser   = `SELECT DISTINCT ON (name) name, power_type, chara_id
							FROM userlocation
							LEFT JOIN users ON users.user_id = userlocation.user_id
							WHERE ST_DWithin(coord, ST_MakePoint($1, $2)::geography, $3)
							AND users.user_id != $4 
							ORDER BY name, created_at
							LIMIT $5;`
	insertUserScore = `INSERT INTO userscore (user_id, score, user_level)
							VALUES ($1, $2, $3)
							ON CONFLICT (user_id) DO UPDATE
							SET score = GREATEST(EXCLUDED.score, userscore.score),
							user_level = EXCLUDED.user_level;`

	MaxNearbyUserLimit    = 5
	MaxNearbyUserDistance = 10000
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

func (c *Client) InsertUserLocation(ctx context.Context, user_id uuid.UUID, coord models.Coordinate) error {
	_, err := c.conn.Exec(ctx, insertUserLocation, user_id, coord.Latitude, coord.Longitude)
	return err
}

func (c *Client) SelectNearbyUser(ctx context.Context, user_id uuid.UUID, coord models.Coordinate) ([]models.User, error) {
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
		err = rows.Scan(&user.Name, &ptype, &charaID)
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
