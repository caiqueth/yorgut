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
	"yorgut/internal/helpers"
	"yorgut/internal/repository"
	"yorgut/internal/service"
)

func NewServer() *http.Server {
	log.Info().Msg("Starting the server")
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	srvAddr := fmt.Sprintf("0.0.0.0:%d", port)

	if !helpers.SessionSecretConfigured() {
		log.Fatal().Msg("SESSION_SECRET is not set; it signs the login session tokens")
	}
	// Without it, Google ID tokens issued to any app would be accepted.
	if os.Getenv("GOOGLE_CLIENT_ID") == "" {
		log.Fatal().Msg("GOOGLE_CLIENT_ID is not set; it is needed to validate Google logins")
	}

	db := database.DatabaseConnection()
	validate := validator.New()

	database.Migrate(db)

	repositories := repository.NewRepositories(db)
	services := service.NewServices(repositories, validate)

	// Controllers
	authController := controller.NewAuthController(*services.Auth)
	meController := controller.NewMeController(
		*services.User,
		*services.Friendship,
		*services.Community,
		*services.CommunityMembership,
		*services.UserFan,
		*services.Scrap,
		*services.UserProfile,
	)
	userController := controller.NewUserController(
		*services.User,
		*services.Friendship,
		*services.CommunityMembership,
	)
	communityMembershipController := controller.NewCommunityMembershipController(*services.CommunityMembership)
	communityController := controller.NewCommunityController(*services.Community, *services.Topic)
	friendshipController := controller.NewFriendshipController(*services.Friendship)
	photoController := controller.NewPhotoController(*services.Photo)
	scrapController := controller.NewScrapController(*services.Scrap)
	testimonialController := controller.NewTestimonialController(*services.Testimonial)
	topicController := controller.NewTopicController(*services.Topic)
	userReviewController := controller.NewUserReviewController(*services.UserReview)
	userFanController := controller.NewUserFanController(*services.UserFan)

	// Router
	routes := NewRouter(
		authController,
		meController,
		userController,
		communityController,
		communityMembershipController,
		friendshipController,
		photoController,
		scrapController,
		testimonialController,
		topicController,
		userReviewController,
		userFanController,
		services.Auth,
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
