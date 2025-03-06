package dto

import (
	"time"
	"yorgut/internal/model"
)

type FriendshipDto struct {
	UserRequestedID uint
	UserRequesteeID uint
	CreatedAt       time.Time
}

type CreateFriendshipRequest struct {
	UserRequested uint
	UserRequestee uint
}

type FindFriendshipOfUserResponse struct {
	Friendships []FriendshipDto `json:"friendships"`
}

func (r *FindFriendshipOfUserResponse) FromFriendships(friendships []model.Friendship) {
	r.Friendships = make([]FriendshipDto, len(friendships))

	for i, f := range friendships {
		r.Friendships[i] = FriendshipDto{
			UserRequestedID: f.UserRequestedID,
			UserRequesteeID: f.UserRequesteeID,
			CreatedAt:       f.CreatedAt,
		}
	}
}
