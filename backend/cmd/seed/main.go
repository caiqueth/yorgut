package main

import (
	"fmt"
	"log"
	"time"
	"yorgut/internal/database"
	"yorgut/internal/model"
	"yorgut/internal/repository"
)

func main() {
	db := database.DatabaseConnection()
	database.Migrate(db)

	userRepo := repository.NewUserRepository(db)
	communityRepo := repository.NewCommunityRepository(db)

	seedUsers(userRepo)
	seedFriendships(communityRepo, userRepo)

	fmt.Println("✅ Seeding complete!")
}

func seedUsers(r *repository.UserRepository) {

	users := []*model.User{
		model.NewUser("Caique", "Test", "caique@test.com", time.Date(1950, 12, 9, 0, 0, 0, 0, time.UTC)),
		model.NewUser("Lígia", "Test", "ligia@test.com", time.Date(1950, 1, 10, 0, 0, 0, 0, time.UTC)),
		model.NewUser("Aline", "Test", "aline@test.com", time.Date(1950, 10, 5, 0, 0, 0, 0, time.UTC)),
		model.NewUser("Renato", "Test", "renato@test.com", time.Date(1950, 3, 20, 0, 0, 0, 0, time.UTC)),
		model.NewUser("Diego", "Test", "diego@test.com", time.Date(1950, 10, 23, 0, 0, 0, 0, time.UTC)),
		model.NewUser("Luciana", "Test", "luciana@test.com", time.Date(1950, 5, 4, 0, 0, 0, 0, time.UTC)),
		model.NewUser("Weslley", "Test", "weslley@test.com", time.Date(1950, 10, 8, 0, 0, 0, 0, time.UTC)),
		model.NewUser("Carol", "Test", "carol@test.com", time.Date(1950, 12, 26, 0, 0, 0, 0, time.UTC)),
	}

	for _, user := range users {
		err := r.Create(user)
		if err != nil {
			log.Printf("❌ Failed to insert user %s: %v\n", user.FirstName, err)
		} else {
			fmt.Printf("✅ Inserted user: %s\n", user.FirstName)
		}

	}

}

func seedFriendships(r *repository.CommunityRepository, ur *repository.UserRepository) {
	userId := 1
	owner, err := ur.FindById(uint(userId))
	if err != nil {
		log.Printf("❌ Failed to find user %s: %v\n", 1, err)
	}

	communities := []*model.Community{
		model.NewCommunity("Deus me disse: desce e arrasa", "Descricao", "https://example.test.com", owner),
		model.NewCommunity("Eu nunca terminei uma borracha", "Descricao", "https://example.test.com", owner),
		model.NewCommunity("Tocava campainha e corria", "Descricao", "https://example.test.com", owner),
		model.NewCommunity("Só mais 5 minutinhos", "Descricao", "https://example.test.com", owner),
		model.NewCommunity("Não sei individualizar duplas", "Descricao", "https://example.test.com", owner),
		model.NewCommunity("Tive medo da Gina dos palitos", "Descricao", "https://example.test.com", owner),
		model.NewCommunity("Eu odeio acordar cedo", "Descricao", "https://example.test.com", owner),
		model.NewCommunity("Queria sorvete, mas era feijão", "Descricao", "https://example.test.com", owner),
		model.NewCommunity("Eu colho flores e o Tony Ramos", "Descricao", "https://example.test.com", owner),
	}

	for _, community := range communities {
		err := r.Create(community)
		if err != nil {
			log.Printf("❌ Failed to insert community %s: %v\n", community.Name, err)
		} else {
			fmt.Printf("✅ Inserted community: %s\n", community.Name)
		}

	}

}
