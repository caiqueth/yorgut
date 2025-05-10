package controller

import (
	"net/http"
	"strconv"
	"yorgut/internal/dto"
	"yorgut/internal/helpers"
	"yorgut/internal/model"
	"yorgut/internal/service"

	"github.com/gin-gonic/gin"
)

type MeController struct {
	userService service.UserService
}

func NewMeController(service service.UserService) *MeController {
	return &MeController{
		userService: service,
	}
}

// FindMe	 			godoc
// @Summary				Get Single user by id.
// @Description			Returns the user matching the requested id
// @Produce				application/json
// @Tags				Me
// @Success				200 {object} dto.FindUserResponse{}
// @Router				/me [get]
func (uc *MeController) FindMe(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	resp := dto.FindUserResponse{}
	resp.FromUser(user)

	ctx.JSON(http.StatusOK, resp)
}

// Update 			godoc
// @Summary			Update user
// @Description		Returns null
// @Produce			application/json
// @Tags			Me
// @Success			200 {object} dto.Response{}
// @Router			/me [put]
func (uc *MeController) UpdateMe(ctx *gin.Context) {
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

// FindLoggedUserProfile 	godoc
// @Summary					Get logged user with it's profile by their Token.
// @Description				Returns the user matching the token
// @Produce					application/json
// @Tags					Me
// @Success					200 {object} dto.FindUserWithProfileByIdResponse{}
// @Router					/me/profile [get]
func (uc *MeController) FindMyProfile(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	userWithProfile, err := uc.userService.FindUserWithProfileById(user.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, userWithProfile)

}

// FindMyPhotoAlbums		godoc
// @Summary					Gets a list of photo albums of the user.
// @Description				Returns a list of photo albums
// @Produce					application/json
// @Tags					Me
// @Success					200 {object} dto.FindPhotoAlbumOfUserResponse{}
// @Router					/me/photoAlbums [get]
func (uc *MeController) FindMyPhotoAlbums(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	photoAlbums, err := uc.userService.FindUserPhotoAlbums(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, photoAlbums)
}

// FindMyFriends	godoc
// @Summary			Gets a list of friends of the user.
// @Description		Returns a list of friends
// @Produce			application/json
// @Tags			Me
// @Success			200 {object} dto.FindFriendsOfUserResponse{}
// @Router			/me/friends [get]
func (uc *MeController) FindMyFriends(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	friendships, err := uc.userService.FindFriendsOfUser(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, friendships)
}

// FindMyPendingFriendships	godoc
// @Summary					Gets a list of friendship invites of the user.
// @Description				Returns a list of users
// @Produce					application/json
// @Tags					Me
// @Success					200 {object} dto.FindPendingFriendshipOfUserResponse{}
// @Router					/me/friends/pending [get]
func (uc *MeController) FindMyPendingFriendships(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	friendships, err := uc.userService.FindPendingFriendshipsOfUser(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, friendships)
}

// JoinCommunity	godoc
// @Summary			User joins a community.
// @Param			communityId path string true "Community ID"
// @Description		Joins a community.
// @Produce			application/json
// @Tags			Me
// @Success			201 {object} dto.Response{}
// @Router			/me/join/{communityId} [post]
func (uc *MeController) JoinCommunity(ctx *gin.Context) {
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
// @Param			communityId path string true "Community ID"
// @Description		Leaves a community.
// @Produce			application/json
// @Tags			Me
// @Success			201 {object} dto.Response{}
// @Router			/me/leave/{communityId} [post]
func (uc *MeController) LeaveCommunity(ctx *gin.Context) {
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

// FindMyMemberships		godoc
// @Summary					Gets a list of communities that user is member of.
// @Description				Returns a list of communities
// @Produce					application/json
// @Tags					Me
// @Success					200 {object} dto.FindUserMembershipsResponse{}
// @Router					/me/memberships [get]
func (uc *MeController) FindMyMemberships(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	memberships, err := uc.userService.FindUserMemberships(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	memberships.MembershipCount = uc.userService.CountUserMemberships(user.ID)

	ctx.JSON(http.StatusOK, memberships)
}

// AmIMemberOf		godoc
// @Summary			Evaluates if the logged user is member of specified community
// @Description		Returns true or false
// @Produce			application/json
// @Tags			Me
// @Success			200 {object} dto.AmIMemberOfResponse{}
// @Router			/me/memberships/{communityId} [get]
func (uc *MeController) AmIMemberOf(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	communityId := helpers.MustParseUint(ctx.Param("communityId"))

	isMember := uc.userService.IsUserMemberOf(communityId, user.ID)

	ctx.JSON(http.StatusOK, dto.AmIMemberOfResponse{
		AmIMember: isMember,
	})
}
