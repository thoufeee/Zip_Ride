package main

import (
	"log"
	"zipride/database"
	"zipride/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// env load
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load env")
	}

	// database connect
	database.Connect()

	// redis connect
	database.InitRedis()

	r := gin.Default()

	routes.User(r)

	r.Run(":8080")
}
