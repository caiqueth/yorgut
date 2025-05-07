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
// @in header
// @name token// @license.name Apache 2.0
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
func NewRouter(
	meController *controller.MeController,
	userController *controller.UserController,
	communityController *controller.CommunityController,
	communityMembershipController *controller.CommunityMembershipController,
	friendshipController *controller.FriendshipController,
	photoAlbumController *controller.PhotoAlbumController,
	photoController *controller.PhotoController,
	scrapController *controller.ScrapController,
	testimonialController *controller.TestimonialController,
	userService *service.UserService,
) *gin.Engine {
	router := gin.Default()

	router.Static("/storage/uploads", helpers.GetUploadsPath(true))

	userInjector := InjectUserMiddleware(userService)
	jwtInjector := JWTMiddleware()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/api", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Yorgut")
	})

	v1 := router.Group("/api/v1", jwtInjector, userInjector)
	{
		me := v1.Group("/me")
		{
			me.GET("", meController.FindMe)
			me.PUT("", meController.UpdateMe)
			me.GET("/profile", meController.FindMyProfile)
			me.GET("/memberships", meController.FindMyMemberships)
			me.GET("/friends", meController.FindMyFriends)
			me.GET("/friends/pending", meController.FindMyPendingFriendships)
			me.GET("/photoAlbums", meController.FindMyPhotoAlbums)
			me.POST("/join/:communityId", meController.JoinCommunity)
			me.POST("/leave/:communityId", meController.LeaveCommunity)
		}

		users := v1.Group("/users")
		{
			users.POST("", userController.Create)
			users.GET("", userController.FindMany)
			users.GET("/:userId", userController.FindById)
			users.GET("/:userId/profile", userController.FindUserProfile)
			users.GET("/:userId/memberships", userController.FindMemberships)
			users.GET("/:userId/friends", userController.FindFriends)
			users.GET("/:userId/photoAlbums", userController.FindPhotoAlbums)
		}

		communities := v1.Group("/community")
		{
			communities.POST("", communityController.Create)
			communities.PUT("", communityController.Update)
			communities.GET("", communityController.FindMany)
			communities.DELETE("/:communityId", communityController.Delete)
			communities.GET("/:communityId", communityController.FindById)
			communities.GET("/:communityId/members", communityController.FindMembers)
			communities.GET("/:communityId/profile", communityController.FindProfileById)
		}

		community_memberships := v1.Group("/community_membership")
		{
			community_memberships.POST("", communityMembershipController.Create)
			community_memberships.DELETE("/:userId/:communityId", communityMembershipController.Delete)
		}

		friendships := v1.Group("/friendship")
		{
			friendships.POST("", friendshipController.Create)
			friendships.DELETE("", friendshipController.Delete)
		}

		photoAlbums := v1.Group("/photoAlbum")
		{
			photoAlbums.POST("", photoAlbumController.Create)
			photoAlbums.DELETE("", photoAlbumController.Delete)
		}

		photos := v1.Group("/photo")
		{
			photos.POST("", photoController.Create)
			photos.DELETE("", photoController.Delete)
		}

		scraps := v1.Group("/scrap")
		{
			scraps.GET("/user/:userId", scrapController.FindMany)
			scraps.POST("", scrapController.Create)
			scraps.DELETE("", scrapController.Delete)
		}

		testimonials := v1.Group("/testimonial")
		{
			testimonials.POST("", testimonialController.Create)
			testimonials.DELETE("", testimonialController.Delete)
		}

	}

	return router
}
