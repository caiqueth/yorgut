package main

import (
	"log"
	"os"
	"path"
	"yorgut/internal/database"
	"yorgut/internal/helpers"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/seed"
	"yorgut/internal/service"

	"github.com/go-playground/validator/v10"
)

func main() {
	clearStorageFolder()
	deleteDb()

	db := database.DatabaseConnection()
	database.Migrate(db)

	validate := validator.New()

	userRepository := repository.NewUserRepository(db)
	userProfileRepository := repository.NewUserProfileRepository(db)
	communityMembershipRepository := repository.NewCommunityMembershipRepository(db)
	communityRepository := repository.NewCommunityRepository(db)
	communityProfileRepository := repository.NewCommunityProfileRepository(db)
	friendshipRepository := repository.NewFriendshipRepository(db)
	photoAlbumRepository := repository.NewPhotoAlbumRepository(db)
	// photoRepository := repository.NewPhotoRepository(db)
	// scrapRepository := repository.NewScrapRepository(db)
	// testimonialRepository := repository.NewTestimonialRepository(db)

	userService := service.NewUserService(
		*userRepository,
		*userProfileRepository,
		*communityRepository,
		*communityMembershipRepository,
		*photoAlbumRepository,
		*friendshipRepository,
		validate,
	)

	seedUsers(userRepository)
	seedCommunities(communityRepository, userRepository)
	seedFriendships(friendshipRepository, userRepository)
	seedMemberships(userRepository, communityRepository, userService)
	seedUserProfiles(userProfileRepository, userRepository)
	seedCommunityProfiles(communityProfileRepository, communityRepository)

	log.Printf("✅ Seeding complete!")
}

func clearStorageFolder() {
	os.RemoveAll(path.Join("storage", "uploads"))
}

func deleteDb() {
	os.Remove("test.db")
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

	communities := seed.GetCommunitiesSeed()

	for _, community := range communities {
		err := cr.Create(model.NewCommunity(
			community.Name,
			community.Category,
			seed.RandomChoice(owners),
		))
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
	ur *repository.UserRepository,
	cr *repository.CommunityRepository,
	us *service.UserService,
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
			if !seed.RandomBoolean() {
				continue
			}
			err := us.JoinCommunity(user.ID, community.ID)
			if err != nil {
				log.Printf("❌ Failed to create %s membership for user %s: %v\n", community.Name, user.FirstName, err)
			} else {
				log.Printf("✅ %s Membership created for user %s\n", community.Name, user.FirstName)
			}
		}
	}
}

func seedUserProfiles(
	upr *repository.UserProfileRepository,
	ur *repository.UserRepository,
) {
	users, err := ur.FindMany()
	if err != nil {
		panic(err)
	}

	for _, user := range users {
		profile := model.NewUserProfile(user)

		seedImagePath := seed.GetProfilePictureOf(user)
		pictureAsBytes, err := helpers.ReadLocalFile(seedImagePath)
		if err != nil {
			log.Printf("❌ Failed to read profile picture: %s\n", seedImagePath)
			continue
		}

		profilePicturePath, err := helpers.UploadProfilePicture(pictureAsBytes)
		if err != nil {
			log.Printf("❌ Failed to save profile picture to storage: %s. Error: %v\n", seedImagePath, err.Error())
			continue
		}

		model.AddUserProfilePicture(profile, *profilePicturePath)
		err = upr.Create(profile)

		if err != nil {
			log.Printf("❌ Failed to create profile for user %s: %v\n", user.FirstName, err)
		} else {
			log.Printf("✅ Profile created for user %s\n", user.FirstName)
		}
	}
}

func seedCommunityProfiles(
	cpr *repository.CommunityProfileRepository,
	cr *repository.CommunityRepository,
) {
	communities, err := cr.FindMany()
	if err != nil {
		panic(err)
	}

	for _, community := range communities {
		profile := model.NewCommunityProfile(community)

		seedImagePath, err := seed.GetCommunityPictureOf(community)
		if err != nil {
			log.Printf("❌ Could not find a matching picture file for the community: %s\n", community.Name)
			continue
		}

		pictureAsBytes, err := helpers.ReadLocalFile(seedImagePath)
		if err != nil {
			log.Printf("❌ Failed to read community picture: %s\n", seedImagePath)
			continue
		}

		communityPicturePath, err := helpers.UploadCommunityPicture(pictureAsBytes)
		if err != nil {
			log.Printf("❌ Failed to save community picture to storage: %s. Error: %v\n", seedImagePath, err.Error())
			continue
		}

		model.AddCommunityProfilePicture(profile, *communityPicturePath)
		err = cpr.Create(profile)

		if err != nil {
			log.Printf("❌ Failed to create profile for community %s: %v\n", community.Name, err)
		} else {
			log.Printf("✅ Profile created for community %s\n", community.Name)
		}
	}

}
