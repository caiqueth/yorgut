package controller

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"yorgut/internal/helpers"
	"yorgut/internal/model"
	"yorgut/internal/service"
	"yorgut/internal/transport"

	"github.com/gin-gonic/gin"
)

type MeController struct {
	userService                service.UserService
	friendshipService          service.FriendshipService
	communityService           service.CommunityService
	communityMembershipService service.CommunityMembershipService
	userFanService             service.UserFanService
	scrapService               service.ScrapService
	userProfileService         service.UserProfileService
}

func NewMeController(
	userService service.UserService,
	friendshipService service.FriendshipService,
	communityService service.CommunityService,
	communityMembershipService service.CommunityMembershipService,
	userFanService service.UserFanService,
	scrapService service.ScrapService,
	userProfileService service.UserProfileService,
) *MeController {
	return &MeController{
		userService:                userService,
		friendshipService:          friendshipService,
		communityService:           communityService,
		communityMembershipService: communityMembershipService,
		userFanService:             userFanService,
		scrapService:               scrapService,
		userProfileService:         userProfileService,
	}
}

// FindMe	 			godoc
// @ID					findMe
// @Summary				Get the logged user.
// @Description			Returns the user that owns the session token
// @Produce				application/json
// @Tags				Me
// @Success				200 {object} transport.FindUserResponse{}
// @Security			BearerAuth
// @Router				/me [get]
func (uc *MeController) FindMe(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	resp := transport.FindUserResponse{}
	resp.FromUser(user)

	ctx.JSON(http.StatusOK, resp)
}

