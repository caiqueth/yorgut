package service

import (
	"yorgut/internal/model"
	"yorgut/internal/repository"
)

// What a profile keeps to friends: the user and their friends see everything;
// anyone else sees only what the user opened (model.User.Public*). A viewer id
// of 0 is "anyone": the user previewing their profile as a non-friend sees.

var ErrPrivateSection = newError(ErrForbidden, "Esta parte do perfil é visível só para amigos.")

// profileSection is a part of a profile its user may keep to friends.
type profileSection int

const (
	sectionAlbum profileSection = iota
	sectionCommunities
	sectionScraps
	sectionTestimonials
)

// profileAccess is what one viewer may see of one user's profile.
type profileAccess struct {
	owner *model.User
	// friend is true for the user themselves and their friends.
	friend bool
}

func (a profileAccess) sees(section profileSection) bool {
	if a.friend {
		return true
	}
	switch section {
	case sectionAlbum:
		return a.owner.PublicAlbum
	case sectionCommunities:
		return a.owner.PublicCommunities
	case sectionScraps:
		return a.owner.PublicScraps
	case sectionTestimonials:
		return a.owner.PublicTestimonials
	}
	return false
}

// accessTo loads ownerId's settings and whether viewerId is them or a friend.
func accessTo(users repository.UserRepository, friendships repository.FriendshipRepository, viewerId uint, ownerId uint) (profileAccess, error) {
	owner, err := users.FindById(ownerId)
	if err != nil {
		return profileAccess{}, fromRepository(err, "o usuário")
	}
	if viewerId != 0 && viewerId == ownerId {
		return profileAccess{owner: owner, friend: true}, nil
	}
	friends, err := areFriends(friendships, viewerId, ownerId)
	if err != nil {
		return profileAccess{}, err
	}
	return profileAccess{owner: owner, friend: friends}, nil
}

// requireAccess returns ErrPrivateSection when viewerId may not see section of
// ownerId's profile.
func requireAccess(users repository.UserRepository, friendships repository.FriendshipRepository, viewerId uint, ownerId uint, section profileSection) error {
	access, err := accessTo(users, friendships, viewerId, ownerId)
	if err != nil {
		return err
	}
	if !access.sees(section) {
		return ErrPrivateSection
	}
	return nil
}
