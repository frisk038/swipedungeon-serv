package user

import (
	"context"

	"github.com/frisk038/swipe_dungeon/business/models"
	"github.com/google/uuid"
)

type store interface {
	InsertUser(ctx context.Context, user models.User) (uuid.UUID, error)
	SelectUserID(ctx context.Context, playerID string) (uuid.UUID, error)
	UpdateUserInfo(ctx context.Context, user models.User) error
	SelectNearbyUser(ctx context.Context, user_id uuid.UUID, coord models.Coordinate) ([]models.User, error)
	InsertUserLocation(ctx context.Context, user_id uuid.UUID, coord models.Coordinate) error
}

type UserBusiness struct {
	store store
}

func New(store store) *UserBusiness {
	return &UserBusiness{store: store}
}

func (ub *UserBusiness) StoreUser(ctx context.Context, user models.User) (uuid.UUID, error) {
	return ub.store.InsertUser(ctx, user)
}

func (ub *UserBusiness) GetUserID(ctx context.Context, playerID string) (uuid.UUID, error) {
	return ub.store.SelectUserID(ctx, playerID)
}

func (ub *UserBusiness) UpdateUserInfo(ctx context.Context, user models.User) error {
	// todo valide type
	return ub.store.UpdateUserInfo(ctx, user)
}

func (ub *UserBusiness) GetNearbyUser(ctx context.Context, user_id uuid.UUID, coord models.Coordinate) ([]models.User, error) {
	// todo valide type
	return ub.store.SelectNearbyUser(ctx, user_id, coord)
}

func (ub *UserBusiness) StoreUserLocation(ctx context.Context, user_id uuid.UUID, coord models.Coordinate) error {
	// todo valide type
	return ub.store.InsertUserLocation(ctx, user_id, coord)
}
