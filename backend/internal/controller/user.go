package controller

import (
	"net/http"
	"strconv"
	"yorgut/internal/helpers"
	"yorgut/internal/service"
	"yorgut/internal/transport"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService                service.UserService
	friendshipService          service.FriendshipService
	communityMembershipService service.CommunityMembershipService
}

func NewUserController(
	userService service.UserService,
	friendshipService service.FriendshipService,
	communityMembershipService service.CommunityMembershipService,
) *UserController {
	return &UserController{
		userService:                userService,
		friendshipService:          friendshipService,
		communityMembershipService: communityMembershipService,
	}
}

// Create			godoc
// @Summary			Create user
// @Decription		Save user in Db.
// @Param			user body transport.CreateUserRequest true "Create User"
// @Produce			application/json
// @Tags			Users
// @Success			201 {object} transport.Response{}
// @Router			/users [post]
func (uc *UserController) Create(ctx *gin.Context) {
	req := transport.CreateUserRequest{}

	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = uc.userService.CreateUser(req.FirstName, req.LastName, req.Email, req.BirthDate, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	webResponse := transport.Response{
		Data: nil,
	}

	ctx.JSON(http.StatusCreated, webResponse)

}

// FindMany				godoc
// @Summary				Gets a list of users.
// @Description			Returns a list of users
// @Produce				application/json
// @Tags				Users
// @Success				200 {object} transport.FindManyUsersResponse{}
// @Router				/users [get]
func (uc *UserController) FindMany(ctx *gin.Context) {
	users, err := uc.userService.FindManyUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)

}

// SearchUser			godoc
// @Summary				Searches for a user by their name.
// @Description			Returns a list of users
// @Produce				application/json
// @Tags				Users
// @Success				200 {object} transport.SearchUsersResponse{}
// @Router				/users/search-users [get]
func (uc *UserController) SearchUsers(ctx *gin.Context) {
	name := ctx.Query("name")
	page := helpers.MustParseInt(ctx.Query("page"))
	pageSize := helpers.MustParseInt(ctx.Query("pageSize"))

	users := uc.userService.SearchUsers(name, page, pageSize)

	ctx.JSON(http.StatusOK, users)

}

// FindById 			godoc
// @Summary				Get Single user by id.
// @Param				userId path string true "User ID"
// @Description			Returns the user matching the requested id
// @Produce				application/json
// @Tags				Users
// @Success				200 {object} transport.FindUserResponse{}
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
// @Success				200 {object} transport.FindUserWithProfileByIdResponse{}
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
// @Success					200 {object} transport.FindUserMembershipsResponse{}
// @Router					/users/{userId}/memberships [get]
func (uc *UserController) FindMemberships(ctx *gin.Context) {
	userId := helpers.MustParseUint(ctx.Param("userId"))

	memberships, err := uc.communityMembershipService.FindUserMemberships(userId)
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
// @Success					200 {object} transport.FindFriendsOfUserResponse{}
// @Router					/users/{userId}/friends [get]
func (uc *UserController) FindFriends(ctx *gin.Context) {
	userId := helpers.MustParseUint(ctx.Param("userId"))
	friendships, err := uc.friendshipService.FindFriendsOfUser(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, friendships)
}
