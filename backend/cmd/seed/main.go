package main

import (
	"fmt"
	"log"
	"yorgut/internal/database"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/seed"
)

func main() {
	db := database.DatabaseConnection()
	database.Migrate(db)

	userRepo := repository.NewUserRepository(db)
	communityRepo := repository.NewCommunityRepository(db)
	communityMembershipRepo := repository.NewCommunityMembershipRepository(db)
	friendshipRepo := repository.NewFriendshipRepository(db)

	seedUsers(userRepo)
	seedCommunities(communityRepo, userRepo)
	seedFriendships(friendshipRepo, userRepo)
	seedMemberships(communityMembershipRepo, userRepo, communityRepo)

	fmt.Println("✅ Seeding complete!")
}

func seedUsers(ur *repository.UserRepository) {

	users := seed.GetUsersSeed()

	for _, user := range users {
		err := ur.Create(user)
		if err != nil {
			log.Printf("❌ Failed to insert user %s: %v\n", user.FirstName, err)
		} else {
			log.Printf("✅ Inserted user: %s\n", user.FirstName)
		}

	}

}

func seedCommunities(
	cr *repository.CommunityRepository,
	ur *repository.UserRepository,
) {
	owners, err := ur.FindMany()
	if err != nil {
		panic(err)
	}

	communities := seed.GetCommunitiesSeed(owners)

	for _, community := range communities {
		err := cr.Create(community)
		if err != nil {
			log.Printf("❌ Failed to insert community %s: %v\n", community.Name, err)
		} else {
			log.Printf("✅ Inserted community: %s\n", community.Name)
		}

	}

}

func seedFriendships(
	fr *repository.FriendshipRepository,
	ur *repository.UserRepository,
) {
	users, err := ur.FindMany()
	if err != nil {
		panic(err)
	}

	for _, user := range users {
		for _, friend := range users {

			if user.ID == friend.ID || user.ID > friend.ID {
				continue
			}

			if fr.UsersAreFriends(user.ID, friend.ID) {
				continue
			}

			friendship := model.NewFriendship(user, friend)
			friendship.Accepted = seed.RandomBoolean()

			err := fr.Create(friendship)
			if err != nil {
				log.Printf("❌ Failed to create friendship between %s and %s: %v\n", user.FirstName, friend.FirstName, err)
			} else {
				log.Printf("✅ Friendship created between %s and %s\n", user.FirstName, friend.FirstName)
			}
		}

	}

}

func seedMemberships(
	cmr *repository.CommunityMembershipRepository,
	ur *repository.UserRepository,
	cr *repository.CommunityRepository,
) {
	users, err := ur.FindMany()
	if err != nil {
		panic(err)
	}

	communities, err := cr.FindMany()
	if err != nil {
		panic(err)
	}

	for _, user := range users {
		for _, community := range communities {
			if seed.RandomBoolean() {
				membership := model.NewCommunityMembership(community, user)
				err := cmr.Create(membership)
				if err != nil {
					log.Printf("❌ Failed to create %s membership for user %s: %v\n", community.Name, user.FirstName, err)
				} else {
					log.Printf("✅ %s Membership created for user %s\n", community.Name, user.FirstName)
				}
			}
		}
	}
}
