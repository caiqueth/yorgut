package seed

import (
	"math/rand"
	"path"
	"time"
	"yorgut/internal/helpers"
	"yorgut/internal/model"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func ptr[T any](v T) *T {
	return &v
}

func RandomChoice[T any](arr []T) T {
	if len(arr) == 0 {
		panic("RandomChoice called with empty slice")
	}
	return arr[rng.Intn(len(arr))]
}

func RandomBoolean() bool {
	return rng.Intn(2) == 1
}

func uploadFromSeed(seedPath string) *string {
	seedImagePath := path.Join("internal", "seed", seedPath)
	pictureAsBytes, _ := helpers.ReadLocalFile(seedImagePath)
	picturePath, _ := helpers.UploadPicture(pictureAsBytes)
	return picturePath
}

func GetUsersSeed() []*model.User {
	users := []*model.User{
		{
			Email:     "caiqueth@gmail.com",
			FirstName: "Caique",
			LastName:  "Test",
			BirthDate: time.Date(1950, 12, 9, 0, 0, 0, 0, time.Local),
			Profile: &model.UserProfile{
				ProfilePicture:     uploadFromSeed("profile_pictures/caique.jpg"),
				Description:        ptr("Descrição para testes\nPulo de linha para testes"),
				RelationshipStatus: ptr("casado (a)"),
				Interests:          ptr("Videogames"),
				Fashion:            ptr("Calça Baggy"),
				Website:            ptr("github.com/caiqueth"),
			},
		},
		{
			FirstName: "Lígia",
			LastName:  "Test",
			Email:     "ligiasalton@gmail.com",
			BirthDate: time.Date(1950, 1, 10, 0, 0, 0, 0, time.Local),
			Profile: &model.UserProfile{
				ProfilePicture: uploadFromSeed("profile_pictures/ligia.jpg"),
			},
		},
		{
			FirstName: "Aline",
			LastName:  "Test",
			Email:     "aline@test.com",
			BirthDate: time.Date(1950, 10, 5, 0, 0, 0, 0, time.Local),
			Profile: &model.UserProfile{
				ProfilePicture: uploadFromSeed("profile_pictures/aline.jpg"),
			},
		},
		{
			FirstName: "Renato",
			LastName:  "Test",
			Email:     "renato@test.com",
			BirthDate: time.Date(1950, 3, 20, 0, 0, 0, 0, time.Local),
			Profile: &model.UserProfile{
				ProfilePicture: uploadFromSeed("profile_pictures/renato.jpg"),
			},
		},
		{
			FirstName: "Diego",
			LastName:  "Test",
			Email:     "diego@test.com",
			BirthDate: time.Date(1950, 10, 23, 0, 0, 0, 0, time.Local),
			Profile: &model.UserProfile{
				ProfilePicture: uploadFromSeed("profile_pictures/diego.jpg"),
			},
		},
		{
			FirstName: "Luciana",
			LastName:  "Test",
			Email:     "luciana@test.com",
			BirthDate: time.Date(1950, 5, 4, 0, 0, 0, 0, time.Local),
			Profile: &model.UserProfile{
				ProfilePicture: uploadFromSeed("profile_pictures/luciana.jpg"),
			},
		},
		{
			FirstName: "Weslley",
			LastName:  "Test",
			Email:     "weslley@test.com",
			BirthDate: time.Date(1950, 10, 8, 0, 0, 0, 0, time.Local),
			Profile: &model.UserProfile{
				ProfilePicture: uploadFromSeed("profile_pictures/weslley.jpg"),
			},
		},
		{
			FirstName: "Carol",
			LastName:  "Test",
			Email:     "carol@test.com",
			BirthDate: time.Date(1950, 12, 26, 0, 0, 0, 0, time.Local),
			Profile: &model.UserProfile{
				ProfilePicture: uploadFromSeed("profile_pictures/carol.jpg"),
			},
		},
	}

	for _, user := range users {
		user.FriendCount = 0
	}

	return users
}

func GetCommunitiesSeed(owners []*model.User) []*model.Community {
	communities := []*model.Community{
		{
			Name:     "Deus me disse: desce e arrasa",
			Category: "Computadores e Internet",
			Profile: &model.CommunityProfile{
				ProfilePicture: uploadFromSeed("community_pictures/arrasa.png"),
				Description:    ptr("Descrição da comunidade para testes."),
			},
		},
		{
			Name:     "Eu nunca terminei uma borracha",
			Category: "Pessoas",
			Profile: &model.CommunityProfile{
				ProfilePicture: uploadFromSeed("community_pictures/borracha.png"),
				Description:    ptr("Descrição da comunidade para testes."),
			},
		},
		{
			Name:     "Tocava campainha e corria",
			Category: "Outros",
			Profile: &model.CommunityProfile{
				ProfilePicture: uploadFromSeed("community_pictures/campainha.png"),
				Description:    ptr("Descrição da comunidade para testes."),
			},
		},
		{
			Name:     "Só mais 5 minutinhos",
			Category: "Viagens",
			Profile: &model.CommunityProfile{
				ProfilePicture: uploadFromSeed("community_pictures/cinco.png"),
				Description:    ptr("Descrição da comunidade para testes."),
			},
		},
		{
			Name:     "Não sei individualizar duplas",
			Category: "Animais e bichos de estimacão",
			Profile: &model.CommunityProfile{
				ProfilePicture: uploadFromSeed("community_pictures/duplas.png"),
				Description:    ptr("Descrição da comunidade para testes."),
			},
		},
		{
			Name:     "Tive medo da Gina dos palitos",
			Category: "Artes e Entretenimento",
			Profile: &model.CommunityProfile{
				ProfilePicture: uploadFromSeed("community_pictures/gina.png"),
				Description:    ptr("Descrição da comunidade para testes."),
			},
		},
		{
			Name:     "Eu odeio acordar cedo",
			Category: "Ex-alunos e Escolas",
			Profile: &model.CommunityProfile{
				ProfilePicture: uploadFromSeed("community_pictures/odeio_acordar_cedo.png"),
				Description:    ptr("Descrição da comunidade para testes."),
			},
		},
		{
			Name:     "Queria sorvete, mas era feijão",
			Category: "Romances e Relacionamentos",
			Profile: &model.CommunityProfile{
				ProfilePicture: uploadFromSeed("community_pictures/sorvete.png"),
				Description:    ptr("Descrição da comunidade para testes."),
			},
		},
		{
			Name:     "Eu colho flores e o Tony Ramos",
			Category: "Saúde, Bem-estar e Fitness",
			Profile: &model.CommunityProfile{
				ProfilePicture: uploadFromSeed("community_pictures/tony.png"),
				Description:    ptr("Descrição da comunidade para testes."),
			},
		},
	}

	for _, community := range communities {
		owner := RandomChoice(owners)
		community.Owner = owner
		community.OwnerID = owner.ID
		community.CreatedBy = owner
		community.CreatedByID = owner.ID
		community.MemberCount = 0
	}
	return communities
}
