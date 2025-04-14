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
// @Tags			User
// @Success			201 {object} dto.Response{}
// @Router			/user [post]
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

// FindById 			godoc
// @Summary				Get Single user by id.
// @Param				userId path string true "get user by id"
// @Description			Returns the user matching the requested id
// @Produce				application/json
// @Tags				User
// @Success				200 {object} dto.FindUserByIdResponse{}
// @Router				/user/{userId} [get]
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

// FindWithProfileById 	godoc
// @Summary				Get Single user with it's profile by id.
// @Param				userId path string true "get user by id"
// @Description			Returns the user matching the requested id
// @Produce				application/json
// @Tags				User
// @Success				200 {object} dto.FindUserWithProfileByIdResponse{}
// @Router				/user/{userId}/profile [get]
func (uc *UserController) FindWithProfileById(ctx *gin.Context) {
	userId := ctx.Param("userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.userService.FindUserWithProfileById(uint(id))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)

}

// FindMany				godoc
// @Summary				Gets a list of users.
// @Description			Returns a list of users
// @Produce				application/json
// @Tags				User
// @Success				200 {object} dto.Response{}
// @Router				/user [get]
func (uc *UserController) FindMany(ctx *gin.Context) {
	users, err := uc.userService.FindManyUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)

}

// Delete 		godoc
// @Summary			Delete User by id
// @Param			userId path string true "delete user by id"
// @Description		Deletes the user matching the requested id
// @Produce			application/json
// @Tags			User
// @Success			200 {object} dto.Response{}
// @Router			/user/{userId} [delete]
func (uc *UserController) Delete(ctx *gin.Context) {
	userId := ctx.Param("userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	uc.userService.DeleteUser(uint(id))

	webResponse := dto.Response{
		Data: nil,
	}
	ctx.JSON(http.StatusOK, webResponse)

}

// Update 			godoc
// @Summary			Update user
// @Param			userId path string true "update user"
// @Description		Returns null
// @Produce			application/json
// @Tags			User
// @Success			200 {object} dto.Response{}
// @Router			/user [put]
func (uc *UserController) Update(ctx *gin.Context) {
	request := dto.UpdateUserRequest{}

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = uc.userService.UpdateUser(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	webResponse := dto.Response{
		Data: nil,
	}

	ctx.JSON(http.StatusOK, webResponse)

}

// FindPhotoAlbums	godoc
// @Summary					Gets a list of photo albums of the user.
// @Param					userId path string true "user id"
// @Description				Returns a list of photo albums
// @Produce					application/json
// @Tags					User
// @Success					200 {object} dto.FindPhotoAlbumOfUserResponse{}
// @Router					/user/{userId}/photoAlbums [get]
func (uc *UserController) FindPhotoAlbums(ctx *gin.Context) {
	userId := ctx.Param("userId")
	id, err := strconv.Atoi(userId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	photoAlbums, err := uc.userService.FindUserPhotoAlbums(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, photoAlbums)
}

// FindFriends	godoc
// @Summary					Gets a list of friends of the user.
// @Param					userId path string true "user id"
// @Description				Returns a list of friends
// @Produce					application/json
// @Tags					User
// @Success					200 {object} dto.FindFriendsOfUserResponse{}
// @Router					/user/{userId}/friends [get]
func (uc *UserController) FindFriends(ctx *gin.Context) {
	userId := helpers.MustParseUint(ctx.Param("userId"))

	friendships, err := uc.userService.FindFriendsOfUser(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, friendships)
}

// FindPendingFriendships	godoc
// @Summary					Gets a list of friendship invites of the user.
// @Param					userId path string true "user id"
// @Description				Returns a list of users
// @Produce					application/json
// @Tags					User
// @Success					200 {object} dto.FindFriendsOfUserResponse{}
// @Router					/user/{userId}/friends [get]
func (uc *UserController) FindPendingFriendships(ctx *gin.Context) {
	userId := helpers.MustParseUint(ctx.Param("userId"))

	friendships, err := uc.userService.FindPendingFriendshipsOfUser(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, friendships)
}

func (uc *UserController) FindCommunitiesOwned(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("userId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	communities, err := uc.userService.FindUserOwnedCommunities(uint(userId))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, communities)
}

// JoinCommunity	godoc
// @Summary			User joins a community.
// @Param			userId path string true "user Id"
// @Param			communityId path string true "community Id"
// @Description		Joins a community.
// @Produce			application/json
// @Tags			User
// @Success			201 {object} dto.Response{}
// @Router			/user/{userId}/join/{communityId} [post]
func (uc *UserController) JoinCommunity(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("userId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	communityId, err := strconv.Atoi(ctx.Param("communityId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = uc.userService.JoinCommunity(uint(userId), uint(communityId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, dto.Response{})
}

// FindCommunitiesMemberOf	godoc
// @Summary					Gets a list of communities that user is member of.
// @Param					userId path string true "get user memberships by id"
// @Description				Returns a list of communities
// @Produce					application/json
// @Tags					User
// @Success					200 {object} dto.FindUserMembershipsResponse{}
// @Router					/user/{userId}/memberships [get]
func (uc *UserController) FindMemberships(ctx *gin.Context) {
	userId := ctx.Param("userId")
	id, err := strconv.Atoi(userId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	memberships, err := uc.userService.FindUserMemberships(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, memberships)
}

func (uc *UserController) FindAuthoredTestimonials(ctx *gin.Context) {
	panic("unimplemented")
}

func (uc *UserController) FindReceivedTestimonials(ctx *gin.Context) {
	panic("unimplemented")
}
