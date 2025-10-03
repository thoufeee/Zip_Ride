package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"time"
	"zipride/database"
	"zipride/routes"
	"zipride/seeders"
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

	// seeder admin
	seeders.SeedAdmin()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500", "http://localhost:5500"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// routes connection
	routes.PublicRoutes(r)
	routes.UserRoutes(r)
	routes.DriverRoutes(r)

	r.Run(":8080")
}
