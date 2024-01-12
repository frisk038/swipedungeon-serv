package store

import (
	"context"
	"os"

	"github.com/frisk038/swipe_dungeon/business/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	conn *pgxpool.Pool
}

const (
	insertUser     = "INSERT INTO users(name, player_id) VALUES($1, $2) RETURNING user_id;"
	selectUserID   = "SELECT user_id from users WHERE player_id=$1;"
	updateUserType = "UPDATE users SET power_type=$1 WHERE user_id=$2;"
)

func New() (*Client, error) {
	conn, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

func (c *Client) InsertUser(ctx context.Context, user models.User) (int64, error) {
	row := c.conn.QueryRow(ctx, insertUser, user.Name, user.PlayerID)
	var id int64
	if err := row.Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}

func (c *Client) SelectUserID(ctx context.Context, playerID string) (int64, error) {
	row := c.conn.QueryRow(ctx, selectUserID, playerID)
	var id int64
	if err := row.Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}

func (c *Client) UpdateUserType(ctx context.Context, user models.User) error {
	_, err := c.conn.Exec(ctx, updateUserType, user.PowerType, user.UserID)
	return err
}
