package main

import (
	"fmt"
	"log"
	"os"

	"github.com/frisk038/swipe_dungeon/app/handlers"
	"github.com/frisk038/swipe_dungeon/business/user"
	"github.com/frisk038/swipe_dungeon/infra/store"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func initRoutes(um handlers.UserManager) {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r := gin.Default()
	r.Use(cors.Default(), gin.BasicAuth(gin.Accounts{
		os.Getenv("USER"): os.Getenv("PASS"),
	}))

	r.POST("/user", handlers.PostUser(um))
	r.GET("/user/:player_id", handlers.GetUser(um))
	r.PUT("/user", handlers.UpdateUserInfo(um))

	r.POST("/geo", handlers.POSTUserCoordinates(um))
	r.GET("/geo", handlers.GetNearbyUser(um))

	r.Run(":" + port)
}

func main() {
	repo, err := store.New()
	if err != nil {
		fmt.Println(err)
	}

	ub := user.New(repo)
	initRoutes(ub)
}
