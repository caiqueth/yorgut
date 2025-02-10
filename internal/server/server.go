package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"

	"yorgut/internal/controller"
	"yorgut/internal/database"
	"yorgut/internal/repository"
	"yorgut/internal/service"
)

func NewServer() *http.Server {
	log.Info().Msg("Starting the server")
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	db := database.DatabaseConnection()
	validate := validator.New()

	database.Migrate(db)

	// Repository
	userRepository := repository.NewUserRepositoryImpl(db)

	// Service
	userService := service.NewUserServiceImpl(userRepository, validate)

	// Controller
	userController := controller.NewUserController(userService)

	// Router
	routes := NewRouter(userController)

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      routes,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
