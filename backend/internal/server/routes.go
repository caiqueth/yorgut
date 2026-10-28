package server

import (
	"net/http"
	"os"
	"yorgut/internal/controller"
	"yorgut/internal/service"
	"yorgut/internal/storage"

	_ "yorgut/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Yorgut APIs
// @version 1.0
// @description Our very own OG Social Media.
// @license.name Apache 2.0
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Yorgut session token from POST /auth/google, sent as "Bearer <token>".
func NewRouter(
	authController *controller.AuthController,
	meController *controller.MeController,
	userController *controller.UserController,
	communityController *controller.CommunityController,
	communityMembershipController *controller.CommunityMembershipController,
	friendshipController *controller.FriendshipController,
	photoController *controller.PhotoController,
	scrapController *controller.ScrapController,
	testimonialController *controller.TestimonialController,
	topicController *controller.TopicController,
	userReviewController *controller.UserReviewController,
	userFanController *controller.UserFanController,
	pollController *controller.PollController,
	infoController *controller.InfoController,
	authService *service.AuthService,
	store storage.ObjectStore,
) *gin.Engine {
	router := gin.Default()
	// Behind Cloudflare (production), the visitor's address is in the
	// CF-Connecting-IP header; nginx's and the tunnel's addresses aren't
	// trusted for anything else. Used by the login rate limit.
	router.TrustedPlatform = gin.PlatformCloudflare
	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatal().Err(err).Msg("could not set trusted proxies")
	}

	// A disk store's pictures are served by the API, at the path of
	// MEDIA_BASE_URL (e.g. /media). Other stores serve their own.
	if disk, ok := store.(*storage.Disk); ok {
		router.Static(disk.Path, disk.Dir)
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// DISABLE_SWAGGER=true (production) leaves the API docs unpublished.
	if os.Getenv("DISABLE_SWAGGER") != "true" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	router.GET("/api", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Yorgut")
	})

	// API_DELAY (development only) slows the API down, to see loading states.
	api := router.Group("")
	if shortest, longest, ok := delayFromEnv(); ok {
		log.Warn().Dur("shortest", shortest).Dur("longest", longest).Msg("API_DELAY is set: every API response is delayed")
		api.Use(DelayMiddleware(shortest, longest))
	}

	// Public: the only way in, so it can't require a session.
	auth := api.Group("/api/v1/auth")
	{
		auth.POST("/google", RateLimit(loginLimit, byClientAddress), authController.LoginWithGoogle)
	}

	// Public: the about page shows it before logging in.
	api.GET("/api/v1/info", infoController.FindInfo)

	v1 := api.Group("/api/v1", AuthMiddleware(authService), RateLimit(writeLimit, writesByUser))
	// Uploads also count against their own, smaller limit.
	uploads := RateLimit(uploadLimit, byUser)
	{
		me := v1.Group("/me")
		{
			me.GET("", meController.FindMe)
			me.PUT("", meController.UpdateMe)
			me.DELETE("", meController.DeleteMe)
			me.PUT("/profile-picture", uploads, meController.UpdateMyProfilePicture)
			me.PUT("/privacy", meController.UpdateMyPrivacy)
			me.PUT("/theme", meController.UpdateMyTheme)
			me.GET("/welcome", meController.Welcome)
			me.GET("/profile", meController.FindMyProfile)
			me.GET("/my-communities", meController.FindMyCommunities)
			me.GET("/memberships", meController.FindMyMemberships)
			me.GET("/memberships/:communityId", meController.AmIMemberOf)
			me.GET("/friends", meController.FindMyFriends)
			me.GET("/friends/pending", meController.FindMyPendingFriendships)
			me.GET("/birthdays", meController.FindMyBirthdays)
			me.GET("/friends/:userId", meController.AmIFriendOf)
			me.POST("/join/:communityId", meController.JoinCommunity)
			me.POST("/leave/:communityId", meController.LeaveCommunity)
			me.DELETE("/join-requests/:communityId", meController.CancelJoinRequest)
			me.GET("/community-invitations", meController.FindMyCommunityInvitations)
			me.POST("/community-invitations/:invitationId/accept", meController.AcceptCommunityInvitation)
			me.DELETE("/community-invitations/:invitationId", meController.DeclineCommunityInvitation)
		}

		users := v1.Group("/users")
		{
			users.GET("/search-users", userController.SearchUsers)
			users.GET("/:userId", userController.FindById)
			users.GET("/:userId/profile", userController.FindUserProfile)
			users.GET("/:userId/memberships", userController.FindMemberships)
			users.GET("/:userId/friends", userController.FindFriends)
			users.GET("/:userId/reviews", userReviewController.FindSummary)
			users.PUT("/:userId/review", userReviewController.Upsert)
			users.GET("/:userId/fans", userFanController.FindSummary)
			users.GET("/:userId/fans/list", userFanController.FindList)
			users.PUT("/:userId/fan", userFanController.BecomeFan)
			users.DELETE("/:userId/fan", userFanController.StopBeingFan)
		}

		communities := v1.Group("/communities")
		{
			communities.POST("", uploads, communityController.Create)
			communities.PUT("/:communityId", communityController.Update)
			communities.PUT("/:communityId/picture", uploads, communityController.UpdatePicture)
			communities.GET("/search", communityController.Search)
			communities.GET("/:communityId", communityController.FindById)
			communities.GET("/:communityId/members", communityController.FindMembers)
			communities.GET("/:communityId/profile", communityController.FindProfileById)
			communities.GET("/:communityId/related", communityController.FindRelated)
			communities.GET("/:communityId/topics", communityController.FindAllTopics)
			communities.GET("/:communityId/hot-topics", communityController.FindHotTopics)
			communities.GET("/:communityId/join-requests", communityMembershipController.FindJoinRequests)
			communities.POST("/:communityId/join-requests/:userId/reply", communityMembershipController.ReplyJoinRequest)
			communities.GET("/:communityId/invitable-friends", communityMembershipController.FindInvitableFriends)
			communities.POST("/:communityId/invitations", communityMembershipController.InviteToCommunity)
			communities.GET("/:communityId/polls", pollController.FindPolls)
			communities.POST("/:communityId/polls", pollController.CreatePoll)
		}

		friendships := v1.Group("/friendships")
		{
			friendships.POST("", friendshipController.AddFriend)
			friendships.POST("/reply", friendshipController.ReplyFriendship)
			friendships.DELETE("/friend/:userId", friendshipController.RemoveFriendByUserID)
		}

		photos := v1.Group("/photos")
		{
			photos.GET("/user/:userId", photoController.FindOfUser)
			photos.POST("", uploads, photoController.Create)
			photos.DELETE("/:photoId", photoController.Delete)
		}

		scraps := v1.Group("/scraps")
		{
			scraps.GET("/user/:userId", scrapController.FindOfUser)
			scraps.POST("", scrapController.Create)
			scraps.DELETE("/:scrapId", scrapController.Delete)
		}

		testimonials := v1.Group("/testimonials")
		{
			testimonials.GET("/user/:userId", testimonialController.FindOfUser)
			testimonials.POST("", testimonialController.Create)
			testimonials.POST("/:testimonialId/approve", testimonialController.Approve)
			testimonials.DELETE("/:testimonialId", testimonialController.Delete)
		}

		polls := v1.Group("/polls")
		{
			polls.GET("/:pollId", pollController.FindPoll)
			polls.POST("/:pollId/vote", pollController.Vote)
			polls.POST("/:pollId/close", pollController.Close)
		}

		topics := v1.Group("/topics")
		{
			topics.GET("/:topicId/replies", topicController.FindTopicAndReplies)
			topics.POST("", topicController.CreateTopic)
			topics.POST("/:topicId", topicController.ReplyTopic)
		}

	}

	return router
}
