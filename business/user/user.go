package user

import (
	"context"

	"github.com/frisk038/swipe_dungeon/business/models"
)

type store interface {
	InsertUser(ctx context.Context, user models.User) (int64, error)
	SelectUserID(ctx context.Context, playerID string) (int64, error)
	UpdateUserType(ctx context.Context, user models.User) error
}

type UserBusiness struct {
	store store
}

func New(store store) *UserBusiness {
	return &UserBusiness{store: store}
}

func (ub *UserBusiness) StoreUser(ctx context.Context, user models.User) (int64, error) {
	// todo valide type
	return ub.store.InsertUser(ctx, user)
}

func (ub *UserBusiness) GetUserID(ctx context.Context, playerID string) (int64, error) {
	return ub.store.SelectUserID(ctx, playerID)
}

func (ub *UserBusiness) UpdateUserType(ctx context.Context, user models.User) error {
	return ub.store.UpdateUserType(ctx, user)
}
