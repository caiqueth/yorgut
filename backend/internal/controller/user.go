package controller

import (
	"net/http"
	"yorgut/internal/model"
	"yorgut/internal/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService                service.UserService
	friendshipService          service.FriendshipService
	communityMembershipService service.CommunityMembershipService
	profileVisitService        service.ProfileVisitService
}

func NewUserController(
	userService service.UserService,
	friendshipService service.FriendshipService,
	communityMembershipService service.CommunityMembershipService,
	profileVisitService service.ProfileVisitService,
) *UserController {
	return &UserController{
		userService:                userService,
		friendshipService:          friendshipService,
		communityMembershipService: communityMembershipService,
		profileVisitService:        profileVisitService,
	}
}

// SearchUsers			godoc
// @ID					searchUsers
// @Summary				Searches for a user by their name.
// @Param				name query string false "Part of the user's full name (omit to list all)"
// @Param				page query int false "Page, starting at 1 (default 1)"
// @Param				pageSize query int false "Users per page (default 20, max 50)"
// @Description			Returns a page of users
// @Produce				application/json
// @Tags				Users
// @Success				200 {object} transport.SearchUsersResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/users/search-users [get]
func (u *UserController) SearchUsers(ctx *gin.Context) {
	name := ctx.Query("name")
	page, pageSize := pageQuery(ctx)

	users, err := u.userService.SearchUsers(name, page, pageSize)
	if err != nil {
		RespondError(ctx, err)
		return
	}

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
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/users/{userId} [get]
func (u *UserController) FindById(ctx *gin.Context) {
	id, ok := idParam(ctx, "userId")
	if !ok {
		return
	}

	user, err := u.userService.FindUserById(id)

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
// @Param				as query string false "\"stranger\": as someone who isn't a friend sees it (the profile preview)"
// @Description			Returns the user matching the requested id, and counts a visit to their profile (see GET /me/welcome)
// @Produce				application/json
// @Tags				Users
// @Success				200 {object} transport.FindUserWithProfileByIdResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/users/{userId}/profile [get]
func (u *UserController) FindUserProfile(ctx *gin.Context) {
	viewer := ctx.MustGet("user").(*model.User)
	userId, ok := idParam(ctx, "userId")
	if !ok {
		return
	}
	userWithProfile, err := u.userService.FindUserWithProfileById(viewerOf(ctx), userId)

	if err != nil {
		RespondError(ctx, err)
		return
	}

	u.profileVisitService.RecordVisit(viewer.ID, userId)

	ctx.JSON(http.StatusOK, userWithProfile)
}

// FindMemberships	godoc
// @ID					findUserMemberships
// @Summary					Gets a page of the communities a user is a member of.
// @Param					userId path int true "User ID"
// @Param				as query string false "\"stranger\": as someone who isn't a friend sees it (the profile preview)"
// @Param				page query int false "Page, starting at 1 (default 1)"
// @Param				pageSize query int false "Items per page (default 20, max 50)"
// @Param				random query bool false "Shuffled instead of by name, the same way all day (for the profile box)"
// @Description				Returns a page of communities
// @Produce					application/json
// @Tags					Users
// @Success					200 {object} transport.FindUserMembershipsResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router					/users/{userId}/memberships [get]
func (u *UserController) FindMemberships(ctx *gin.Context) {
	userId, ok := idParam(ctx, "userId")
	if !ok {
		return
	}

	page, pageSize := pageQuery(ctx)
	random := ctx.Query("random") == "true"
	memberships, err := u.communityMembershipService.FindUserMemberships(viewerOf(ctx), userId, page, pageSize, random)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, memberships)
}

// FindFriends	godoc
// @ID					findUserFriends
// @Summary					Gets a page of a user's friends.
// @Param					userId path int true "User ID"
// @Param				page query int false "Page, starting at 1 (default 1)"
// @Param				pageSize query int false "Items per page (default 20, max 50)"
// @Description				Returns a page of friends, most recently seen first
// @Produce					application/json
// @Tags					Users
// @Success					200 {object} transport.FindFriendsOfUserResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router					/users/{userId}/friends [get]
func (u *UserController) FindFriends(ctx *gin.Context) {
	userId, ok := idParam(ctx, "userId")
	if !ok {
		return
	}
	page, pageSize := pageQuery(ctx)
	friendships, err := u.friendshipService.FindFriendsOfUser(userId, page, pageSize)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, friendships)
}
