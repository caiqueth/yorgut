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
// @ID					createUser
// @Summary			Create user
// @Description		Save user in Db.
// @Param			user body transport.CreateUserRequest true "Create User"
// @Accept			application/json
// @Produce			application/json
// @Tags			Users
// @Success			201 {object} transport.Response{}
// @Failure			400
// @Failure			401
// @Failure			409
// @Security			BearerAuth
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
		RespondError(ctx, err)
		return
	}

	webResponse := transport.Response{
		Data: nil,
	}

	ctx.JSON(http.StatusCreated, webResponse)

}

// FindMany				godoc
// @ID					findUsers
// @Summary				Gets a list of users.
// @Description			Returns a list of users
// @Produce				application/json
// @Tags				Users
// @Success				200 {object} transport.FindManyUsersResponse{}
// @Failure			401
// @Security			BearerAuth
// @Router				/users [get]
func (uc *UserController) FindMany(ctx *gin.Context) {
	users, err := uc.userService.FindManyUsers()
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, users)

}

// SearchUsers			godoc
// @ID					searchUsers
// @Summary				Searches for a user by their name.
// @Param				name query string false "Part of the user's full name (omit to list all)"
// @Param				page query int false "Page, starting at 1"
// @Param				pageSize query int false "Users per page"
// @Description			Returns a page of users
// @Produce				application/json
// @Tags				Users
// @Success				200 {object} transport.SearchUsersResponse{}
// @Failure			401
// @Security			BearerAuth
// @Router				/users/search-users [get]
func (uc *UserController) SearchUsers(ctx *gin.Context) {
	name := ctx.Query("name")
	page := helpers.MustParseInt(ctx.Query("page"))
	pageSize := helpers.MustParseInt(ctx.Query("pageSize"))

	users := uc.userService.SearchUsers(name, page, pageSize)

	ctx.JSON(http.StatusOK, users)

}

// FindById 			godoc
// @ID					findUser
// @Summary				Get Single user by id.
// @Param				userId path int true "User ID"
// @Description			Returns the user matching the requested id
// @Produce				application/json
// @Tags				Users
// @Success				200 {object} transport.FindUserResponse{}
// @Failure			400
// @Failure			401
// @Failure			404
// @Security			BearerAuth
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
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, user)

}

// FindUserProfile 	godoc
// @ID					findUserProfile
// @Summary				Get Single user with it's profile by id.
// @Param				userId path int true "User ID"
// @Description			Returns the user matching the requested id
// @Produce				application/json
// @Tags				Users
// @Success				200 {object} transport.FindUserWithProfileByIdResponse{}
// @Failure			401
// @Failure			404
// @Security			BearerAuth
// @Router				/users/{userId}/profile [get]
func (uc *UserController) FindUserProfile(ctx *gin.Context) {
	userId := helpers.MustParseUint(ctx.Param("userId"))
	userWithProfile, err := uc.userService.FindUserWithProfileById(userId)

	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, userWithProfile)

}

// FindMemberships	godoc
// @ID					findUserMemberships
// @Summary					Gets a list of communities that user is member of.
// @Param					userId path int true "User ID"
// @Description				Returns a list of communities
// @Produce					application/json
// @Tags					Users
// @Success					200 {object} transport.FindUserMembershipsResponse{}
// @Failure			401
// @Security			BearerAuth
// @Router					/users/{userId}/memberships [get]
func (uc *UserController) FindMemberships(ctx *gin.Context) {
	userId := helpers.MustParseUint(ctx.Param("userId"))

	memberships, err := uc.communityMembershipService.FindUserMemberships(userId)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, memberships)
}

// FindFriends	godoc
// @ID					findUserFriends
// @Summary					Gets a list of friends of the user.
// @Param					userId path int true "User ID"
// @Description				Returns a list of friends
// @Produce					application/json
// @Tags					Users
// @Success					200 {object} transport.FindFriendsOfUserResponse{}
// @Failure			401
// @Security			BearerAuth
// @Router					/users/{userId}/friends [get]
func (uc *UserController) FindFriends(ctx *gin.Context) {
	userId := helpers.MustParseUint(ctx.Param("userId"))
	friendships, err := uc.friendshipService.FindFriendsOfUser(userId)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, friendships)
}
