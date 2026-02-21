package main

import (
	"log"
	"todolist/internal/config"
	"todolist/internal/user/delivery/http"
	"todolist/internal/user/delivery/http/route"
	"todolist/internal/user/repository/postgres"
	"todolist/internal/user/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	db := config.NewDB()
	defer db.Close()
	rd := config.NewRedisClient()

	userRepo := postgres.NewUserRepository(db, rd)
	userUse := usecase.NewUserUsecase(userRepo)
	userHandle := http.NewUserHandler(userUse)

	r := gin.Default()
	route.SetupUserRoutes(r, userHandle)
	port := config.GetEnv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	// Pastikan tidak ada karakter tersembunyi
	address := ":" + port
	log.Printf("Aplikasi berjalan di %s", address)
	r.Run(address)
}
