package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"

	"yorgut/internal/controller"
	"yorgut/internal/database"
	"yorgut/internal/helpers"
	"yorgut/internal/repository"
	"yorgut/internal/service"
	"yorgut/internal/storage"
	"yorgut/internal/transport"
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

	database.Migrate(db)

	store, err := storage.NewFromEnv(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("could not set up the uploads storage")
	}
	log.Info().Msgf("Storing uploads in %v", store)
	transport.SetMediaURL(store.URL)

	repositories := repository.NewRepositories(db)
	services := service.NewServices(repositories, store)

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
		*services.ProfileVisit,
	)
	userController := controller.NewUserController(
		*services.User,
		*services.Friendship,
		*services.CommunityMembership,
		*services.ProfileVisit,
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
	pollController := controller.NewPollController(*services.Poll)
	infoController := controller.NewInfoController(*services.Info)

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
		pollController,
		infoController,
		services.Auth,
		store,
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
