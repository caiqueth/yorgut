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
	"gorm.io/gorm"
)

func main() {
	time.Local, _ = time.LoadLocation("America/Sao_Paulo")

	db := localReset()
	database.Migrate(db)

	validate := validator.New()

	repositories := repository.NewRepositories(db)
	services := service.NewServices(repositories, validate)

	seedUsers(repositories.User)
	seedCommunities(repositories.Community, repositories.User)
	seedFriendships(repositories.Friendship, repositories.User, services.Friendship)
	seedOwnerMemberships(repositories.Community, services.CommunityMembership)
	seedMemberships(repositories.User, repositories.Community, services.CommunityMembership)
	seedTopics(repositories.Topic, repositories.Community, repositories.User)

	log.Printf("✅ Seeding complete!")
}

func clearStorageFolder() {
	if err := os.RemoveAll(path.Join("storage", "uploads")); err != nil {
		log.Fatalf("❌ could not clear storage/uploads: %v", err)
	}
}

// deleteSQLiteDb removes the SQLite file. On Windows this fails while another
// process (usually the running API) has it open; seeding on top of the old
// database would duplicate every row, so stop instead.
func deleteSQLiteDb() {
	dsn := os.Getenv("DATABASE_DSN")
	if err := os.Remove(dsn); err != nil && !os.IsNotExist(err) {
		log.Fatalf("❌ could not delete the database %q: %v. Stop the API (it keeps the file open) and run the seed again.", dsn, err)
	}
}

// localReset empties the database and returns a connection to it: SQLite
// (DEV_MODE=true) by deleting the file, Postgres by dropping the schema, which
// database.ResetLocalPostgres only allows on localhost. The database goes
// first: if it can't be reset, nothing (not even the uploaded images the old
// rows point to) has been touched yet.
func localReset() *gorm.DB {
	if os.Getenv("DEV_MODE") == "true" {
		deleteSQLiteDb()
	}

	db := database.DatabaseConnection()
	if db.Dialector.Name() == "postgres" {
		if err := database.ResetLocalPostgres(db); err != nil {
			log.Fatalf("❌ could not reset the database: %v", err)
		}
	}

	clearStorageFolder()
	return db
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

// seedOwnerMemberships makes every owner a member of their community, as
// creating a community through the API does.
func seedOwnerMemberships(
	cr *repository.CommunityRepository,
	ms *service.CommunityMembershipService,
) {
	communities, err := cr.FindMany()
	if err != nil {
		panic(err)
	}

	for _, community := range communities {
		if err := ms.JoinCommunity(community.OwnerID, community.ID); err != nil {
			log.Printf("❌ Failed to add owner to %s: %v\n", community.Name, err)
		} else {
			log.Printf("✅ Owner joined %s\n", community.Name)
		}
	}
}

func seedMemberships(
	ur *repository.UserRepository,
	cr *repository.CommunityRepository,
	ms *service.CommunityMembershipService,
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
			// owners already joined in seedOwnerMemberships
			if community.OwnerID == user.ID || !seed.RandomBoolean() {
				continue
			}
			err := ms.JoinCommunity(user.ID, community.ID)
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
		return
	}

	// Replies after their parents exist, since they point to their ids.
	replies := seed.GetTopicRepliesSeed(topics, users)
	if err := tr.CreateMany(replies); err != nil {
		log.Printf("❌ Failed to create topic replies => %v", err.Error())
	} else {
		log.Println("✅ Topics created successfully")
	}
}
