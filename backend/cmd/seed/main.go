package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"yorgut/internal/database"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/seed"
	"yorgut/internal/service"
	"yorgut/internal/storage"

	"gorm.io/gorm"
)

func main() {
	time.Local, _ = time.LoadLocation("America/Sao_Paulo")

	store, err := storage.NewFromEnv(context.Background())
	if err != nil {
		log.Fatalf("❌ could not set up the uploads storage: %v", err)
	}

	db := localReset(store)
	database.Migrate(db)

	repositories := repository.NewRepositories(db)
	services := service.NewServices(repositories, store)

	seedUsers(repositories.User, store)
	seedCommunities(repositories.Community, repositories.User, store)
	seedFriendships(repositories.Friendship, repositories.User, services.Friendship)
	seedFriendsOfMainUser(repositories.Friendship, repositories.User, services.Friendship)
	seedOwnerMemberships(repositories.Community, services.CommunityMembership)
	seedMemberships(repositories.User, repositories.Community, services.CommunityMembership)
	seedTopics(repositories.Topic, repositories.Community, repositories.User)
	seedLongTopic(repositories.Topic, repositories.Community, repositories.User)
	seedTestimonials(repositories.Testimonial, repositories.Friendship, repositories.User)
	seedScraps(repositories.Scrap, repositories.User)
	seedFans(repositories.UserFan, repositories.Friendship, repositories.User)
	seedProfileVisits(repositories.ProfileVisit, repositories.User)

	log.Printf("✅ Seeding complete!")
}

// clearable is implemented by both stores; S3 only clears a bucket on
// localhost.
type clearable interface {
	CanClear() error
	Clear(context.Context) error
}

