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
	srvAddr := fmt.Sprintf(":%d", port)

	db := database.DatabaseConnection()
	validate := validator.New()

	database.Migrate(db)

	// Repository
	userRepository := repository.NewUserRepository(db)
	communityMembershipRepository := repository.NewCommunityMembershipRepository(db)
	communityRepository := repository.NewCommunityRepository(db)
	friendshipRepository := repository.NewFriendshipRepository(db)
	photoAlbumRepository := repository.NewPhotoAlbumRepository(db)
	photoRepository := repository.NewPhotoRepository(db)
	scrapRepository := repository.NewScrapRepository(db)
	testimonialRepository := repository.NewTestimonialRepository(db)

	// Service
	userService := service.NewUserService(*userRepository, validate)
	communityMembershipService := service.NewCommunityMembershipService(*communityMembershipRepository, validate)
	communityService := service.NewCommunityService(*communityRepository, validate)
	friendshipService := service.NewFriendshipService(*friendshipRepository, validate)
	photoAlbumService := service.NewPhotoAlbumService(*photoAlbumRepository, validate)
	photoService := service.NewPhotoService(*photoRepository, validate)
	scrapService := service.NewScrapService(*scrapRepository, validate)
	testimonialService := service.NewTestimonialService(*testimonialRepository, validate)

	// Controllers
	userController := controller.NewUserController(*userService)
	communityMembershipController := controller.NewCommunityMembershipController(*communityMembershipService)
	communityController := controller.NewCommunityController(*communityService)
	friendshipController := controller.NewFriendshipController(*friendshipService)
	photoAlbumController := controller.NewPhotoAlbumController(*photoAlbumService)
	photoController := controller.NewPhotoController(*photoService)
	scrapController := controller.NewScrapController(*scrapService)
	testimonialController := controller.NewTestimonialController(*testimonialService)

	// Router
	routes := NewRouter(
		userController,
		communityController,
		communityMembershipController,
		friendshipController,
		photoAlbumController,
		photoController,
		scrapController,
		testimonialController,
	)

	// Declare Server config
	server := &http.Server{
		Addr:         srvAddr,
		Handler:      routes,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Info().Msg(fmt.Sprintf("Server started at: %v", srvAddr))

	return server
}
