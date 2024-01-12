package main

import (
	"log"
	"os"

	"github.com/frisk038/swipe_dungeon/app/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func initRoutes() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r := gin.Default()
	r.Use(cors.Default())
	geo := r.Group("/", gin.BasicAuth(gin.Accounts{
		os.Getenv("USER"): os.Getenv("PASS"),
	}))

	geo.POST("/user", handlers.PostUser())

	r.Run(":" + port)
}

func main() {
	initRoutes()
}
