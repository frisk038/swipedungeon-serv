package handlers

import (
	"context"
	"net/http"

	"github.com/frisk038/swipe_dungeon/business/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserManager interface {
	StoreUser(ctx context.Context, user models.User) (uuid.UUID, error)
	GetUserID(ctx context.Context, playerID string) (uuid.UUID, error)
	UpdateUserInfo(ctx context.Context, user models.User) error
	GetNearbyUser(ctx context.Context, user_id uuid.UUID, coord models.Coordinate) ([]models.User, error)
	StoreUserLocation(ctx context.Context, user_id uuid.UUID, coord models.Coordinate) error
	StoreUserScore(ctx context.Context, user_id uuid.UUID, score models.Score) error
}

type nearByResp struct {
	Name      string `json:"name"`
	PowerType string `json:"power_type"`
	CharaID   int64  `json:"char_id"`
}

func PostUser(um UserManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user = struct {
			Name      string `json:"name" binding:"required"`
			PlayerID  string `json:"player_id" binding:"required"`
			CharaID   int64  `json:"chara_id"`
			PowerType string `json:"power_type" binding:"required"`
		}{}
		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		id, err := um.StoreUser(c.Request.Context(), models.User{
			Name:      user.Name,
			PlayerID:  user.PlayerID,
			PowerType: models.PowerType(user.PowerType),
			CharaID:   user.CharaID,
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

func UpdateUserInfo(um UserManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user = struct {
			UserID    uuid.UUID `json:"user_id" binding:"required"`
			PowerType string    `json:"power_type" binding:"required"`
			CharaID   int64     `json:"chara_id"`
		}{}
		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		err = um.UpdateUserInfo(c.Request.Context(), models.User{
			UserID:    user.UserID,
			PowerType: models.PowerType(user.PowerType),
			CharaID:   user.CharaID,
		})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func GetNearbyUser(um UserManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req = struct {
			UserID uuid.UUID `json:"user_id" binding:"required"`
			Coord  struct {
				Latitude  string `json:"latitude" binding:"required"`
				Longitude string `json:"longitude" binding:"required"`
			} `json:"coordinates" binding:"required"`
		}{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		users, err := um.GetNearbyUser(c.Request.Context(), req.UserID, models.Coordinate{
			Longitude: req.Coord.Longitude,
			Latitude:  req.Coord.Latitude,
		})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		userResp := make([]nearByResp, 0)
		for _, u := range users {
			userResp = append(userResp, nearByResp{
				Name:      u.Name,
				PowerType: string(u.PowerType),
				CharaID:   u.CharaID,
			})
		}
		c.JSON(http.StatusOK, gin.H{"users": userResp})
	}
}

func POSTUserCoordinates(um UserManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req = struct {
			UserID uuid.UUID `json:"user_id" binding:"required"`
			Coord  struct {
				Latitude  string `json:"latitude" binding:"required"`
				Longitude string `json:"longitude" binding:"required"`
			} `json:"coordinates" binding:"required"`
		}{}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		err := um.StoreUserLocation(c.Request.Context(), req.UserID, models.Coordinate{
			Longitude: req.Coord.Longitude,
			Latitude:  req.Coord.Latitude,
		})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func POSTUserScore(um UserManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req = struct {
			Floor int `json:"floor" binding:"required"`
			Level int `json:"level" binding:"required"`
		}{}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		userIDstr := c.Param("user_id")
		userID, err := uuid.Parse(userIDstr)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if err = um.StoreUserScore(c.Request.Context(), userID, models.Score{
			Level: req.Level,
			Floor: req.Floor,
		}); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusOK)
	}
}
