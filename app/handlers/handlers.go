package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user = struct {
			UserID    int64  `json:"user_id"`
			Name      string `json:"name"`
			PlayerID  string `json:"player_id"`
			PowerType string `json:"power_type"`
		}{}
		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.Status(http.StatusOK)
	}
}
