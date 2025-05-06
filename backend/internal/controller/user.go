package controller

import (
	"net/http"
	"strconv"
	"yorgut/internal/dto"
	"yorgut/internal/model"
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

// FindWithToken		godoc
// @Summary				Get Single user with the token.
// @Param				userId path string true "get user by id"
// @Description			Returns the user matching the requested id
// @Produce				application/json
// @Tags				User
// @Success				200 {object} dto.FindUserByIdResponse{}
// @Router				/user/me [get]
func (uc *UserController) FindWithToken(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	ctx.JSON(http.StatusOK, user)
}

// FindUserProfile 	godoc
// @Summary				Get Single user with it's profile by id.
// @Param				userId path string true "get user by id"
// @Description			Returns the user matching the requested id
// @Produce				application/json
// @Tags				User
// @Success				200 {object} dto.FindUserWithProfileByIdResponse{}
// @Router				/user/profile [get]
func (uc *UserController) FindUserProfile(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	userWithProfile, err := uc.userService.FindUserWithProfileById(user.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, userWithProfile)

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
// @Router			/user [delete]
func (uc *UserController) Delete(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	uc.userService.DeleteUser(user.ID)

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
	// THIS ENDPOINT IS BROKEN
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
// @Router					/user/photoAlbums [get]
func (uc *UserController) FindPhotoAlbums(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	photoAlbums, err := uc.userService.FindUserPhotoAlbums(user.ID)
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
// @Router					/user/friends [get]
func (uc *UserController) FindFriends(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	friendships, err := uc.userService.FindFriendsOfUser(user.ID)
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
// @Router					/user/friends [get]
func (uc *UserController) FindPendingFriendships(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	friendships, err := uc.userService.FindPendingFriendshipsOfUser(user.ID)
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
// @Router			/user/join/{communityId} [post]
func (uc *UserController) JoinCommunity(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	communityId, err := strconv.Atoi(ctx.Param("communityId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = uc.userService.JoinCommunity(user.ID, uint(communityId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, dto.Response{})
}

// LeaveCommunity	godoc
// @Summary			User leaves a community.
// @Param			userId path string true "user Id"
// @Param			communityId path string true "community Id"
// @Description		Leaves a community.
// @Produce			application/json
// @Tags			User
// @Success			201 {object} dto.Response{}
// @Router			/user/leave/{communityId} [post]
func (uc *UserController) LeaveCommunity(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	communityId, err := strconv.Atoi(ctx.Param("communityId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = uc.userService.LeaveCommunity(user.ID, uint(communityId))
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
// @Router					/user/memberships [get]
func (uc *UserController) FindMemberships(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	memberships, err := uc.userService.FindUserMemberships(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, memberships)
}
