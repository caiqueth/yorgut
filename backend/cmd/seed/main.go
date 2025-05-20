package main

import (
	"log"
	"os"
	"path"
	"time"
	"yorgut/internal/database"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/seed"
	"yorgut/internal/service"

	"github.com/go-playground/validator/v10"
)

func main() {
	time.Local, _ = time.LoadLocation("America/Sao_Paulo")

	localReset()

	db := database.DatabaseConnection()
	database.Migrate(db)

	validate := validator.New()

	userRepository := repository.NewUserRepository(db)
	userProfileRepository := repository.NewUserProfileRepository(db)
	communityMembershipRepository := repository.NewCommunityMembershipRepository(db)
	communityRepository := repository.NewCommunityRepository(db)
	friendshipRepository := repository.NewFriendshipRepository(db)
	photoAlbumRepository := repository.NewPhotoAlbumRepository(db)
	topicRepository := repository.NewTopicRepository(db)
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

	frienshipService := service.NewFriendshipService(
		*friendshipRepository,
		*userRepository,
		validate,
	)

	seedUsers(userRepository)
	seedCommunities(communityRepository, userRepository)
	seedFriendships(friendshipRepository, userRepository, frienshipService)
	seedMemberships(userRepository, communityRepository, userService)
	seedTopics(topicRepository, communityRepository, userRepository)

	log.Printf("✅ Seeding complete!")
}

func clearStorageFolder() {
	os.RemoveAll(path.Join("storage", "uploads"))
}

func deleteDb() {
	os.Remove(os.Getenv("BLUEPRINT_DB_URL"))
}

func localReset() {
	clearStorageFolder()
	deleteDb()
}

func seedUsers(ur *repository.UserRepository) {
	users := seed.GetUsersSeed()

	if err := ur.CreateMany(users); err != nil {
		log.Printf("❌ Failed to insert users => %v\n", err.Error())
	} else {
		log.Printf("✅ Inserted all users\n")
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

	if err := cr.CreateMany(communities); err != nil {
		log.Printf("❌ Failed to insert communities => %v\n", err.Error())
	} else {
		log.Printf("✅ Inserted community\n")
	}

}

func seedFriendships(
	fr *repository.FriendshipRepository,
	ur *repository.UserRepository,
	fs *service.FriendshipService,
) {
	users, err := ur.FindMany()
	if err != nil {
		panic(err)
	}

	for _, user := range users {
		for _, friend := range users {

			if user.ID == friend.ID {
				continue
			}

			if fr.UsersAreFriends(user.ID, friend.ID) {
				continue
			}

			if fr.UsersHavePendingFriendshipInvites(user.ID, friend.ID) {
				continue
			}

			// creates frienship request
			friendship := model.NewFriendship(user, friend)
			err := fr.Create(friendship)

			if seed.RandomBoolean() {
				// ignore frienship request
				continue
			}

			// reply frienship request with random response
			fs.ReplyFriendship(friendship.ID, seed.RandomBoolean())

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

func seedTopics(
	tr *repository.TopicRepository,
	cr *repository.CommunityRepository,
	ur *repository.UserRepository,
) {
	communities, _ := cr.FindMany()
	users, _ := ur.FindMany()

	topics := seed.GetTopicsSeed(communities, users)

	if err := tr.CreateMany(topics); err != nil {
		log.Printf("❌ Failed to create topics => %v", err.Error())
	} else {
		log.Println("✅ Topics created successfully")
	}
}
