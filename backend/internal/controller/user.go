package controller

import (
	"net/http"
	"strconv"
	"yorgut/internal/dto"
	"yorgut/internal/helpers"
	"yorgut/internal/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(service service.UserService) *UserController {
	return &UserController{
		userService: service,
	}
}

// Create			godoc
// @Summary			Create user
// @Decription		Save user in Db.
// @Param			user body dto.CreateUserRequest true "Create User"
// @Produce			application/json
// @Tags			Users
// @Success			201 {object} dto.Response{}
// @Router			/users [post]
func (uc *UserController) Create(ctx *gin.Context) {
	createUserRequest := dto.CreateUserRequest{}

	err := ctx.ShouldBindJSON(&createUserRequest)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	err = uc.userService.CreateUser(createUserRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	webResponse := dto.Response{
		Data: nil,
	}

	ctx.JSON(http.StatusCreated, webResponse)

}

// FindMany				godoc
// @Summary				Gets a list of users.
// @Description			Returns a list of users
// @Produce				application/json
// @Tags				Users
// @Success				200 {object} dto.Response{}
// @Router				/users [get]
func (uc *UserController) FindMany(ctx *gin.Context) {
	users, err := uc.userService.FindManyUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)

}

// FindById 			godoc
// @Summary				Get Single user by id.
// @Param				userId path string true "User ID"
// @Description			Returns the user matching the requested id
// @Produce				application/json
// @Tags				Users
// @Success				200 {object} dto.FindUserResponse{}
// @Router				/users/{userId} [get]
func (uc *UserController) FindById(ctx *gin.Context) {
	userId := ctx.Param("userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.userService.FindUserById(uint(id))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)

}

// FindUserProfile 	godoc
// @Summary				Get Single user with it's profile by id.
// @Param				userId path string true "User ID"
// @Description			Returns the user matching the requested id
// @Produce				application/json
// @Tags				Users
// @Success				200 {object} dto.FindUserWithProfileByIdResponse{}
// @Router				/users/{userId}/profile [get]
func (uc *UserController) FindUserProfile(ctx *gin.Context) {
	userId := helpers.MustParseUint(ctx.Param("userId"))
	userWithProfile, err := uc.userService.FindUserWithProfileById(userId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, userWithProfile)

}

// FindMemberships	godoc
// @Summary					Gets a list of communities that user is member of.
// @Param					userId path string true "User ID"
// @Description				Returns a list of communities
// @Produce					application/json
// @Tags					Users
// @Success					200 {object} dto.FindUserMembershipsResponse{}
// @Router					/users/{userId}/memberships [get]
func (uc *UserController) FindMemberships(ctx *gin.Context) {
	userId := helpers.MustParseUint(ctx.Param("userId"))

	memberships, err := uc.userService.FindUserMemberships(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, memberships)
}

// FindFriends	godoc
// @Summary					Gets a list of friends of the user.
// @Param					userId path string true "User ID"
// @Description				Returns a list of friends
// @Produce					application/json
// @Tags					Users
// @Success					200 {object} dto.FindFriendsOfUserResponse{}
// @Router					/users/{userId}/friends [get]
func (uc *UserController) FindFriends(ctx *gin.Context) {
	userId := helpers.MustParseUint(ctx.Param("userId"))
	friendships, err := uc.userService.FindFriendsOfUser(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, friendships)
}

// FindPhotoAlbums	godoc
// @Summary					Gets a list of photo albums of the user.
// @Param					userId path string true "User ID"
// @Description				Returns a list of photo albums
// @Produce					application/json
// @Tags					Users
// @Success					200 {object} dto.FindPhotoAlbumOfUserResponse{}
// @Router					/users/{userId}/photoAlbums [get]
func (uc *UserController) FindPhotoAlbums(ctx *gin.Context) {
	userId := helpers.MustParseUint(ctx.Param("userId"))
	photoAlbums, err := uc.userService.FindUserPhotoAlbums(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, photoAlbums)
}
