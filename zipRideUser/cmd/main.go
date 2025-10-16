package main

import (
	"log"
	"time"
	"zipride/database"
	"zipride/routes"
	"zipride/seeders"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// env load
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("failed to load env")
	}

	// database connect
	database.Connect()

	// redis connect
	database.InitRedis()

	// seeders run
	seeders.RunAllSeeders()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500", "http://localhost:5500"},
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// routes connection
	routes.PublicRoutes(r)
	routes.UserRoutes(r)
	routes.SuperAdminRoutes(r)

	// port
	r.Run(":8080")
}
