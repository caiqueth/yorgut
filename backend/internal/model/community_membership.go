package model

import "gorm.io/gorm"

type CommunityMembership struct {
	gorm.Model
	CommunityID uint
	Community   *Community `gorm:"foreignKey:CommunityID;references:ID"`
	UserID      uint
	User        *User `gorm:"foreignKey:UserID;references:ID"`
}

func NewCommunityMembership(community *Community, user *User) *CommunityMembership {
	return &CommunityMembership{
		CommunityID: community.ID,
		Community:   community,
		UserID:      user.ID,
		User:        user,
	}
}
