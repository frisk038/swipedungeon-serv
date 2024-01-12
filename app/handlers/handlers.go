package handlers

import (
	"context"
	"net/http"

	"github.com/frisk038/swipe_dungeon/business/models"
	"github.com/gin-gonic/gin"
)

type UserManager interface {
	StoreUser(ctx context.Context, user models.User) (int64, error)
	GetUserID(ctx context.Context, playerID string) (int64, error)
	UpdateUserType(ctx context.Context, user models.User) error
}

func PostUser(um UserManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user = struct {
			Name     string `json:"name" binding:"required"`
			PlayerID string `json:"player_id" binding:"required"`
		}{}
		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		id, err := um.StoreUser(c.Request.Context(), models.User{
			Name:     user.Name,
			PlayerID: user.PlayerID,
		})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": id})
	}
}

func GetUser(um UserManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		player_id := c.Param("player_id")
		user_id, err := um.GetUserID(c.Request.Context(), player_id)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": user_id})
	}
}

func UpdateUserType(um UserManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user = struct {
			UserID    int64  `json:"user_id" binding:"required"`
			PowerType string `json:"power_type" binding:"required"`
		}{}
		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		err = um.UpdateUserType(c.Request.Context(), models.User{
			UserID:    user.UserID,
			PowerType: models.PowerType(user.PowerType),
		})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusOK)
	}
}