// uploadsClearer checks, before anything is reset, that the uploads can be
// cleared.
func uploadsClearer(store storage.ObjectStore) clearable {
	clearer, ok := store.(clearable)
	if !ok {
		log.Fatalf("❌ %v can't be cleared", store)
	}
	if err := clearer.CanClear(); err != nil {
		log.Fatalf("❌ %v", err)
	}
	return clearer
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
// (DEV_MODE=true) by deleting the file, MySQL by dropping its tables, which
// database.ResetLocalMySQL only allows on localhost. It checks first that
// the uploads can be cleared too (an S3 bucket only on localhost). The
// database goes first: if it can't be reset, nothing (not even the uploaded
// images the old rows point to) has been touched yet.
func localReset(store storage.ObjectStore) *gorm.DB {
	uploads := uploadsClearer(store)

	if os.Getenv("DEV_MODE") == "true" {
		deleteSQLiteDb()
	}

	db := database.DatabaseConnection()
	if db.Dialector.Name() == "mysql" {
		if err := database.ResetLocalMySQL(db); err != nil {
			log.Fatalf("❌ could not reset the database: %v", err)
		}
	}

	if err := uploads.Clear(context.Background()); err != nil {
		log.Fatalf("❌ could not clear %v: %v", store, err)
	}
	return db
}

func seedUsers(ur *repository.UserRepository, store storage.ObjectStore) {
	users := seed.GetUsersSeed(store)

	if err := ur.CreateMany(users); err != nil {
		log.Printf("❌ Failed to insert users => %v\n", err.Error())
	} else {
		log.Printf("✅ Inserted all users\n")
	}
}

func seedCommunities(
	cr *repository.CommunityRepository,
	ur *repository.UserRepository,
	store storage.ObjectStore,
) {
	owners, err := ur.FindMany()
	if err != nil {
		panic(err)
	}

	communities := seed.GetCommunitiesSeed(owners, store)

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

			// Skip pairs that are already friends or have a request either way.
			existing, err := fr.FindFriendshipsBetweenUsers(user.ID, friend.ID)
			if err != nil || len(existing) > 0 {
				continue
			}

			friendship := model.NewFriendship(user, friend)
			if err := fr.Create(friendship); err != nil {
				log.Printf("❌ Failed to create friendship between %s and %s: %v\n", user.FirstName, friend.FirstName, err)
				continue
			}
			log.Printf("✅ Friendship requested between %s and %s\n", user.FirstName, friend.FirstName)

			// Leave half of the requests pending; answer the rest at random.
			if seed.RandomBoolean() {
				continue
			}
			if err := fs.ReplyFriendship(friend.ID, friendship.ID, seed.RandomBoolean()); err != nil {
				log.Printf("❌ Failed to answer the request between %s and %s: %v\n", user.FirstName, friend.FirstName, err)
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
		if _, err := ms.JoinCommunity(community.OwnerID, community.ID); err != nil {
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
			if community.OwnerID == user.ID {
				continue
			}
			// Everyone joins the big community. The main user joins fewer of
			// the others, so "descobrir comunidades" has more than a page.
			joinChance := 50
			if community.Name == seed.BigCommunityName {
				joinChance = 100
			} else if user.Email == seed.MainUserEmail() {
				joinChance = 30
			}
			if seed.RandomInt(100) >= joinChance {
				continue
			}
			pending, err := ms.JoinCommunity(user.ID, community.ID)
			if err != nil {
				log.Printf("❌ Failed to create %s membership for user %s: %v\n", community.Name, user.FirstName, err)
			} else if pending {
				log.Printf("✅ %s asked to join %s (private)\n", user.FirstName, community.Name)
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

// seedProfileVisits makes random users visit each other's profiles over the
// last two weeks, for the home page's views and recent visitors.
func seedProfileVisits(pr *repository.ProfileVisitRepository, ur *repository.UserRepository) {
	users, err := ur.FindMany()
	if err != nil {
		panic(err)
	}

	count := 0
	for _, visited := range users {
		for _, visitor := range users {
			if visitor.ID == visited.ID {
				continue
			}
			for daysAgo := 0; daysAgo < 14; daysAgo++ {
				if seed.RandomInt(4) != 0 {
					continue
				}
				at := time.Now().AddDate(0, 0, -daysAgo).Add(-time.Duration(seed.RandomInt(12*60)) * time.Minute)
				day := time.Date(at.Year(), at.Month(), at.Day(), 0, 0, 0, 0, time.Local)
				if err := pr.Record(visitor.ID, visited.ID, day, at); err != nil {
					log.Printf("❌ Failed to create a profile visit: %v", err)
					continue
				}
				count++
			}
		}
	}
	log.Printf("✅ %d profile visits created", count)
}

func findMainUser(ur *repository.UserRepository) *model.User {
	user, err := ur.FindByEmail(seed.MainUserEmail())
	if err != nil {
		panic(err)
	}
	return user
}

// seedFriendsOfMainUser makes most users friends with the main user, so their
// friend list pages; every sixth is left as a pending request to them.
func seedFriendsOfMainUser(
	fr *repository.FriendshipRepository,
	ur *repository.UserRepository,
	fs *service.FriendshipService,
) {
	main := findMainUser(ur)
	users, err := ur.FindMany()
	if err != nil {
		panic(err)
	}

	friends, pending := 0, 0
	for i, user := range users {
		if user.ID == main.ID {
			continue
		}
		if friends, _ := fr.UsersAreFriends(user.ID, main.ID); friends {
			continue
		}
		// The random friendships above may have left a request either way.
		if pending, _ := fr.HasPendingBetweenUsers(user.ID, main.ID); !pending {
			if err := fs.CreateFriendship(user.ID, main.ID); err != nil {
				log.Printf("❌ Failed to ask %s to be friends: %v", user.FirstName, err)
				continue
			}
		}
		list, err := fr.FindFriendshipsBetweenUsers(user.ID, main.ID)
		if err != nil || len(list) == 0 {
			log.Printf("❌ Failed to find the friendship request of %s: %v", user.FirstName, err)
			continue
		}
		request := list[len(list)-1]
		if i%6 == 0 && request.UserRequesteeID == main.ID {
			pending++
			continue
		}
		if err := fs.ReplyFriendship(request.UserRequesteeID, request.ID, true); err != nil {
			log.Printf("❌ Failed to accept %s: %v", user.FirstName, err)
			continue
		}
		friends++
	}
	log.Printf("✅ %s got %d new friends and %d pending requests", main.FirstName, friends, pending)
}

// seedLongTopic adds a topic with enough replies to page, in the community
// everyone joins.
func seedLongTopic(
	tr *repository.TopicRepository,
	cr *repository.CommunityRepository,
	ur *repository.UserRepository,
) {
	main := findMainUser(ur)
	users, err := ur.FindMany()
	if err != nil {
		panic(err)
	}
	communities, err := cr.FindMany()
	if err != nil {
		panic(err)
	}
	var community *model.Community
	for _, c := range communities {
		if c.Name == seed.BigCommunityName {
			community = c
		}
	}
	if community == nil {
		log.Printf("❌ No community named %q for the long topic", seed.BigCommunityName)
		return
	}

	topic := model.NewTopic(main.ID, community.ID, nil, seed.LongTopicTitle, "Conta um pouco sobre você: de onde é, o que faz e por que odeia acordar cedo.")
	if err := tr.Create(topic); err != nil {
		log.Printf("❌ Failed to create the long topic: %v", err)
		return
	}

	replies := make([]*model.Topic, 0, 35)
	for i := 1; i <= 35; i++ {
		author := seed.RandomChoice(users)
		text := fmt.Sprintf("Oi, eu sou %s! Resposta número %d. Também odeio acordar cedo.", author.FirstName, i)
		replies = append(replies, model.NewTopic(author.ID, community.ID, &topic.ID, "", text))
	}
	if err := tr.CreateMany(replies); err != nil {
		log.Printf("❌ Failed to create the long topic's replies: %v", err)
		return
	}
	log.Printf("✅ Created %q with %d replies", seed.LongTopicTitle, len(replies))
}

// seedTestimonials gives the main user 12 approved testimonials (more than a
// page of the profile panel) and 3 waiting for their approval, from friends.
func seedTestimonials(
	tr *repository.TestimonialRepository,
	fr *repository.FriendshipRepository,
	ur *repository.UserRepository,
) {
	main := findMainUser(ur)
	friends, err := fr.FindAllFriendsOfUser(main.ID)
	if err != nil {
		panic(err)
	}

	approved, pending := 0, 0
	for i := range friends {
		if approved+pending == 15 {
			break
		}
		testimonial := model.NewTestimonial(&friends[i], main, seed.TestimonialText(i))
		testimonial.Approved = approved < 12
		if err := tr.Create(testimonial); err != nil {
			log.Printf("❌ Failed to create a testimonial from %s: %v", friends[i].FirstName, err)
			continue
		}
		if testimonial.Approved {
			approved++
		} else {
			pending++
		}
	}
	log.Printf("✅ %s got %d approved and %d pending testimonials", main.FirstName, approved, pending)
}

// seedScraps gives the main user 35 scraps (4 pages of their scrapbook) and
// everyone else a few, from random users, spread over the last two months.
func seedScraps(sr *repository.ScrapRepository, ur *repository.UserRepository) {
	main := findMainUser(ur)
	users, err := ur.FindMany()
	if err != nil {
		panic(err)
	}

	scraps := []*model.Scrap{}
	add := func(recipient *model.User, n int) {
		for i := 0; i < n; i++ {
			author := seed.RandomChoice(users)
			for author.ID == recipient.ID {
				author = seed.RandomChoice(users)
			}
			scrap := model.NewScrap(author, recipient, seed.ScrapText(seed.RandomInt(100)))
			scrap.CreatedAt = time.Now().Add(-time.Duration(seed.RandomInt(60*24*60)) * time.Minute)
			scraps = append(scraps, scrap)
		}
	}
	for _, user := range users {
		if user.ID == main.ID {
			add(user, 35)
		} else {
			add(user, 1+seed.RandomInt(4))
		}
	}

	if err := sr.CreateMany(scraps); err != nil {
		log.Printf("❌ Failed to create scraps: %v", err)
		return
	}
	log.Printf("✅ Created %d scraps", len(scraps))
}

// seedFans makes most of the main user's friends their fans (more than a page
// of "fãs"), and gives everyone else a few fans among their friends.
func seedFans(ufr *repository.UserFanRepository, fr *repository.FriendshipRepository, ur *repository.UserRepository) {
	main := findMainUser(ur)
	users, err := ur.FindMany()
	if err != nil {
		panic(err)
	}

	count := 0
	for _, user := range users {
		friends, err := fr.FindAllFriendsOfUser(user.ID)
		if err != nil {
			panic(err)
		}
		chance := 30
		if user.ID == main.ID {
			chance = 85
		}
		for i := range friends {
			if seed.RandomInt(100) >= chance {
				continue
			}
			if err := ufr.Add(model.NewUserFan(friends[i].ID, user.ID)); err != nil {
				log.Printf("❌ Failed to make %s a fan of %s: %v", friends[i].FirstName, user.FirstName, err)
				continue
			}
			count++
		}
	}
	log.Printf("✅ Created %d fans", count)
}
