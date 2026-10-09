package server

import (
	"net/http"
	"yorgut/internal/controller"
	"yorgut/internal/helpers"
	"yorgut/internal/service"

	_ "yorgut/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	authService *service.AuthService,
) *gin.Engine {
	router := gin.Default()

	router.Static("/storage/uploads", helpers.GetUploadsPath(true))

	// router.Use(DelayMiddleware()) // Used for testing

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://yorgut.onrender.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/api", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Yorgut")
	})

	// Public: the only way in, so it can't require a session.
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/google", authController.LoginWithGoogle)
	}

	v1 := router.Group("/api/v1", AuthMiddleware(authService))
	{
		me := v1.Group("/me")
		{
			me.GET("", meController.FindMe)
			me.PUT("", meController.UpdateMe)
			me.PUT("/profile-picture", meController.UpdateMyProfilePicture)
			me.GET("/welcome", meController.Welcome)
			me.GET("/profile", meController.FindMyProfile)
			me.GET("/my-communities", meController.FindMyCommunities)
			me.GET("/memberships", meController.FindMyMemberships)
			me.GET("/memberships/:communityId", meController.AmIMemberOf)
			me.GET("/friends", meController.FindMyFriends)
			me.GET("/friends/pending", meController.FindMyPendingFriendships)
			me.GET("/friends/:userId", meController.AmIFriendOf)
			me.POST("/join/:communityId", meController.JoinCommunity)
			me.POST("/leave/:communityId", meController.LeaveCommunity)
		}

		users := v1.Group("/users")
		{
			users.POST("", userController.Create)
			users.GET("", userController.FindMany)
			users.GET("/search-users", userController.SearchUsers)
			users.GET("/:userId", userController.FindById)
			users.GET("/:userId/profile", userController.FindUserProfile)
			users.GET("/:userId/memberships", userController.FindMemberships)
			users.GET("/:userId/friends", userController.FindFriends)
			users.GET("/:userId/reviews", userReviewController.FindSummary)
			users.PUT("/:userId/review", userReviewController.Upsert)
			users.GET("/:userId/fans", userFanController.FindSummary)
			users.PUT("/:userId/fan", userFanController.BecomeFan)
			users.DELETE("/:userId/fan", userFanController.StopBeingFan)
		}

		communities := v1.Group("/communities")
		{
			communities.POST("", communityController.Create)
			communities.PUT("/:communityId", communityController.Update)
			communities.PUT("/:communityId/picture", communityController.UpdatePicture)
			communities.GET("", communityController.FindMany)
			communities.GET("/search", communityController.Search)
			communities.DELETE("/:communityId", communityController.Delete)
			communities.GET("/:communityId", communityController.FindById)
			communities.GET("/:communityId/members", communityController.FindMembers)
			communities.GET("/:communityId/profile", communityController.FindProfileById)
			communities.GET("/:communityId/related", communityController.FindRelated)
			communities.GET("/:communityId/topics", communityController.FindAllTopics)
			communities.GET("/:communityId/hot-topics", communityController.FindHotTopics)
		}

		communityMemberships := v1.Group("/community-memberships")
		{
			communityMemberships.POST("", communityMembershipController.Create)
			communityMemberships.DELETE("/:userId/:communityId", communityMembershipController.Delete)
		}

		friendships := v1.Group("/friendships")
		{
			friendships.POST("", friendshipController.AddFriend)
			friendships.POST("/reply", friendshipController.ReplyFriendship)
			friendships.DELETE("/:friendshipId", friendshipController.RemoveFriendByFriendshipID)
			friendships.DELETE("/friend/:userId", friendshipController.RemoveFriendByUserID)
		}

		photos := v1.Group("/photos")
		{
			photos.GET("/user/:userId", photoController.FindOfUser)
			photos.POST("", photoController.Create)
			photos.DELETE("/:photoId", photoController.Delete)
		}

		scraps := v1.Group("/scraps")
		{
			scraps.GET("/user/:userId", scrapController.FindOfUser)
			scraps.POST("", scrapController.Create)
			scraps.DELETE("", scrapController.Delete)
		}

		testimonials := v1.Group("/testimonials")
		{
			testimonials.GET("/user/:userId", testimonialController.FindOfUser)
			testimonials.POST("", testimonialController.Create)
			testimonials.POST("/:testimonialId/approve", testimonialController.Approve)
			testimonials.DELETE("/:testimonialId", testimonialController.Delete)
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
