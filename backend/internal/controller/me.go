package controller

import (
	"net/http"
	"time"
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
	profileVisitService        service.ProfileVisitService
}

func NewMeController(
	userService service.UserService,
	friendshipService service.FriendshipService,
	communityService service.CommunityService,
	communityMembershipService service.CommunityMembershipService,
	userFanService service.UserFanService,
	scrapService service.ScrapService,
	userProfileService service.UserProfileService,
	profileVisitService service.ProfileVisitService,
) *MeController {
	return &MeController{
		userService:                userService,
		friendshipService:          friendshipService,
		communityService:           communityService,
		communityMembershipService: communityMembershipService,
		userFanService:             userFanService,
		scrapService:               scrapService,
		userProfileService:         userProfileService,
		profileVisitService:        profileVisitService,
	}
}

// FindMe	 			godoc
// @ID					findMe
// @Summary				Get the logged user.
// @Description			Returns the user that owns the session token, with their settings
// @Produce				application/json
// @Tags				Me
// @Success				200 {object} transport.FindMeResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/me [get]
func (m *MeController) FindMe(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	resp := transport.FindMeResponse{}
	resp.FromUser(user)

	ctx.JSON(http.StatusOK, resp)
}

// UpdateMyPrivacy		godoc
// @ID					updateMyPrivacy
// @Summary				Changes the logged user's privacy settings.
// @Param				Privacy body transport.UpdateMyPrivacyRequest true "Privacy settings"
// @Description			recordVisits false stops recording the user's visits to profiles and hides them from every visitor list (theirs too).
// @Accept				application/json
// @Produce				application/json
// @Tags				Me
// @Success				204
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/me/privacy [put]
func (m *MeController) UpdateMyPrivacy(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	request := transport.UpdateMyPrivacyRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondInvalidRequest(ctx, err)
		return
	}

	err := m.userService.UpdatePrivacy(user.ID, service.PrivacyInput{
		RecordVisits:        *request.RecordVisits,
		PublicAlbum:         *request.PublicAlbum,
		PublicCommunities:   *request.PublicCommunities,
		PublicScraps:        *request.PublicScraps,
		PublicTestimonials:  *request.PublicTestimonials,
		PublicBirthday:      *request.PublicBirthday,
		HiddenProfileFields: request.HiddenProfileFields,
	})
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// DeleteMe			godoc
// @ID					deleteMe
// @Summary				Deletes the logged user's account, for good.
// @Description			Deletes the user and everything they wrote or did in the app. Each community they own goes to its longest-standing other member, or is deleted when nobody else is in it.
// @Produce				application/json
// @Tags				Me
// @Success				204
// @Failure			401 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/me [delete]
func (m *MeController) DeleteMe(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	if err := m.userService.DeleteAccount(ctx.Request.Context(), user.ID); err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// UpdateMyTheme		godoc
// @ID					updateMyTheme
// @Summary				Changes the look of the app for the logged user.
// @Param				Theme body transport.UpdateMyThemeRequest true "Theme"
// @Accept				application/json
// @Produce				application/json
// @Tags				Me
// @Success				204
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/me/theme [put]
func (m *MeController) UpdateMyTheme(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	request := transport.UpdateMyThemeRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondInvalidRequest(ctx, err)
		return
	}

	if err := m.userService.UpdateTheme(user.ID, request.Theme); err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// Welcome	 			godoc
// @ID					findMyWelcome
// @Summary				Get user stats for the Home Welcome
// @Description			Returns the user stats
// @Produce				application/json
// @Tags				Me
// @Success				200 {object} transport.WelcomeResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router				/me/welcome [get]
func (m *MeController) Welcome(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	numberOfFans, err := m.userFanService.CountFansOfUser(user.ID)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	numberOfScraps, err := m.scrapService.CountScrapsOfUser(user.ID)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	indirectConnections, err := m.friendshipService.CountIndirectConnections(user.ID)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	now := time.Now()
	profileViews, err := m.profileVisitService.FindProfileViews(user.ID, now)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	recentVisitors, err := m.profileVisitService.FindRecentVisitors(user.ID)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	hasRatedFriends, err := m.friendshipService.HasRatedFriends(user.ID)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	communityCount, err := m.communityMembershipService.CountUserMemberships(user.ID)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	resp := transport.WelcomeResponse{
		UserFirstName:       user.FirstName,
		UserCreatedAt:       user.CreatedAt,
		IndirectConnections: indirectConnections,
		DirectConnections:   user.FriendCount,
		NumberOfFans:        numberOfFans,
		NumberOfScraps:      numberOfScraps,
		ProfileViews:        *profileViews,
		RecentVisitors:      recentVisitors,
		TodaysLuck:          service.TodaysLuck(user.ID, now),
		HasRatedFriends:     hasRatedFriends,
		CommunityCount:      communityCount,
		RecordVisits:        user.RecordVisits,
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
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router			/me [put]
func (m *MeController) UpdateMe(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	request := transport.UpdateMeRequest{}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		respondInvalidRequest(ctx, err)
		return
	}

	var birthDate *time.Time
	if request.BirthDate != "" {
		parsed, err := time.ParseInLocation("2006-01-02", request.BirthDate, time.Local)
		if err != nil {
			respondInvalidRequest(ctx, err)
			return
		}
		birthDate = &parsed
	}

	err := m.userService.UpdateUser(user.ID, service.UpdateUserInput{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		BirthDate: birthDate,
	})
	if err != nil {
		RespondError(ctx, err)
		return
	}

	err = m.userProfileService.UpdateProfile(user.ID, service.UpdateUserProfileInput{
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
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// UpdateMyProfilePicture	godoc
// @ID					updateMyProfilePicture
// @Summary			Update the logged user's profile picture
// @Description		Uploads a PNG or JPEG (max 5 MB and 50 megapixels; stored shrunk to 1024px and 1 MB) as the logged user's profile picture.
// @Param			picture formData file true "Profile picture"
// @Accept			multipart/form-data
// @Produce			application/json
// @Tags			Me
// @Success			204
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router			/me/profile-picture [put]
func (m *MeController) UpdateMyProfilePicture(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	picture, ok := readPicture(ctx, false)
	if !ok {
		return
	}

	if err := m.userProfileService.UpdateProfilePicture(ctx.Request.Context(), user.ID, picture); err != nil {
		RespondError(ctx, err)
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
// @Failure			401 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router					/me/profile [get]
func (m *MeController) FindMyProfile(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	userWithProfile, err := m.userService.FindUserWithProfileById(user.ID, user.ID)

	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, userWithProfile)
}

// FindMyFriends	godoc
// @ID					findMyFriends
// @Summary			Gets a page of the logged user's friends, with their ratings and fan stars.
// @Param				page query int false "Page, starting at 1 (default 1)"
// @Param				pageSize query int false "Items per page (default 20, max 50)"
// @Description		Returns a page of friends, most recently seen first, each with the logged user's rating of them and whether they are a fan
// @Produce			application/json
// @Tags			Me
// @Success			200 {object} transport.FindMyFriendsResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router			/me/friends [get]
func (m *MeController) FindMyFriends(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	page, pageSize := pageQuery(ctx)
	friendships, err := m.friendshipService.FindMyFriends(user.ID, page, pageSize)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, friendships)
}

// FindMyBirthdays	godoc
// @ID					findMyBirthdays
// @Summary			Gets the logged user's friends with a birthday in the next 30 days.
// @Description		Returns the friends whose birthday is in the next 30 days (today included), soonest first
// @Produce			application/json
// @Tags			Me
// @Success			200 {object} transport.FindUpcomingBirthdaysResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router			/me/birthdays [get]
func (m *MeController) FindMyBirthdays(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	birthdays, err := m.friendshipService.FindUpcomingBirthdays(user.ID, time.Now())
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, birthdays)
}

// FindMyPendingFriendships	godoc
// @ID					findMyPendingFriendships
// @Summary					Gets a list of friendship invites of the user.
// @Description				Returns a list of users
// @Produce					application/json
// @Tags					Me
// @Success					200 {object} transport.FindPendingFriendshipOfUserResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router					/me/friends/pending [get]
func (m *MeController) FindMyPendingFriendships(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	friendships, err := m.friendshipService.FindPendingFriendshipsOfUser(user.ID)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, friendships)
}

// JoinCommunity	godoc
// @ID					joinCommunity
// @Summary			User joins a community.
// @Param			communityId path int true "Community ID"
// @Description		Joins a public community. For a private one, asks the owner to accept the user and answers pending true.
// @Produce			application/json
// @Tags			Me
// @Success			201 {object} transport.JoinCommunityResponse{}
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Failure			409 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router			/me/join/{communityId} [post]
func (m *MeController) JoinCommunity(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	communityId, ok := idParam(ctx, "communityId")
	if !ok {
		return
	}

	pending, err := m.communityMembershipService.JoinCommunity(user.ID, communityId)
	if err != nil {
		RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, transport.JoinCommunityResponse{Pending: pending})
}

// CancelJoinRequest	godoc
// @ID					cancelJoinRequest
// @Summary			User withdraws their request to join a private community.
// @Param			communityId path int true "Community ID"
// @Description		Deletes the logged user's pending join request (404 when there is none).
// @Produce			application/json
// @Tags			Me
// @Success			204
// @Failure			401 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router			/me/join-requests/{communityId} [delete]
func (m *MeController) CancelJoinRequest(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	communityId, ok := idParam(ctx, "communityId")
	if !ok {
		return
	}

	if err := m.communityMembershipService.CancelJoinRequest(user.ID, communityId); err != nil {
		RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

// LeaveCommunity	godoc
// @ID					leaveCommunity
// @Summary			User leaves a community.
// @Param			communityId path int true "Community ID"
// @Description		Leaves a community. The owner can't leave (403).
// @Produce			application/json
// @Tags			Me
// @Success			204
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Failure			403 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router			/me/leave/{communityId} [post]
func (m *MeController) LeaveCommunity(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	communityId, ok := idParam(ctx, "communityId")
	if !ok {
		return
	}

	if err := m.communityMembershipService.LeaveCommunity(user.ID, communityId); err != nil {
		RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

// FindMyMemberships		godoc
// @ID					findMyMemberships
// @Summary					Gets a page of the communities the logged user is a member of.
// @Param				page query int false "Page, starting at 1 (default 1)"
// @Param				pageSize query int false "Items per page (default 20, max 50)"
// @Param				random query bool false "Shuffled instead of by name, the same way all day (for the profile box)"
// @Description				Returns a page of communities
// @Produce					application/json
// @Tags					Me
// @Success					200 {object} transport.FindUserMembershipsResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router					/me/memberships [get]
func (m *MeController) FindMyMemberships(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	page, pageSize := pageQuery(ctx)
	random := ctx.Query("random") == "true"
	memberships, err := m.communityMembershipService.FindUserMemberships(user.ID, user.ID, page, pageSize, random)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, memberships)
}

// FindMyCommunities		godoc
// @ID					findMyCommunities
// @Summary					Gets a page of the communities the logged user created.
// @Param				page query int false "Page, starting at 1 (default 1)"
// @Param				pageSize query int false "Items per page (default 20, max 50)"
// @Description				Returns a page of communities, by name
// @Produce					application/json
// @Tags					Me
// @Success					200 {object} transport.FindMyOwnedCommunitiesResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router					/me/my-communities [get]
func (m *MeController) FindMyCommunities(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	page, pageSize := pageQuery(ctx)
	communities, err := m.communityService.FindCommunitiesCreatedByUser(user.ID, page, pageSize)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, communities)
}

// AmIMemberOf		godoc
// @ID					amIMemberOf
// @Summary			Evaluates if the logged user is member of specified community
// @Param			communityId path int true "Community ID"
// @Description		Says whether the logged user is a member, or is waiting for a private community's owner to accept them
// @Produce			application/json
// @Tags			Me
// @Success			200 {object} transport.AmIMemberOfResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router			/me/memberships/{communityId} [get]
func (m *MeController) AmIMemberOf(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	communityId, ok := idParam(ctx, "communityId")
	if !ok {
		return
	}

	status, err := m.communityMembershipService.FindMembershipStatus(communityId, user.ID)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, status)
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
// @Failure			401 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router			/me/friends/{userId} [get]
func (m *MeController) AmIFriendOf(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	targetUserId, ok := idParam(ctx, "userId")
	if !ok {
		return
	}
	if user.ID == targetUserId {
		ctx.JSON(http.StatusNoContent, nil)
		return
	}

	response, err := m.friendshipService.UsersFriendshipStatus(user.ID, targetUserId)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// FindMyCommunityInvitations	godoc
// @ID					findMyCommunityInvitations
// @Summary			Gets the logged user's pending community invitations.
// @Description		Returns the invitations friends sent the logged user to join their communities, newest first
// @Produce			application/json
// @Tags			Me
// @Success			200 {object} transport.FindMyCommunityInvitationsResponse{}
// @Failure			401 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router			/me/community-invitations [get]
func (m *MeController) FindMyCommunityInvitations(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)

	invitations, err := m.communityMembershipService.FindMyInvitations(user.ID)
	if err != nil {
		RespondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, invitations)
}

// AcceptCommunityInvitation	godoc
// @ID					acceptCommunityInvitation
// @Summary			User accepts an invitation to join a community.
// @Param			invitationId path int true "Invitation ID"
// @Description		Joins the community. In a private community, an invitation from a member other than the owner asks the owner instead, and answers pending true.
// @Produce			application/json
// @Tags			Me
// @Success			201 {object} transport.JoinCommunityResponse{}
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Failure			409 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router			/me/community-invitations/{invitationId}/accept [post]
func (m *MeController) AcceptCommunityInvitation(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	invitationId, ok := idParam(ctx, "invitationId")
	if !ok {
		return
	}

	pending, err := m.communityMembershipService.AcceptInvitation(user.ID, invitationId)
	if err != nil {
		RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, transport.JoinCommunityResponse{Pending: pending})
}

// DeclineCommunityInvitation	godoc
// @ID					declineCommunityInvitation
// @Summary			User declines an invitation to join a community.
// @Param			invitationId path int true "Invitation ID"
// @Description		Deletes the logged user's invitation (404 when there is none).
// @Produce			application/json
// @Tags			Me
// @Success			204
// @Failure			400 {object} transport.ErrorResponse
// @Failure			401 {object} transport.ErrorResponse
// @Failure			404 {object} transport.ErrorResponse
// @Security			BearerAuth
// @Router			/me/community-invitations/{invitationId} [delete]
func (m *MeController) DeclineCommunityInvitation(ctx *gin.Context) {
	user := ctx.MustGet("user").(*model.User)
	invitationId, ok := idParam(ctx, "invitationId")
	if !ok {
		return
	}

	if err := m.communityMembershipService.DeclineInvitation(user.ID, invitationId); err != nil {
		RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
