package models

import "github.com/google/uuid"

type User struct {
	UserID    uuid.UUID
	Name      string
	PlayerID  string
	PowerType PowerType
	CharaID   int64
	Loc       Location
}