// Welcome	 			godoc
// @ID					findMyWelcome
// @Summary				Get user stats for the Home Welcome
// @Description			Returns the user stats
// @Produce				application/json
// @Tags				Me
// @Success				200 {object} transport.WelcomeResponse{}
// @Security			BearerAuth
// @Router				/me/welcome [get]
func (uc *MeController) Welcome(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	numberOfFans, err := uc.userFanService.CountFansOfUser(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	numberOfScraps, err := uc.scrapService.CountScrapsOfUser(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := transport.WelcomeResponse{
		UserFirstName:       user.FirstName,
		UserCreatedAt:       user.CreatedAt,
		IndirectConnections: 0,
		DirectConnections:   user.FriendCount,
		NumberOfFans:        numberOfFans,
		NumberOfScraps:      numberOfScraps,
		ProfileViews: transport.ProfileViews{
			AllTime:  0,
			LastWeek: 0,
		},
		RecentVisitors: make([]transport.RecentVisitor, 0),
		TodaysLuck:     "Boas notícias chegarão por e-mail.",
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateMe		godoc
// @ID					updateMe
// @Summary			Update the logged user
// @Description		Replaces the logged user's basic data and profile fields. Email can't be changed.
// @Param			Me body transport.UpdateMeRequest true "User data"
// @Accept			application/json
// @Produce			application/json
// @Tags			Me
// @Success			204
// @Security			BearerAuth
// @Router			/me [put]
func (uc *MeController) UpdateMe(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	request := transport.UpdateMeRequest{}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var birthDate *time.Time
	if request.BirthDate != "" {
		parsed, err := time.ParseInLocation("2006-01-02", request.BirthDate, time.Local)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		birthDate = &parsed
	}

	err := uc.userService.UpdateUser(user.ID, service.UpdateUserInput{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		BirthDate: birthDate,
	})
	if errors.Is(err, service.ErrFirstNameRequired) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = uc.userProfileService.UpdateProfile(user.ID, service.UpdateUserProfileInput{
		Description:        request.Description,
		RelationshipStatus: request.RelationshipStatus,
		Interests:          request.Interests,
		About:              request.About,
		Religion:           request.Religion,
		Humor:              request.Humor,
		Fashion:            request.Fashion,
		LivesWith:          request.LivesWith,
		BirthCity:          request.BirthCity,
		Website:            request.Website,
		Passions:           request.Passions,
		Sports:             request.Sports,
		Hobbies:            request.Hobbies,
		Books:              request.Books,
		Music:              request.Music,
		Movies:             request.Movies,
		Cuisine:            request.Cuisine,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// UpdateMyProfilePicture	godoc
// @ID					updateMyProfilePicture
// @Summary			Update the logged user's profile picture
// @Description		Uploads a PNG or JPEG (max 1 MB and 1024x1024) as the logged user's profile picture.
// @Param			picture formData file true "Profile picture"
// @Accept			multipart/form-data
// @Produce			application/json
// @Tags			Me
// @Success			204
// @Security			BearerAuth
// @Router			/me/profile-picture [put]
func (uc *MeController) UpdateMyProfilePicture(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	file, err := ctx.FormFile("picture")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get picture file"})
		return
	}
	// Refuse oversized files before reading them; the service checks again.
	if file.Size > service.MaxPictureBytes {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": service.ErrPictureTooLarge.Error()})
		return
	}

	picture, err := helpers.FileToBytes(file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read picture file"})
		return
	}

	err = uc.userProfileService.UpdateProfilePicture(user.ID, picture)
	if errors.Is(err, service.ErrInvalidPicture) ||
		errors.Is(err, service.ErrPictureTooLarge) ||
		errors.Is(err, service.ErrPictureTooBig) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// FindMyProfile 			godoc
// @ID					findMyProfile
// @Summary					Get logged user with it's profile by their Token.
// @Description				Returns the user matching the token
// @Produce					application/json
// @Tags					Me
// @Success					200 {object} transport.FindUserWithProfileByIdResponse{}
// @Security			BearerAuth
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

// FindMyFriends	godoc
// @ID					findMyFriends
// @Summary			Gets a list of friends of the user.
// @Description		Returns a list of friends
// @Produce			application/json
// @Tags			Me
// @Success			200 {object} transport.FindFriendsOfUserResponse{}
// @Security			BearerAuth
// @Router			/me/friends [get]
func (uc *MeController) FindMyFriends(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	friendships, err := uc.friendshipService.FindFriendsOfUser(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, friendships)
}

// FindMyPendingFriendships	godoc
// @ID					findMyPendingFriendships
// @Summary					Gets a list of friendship invites of the user.
// @Description				Returns a list of users
// @Produce					application/json
// @Tags					Me
// @Success					200 {object} transport.FindPendingFriendshipOfUserResponse{}
// @Security			BearerAuth
// @Router					/me/friends/pending [get]
func (uc *MeController) FindMyPendingFriendships(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	friendships, err := uc.friendshipService.FindPendingFriendshipsOfUser(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, friendships)
}

// JoinCommunity	godoc
// @ID					joinCommunity
// @Summary			User joins a community.
// @Param			communityId path int true "Community ID"
// @Description		Joins a community.
// @Produce			application/json
// @Tags			Me
// @Success			201 {object} transport.Response{}
// @Security			BearerAuth
// @Router			/me/join/{communityId} [post]
func (uc *MeController) JoinCommunity(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	communityId, err := strconv.Atoi(ctx.Param("communityId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = uc.communityMembershipService.JoinCommunity(user.ID, uint(communityId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, transport.Response{})
}

// LeaveCommunity	godoc
// @ID					leaveCommunity
// @Summary			User leaves a community.
// @Param			communityId path int true "Community ID"
// @Description		Leaves a community. The owner can't leave (403).
// @Produce			application/json
// @Tags			Me
// @Success			201 {object} transport.Response{}
// @Failure			403
// @Security			BearerAuth
// @Router			/me/leave/{communityId} [post]
func (uc *MeController) LeaveCommunity(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	communityId, err := strconv.Atoi(ctx.Param("communityId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = uc.communityMembershipService.LeaveCommunity(user.ID, uint(communityId))
	if errors.Is(err, service.ErrOwnerCannotLeave) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, transport.Response{})
}

// FindMyMemberships		godoc
// @ID					findMyMemberships
// @Summary					Gets a list of communities that user is member of.
// @Description				Returns a list of communities
// @Produce					application/json
// @Tags					Me
// @Success					200 {object} transport.FindUserMembershipsResponse{}
// @Security			BearerAuth
// @Router					/me/memberships [get]
func (uc *MeController) FindMyMemberships(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	memberships, err := uc.communityMembershipService.FindUserMemberships(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, memberships)
}

// FindMyCommunities		godoc
// @ID					findMyCommunities
// @Summary					Gets a list of communities that user owns.
// @Description				Returns a list of communities
// @Produce					application/json
// @Tags					Me
// @Success					200 {object} transport.FindMyOwnedCommunitiesResponse{}
// @Security			BearerAuth
// @Router					/me/my-communities [get]
func (uc *MeController) FindMyCommunities(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	communities, err := uc.communityService.FindCommunitiesCreatedByUser(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, communities)
}

// AmIMemberOf		godoc
// @ID					amIMemberOf
// @Summary			Evaluates if the logged user is member of specified community
// @Param			communityId path int true "Community ID"
// @Description		Returns true or false
// @Produce			application/json
// @Tags			Me
// @Success			200 {object} transport.AmIMemberOfResponse{}
// @Security			BearerAuth
// @Router			/me/memberships/{communityId} [get]
func (uc *MeController) AmIMemberOf(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	communityId := helpers.MustParseUint(ctx.Param("communityId"))

	isMember := uc.communityMembershipService.IsUserMemberOf(communityId, user.ID)

	ctx.JSON(http.StatusOK, transport.AmIMemberOfResponse{
		AmIMember: isMember,
	})
}

// AmIFriendOf		godoc
// @ID					amIFriendOf
// @Summary			Evaluates if the logged user is friend of specified user
// @Param			userId path int true "User ID"
// @Description		Returns the friendship status between the logged user and userId. Answers 204 with no body when userId is the logged user.
// @Produce			application/json
// @Tags			Me
// @Success			200 {object} transport.AmIFriendOfResponse{}
// @Success			204 "userId is the logged user"
// @Security			BearerAuth
// @Router			/me/friends/{userId} [get]
func (uc *MeController) AmIFriendOf(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	targetUserId := helpers.MustParseUint(ctx.Param("userId"))
	if user.ID == targetUserId {
		ctx.JSON(http.StatusNoContent, nil)
		return
	}

	response := uc.friendshipService.UsersFriendshipStatus(user.ID, targetUserId)

	ctx.JSON(http.StatusOK, response)
}
