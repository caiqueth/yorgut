package seed

import (
	"fmt"
	"math/rand"
	"time"
	"yorgut/internal/model"
	"yorgut/internal/storage"
	"yorgut/internal/transport"
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

func RandomInt(n int) int {
	return rng.Intn(n)
}

func GetUsersSeed(store storage.ObjectStore) []*model.User {
	users := []*model.User{
		{
			Email:     MainUserEmail(),
			FirstName: "Caique",
			LastName:  "Test",
			BirthDate: ptr(time.Date(1950, 12, 9, 0, 0, 0, 0, time.Local)),
			Profile: &model.UserProfile{
				ProfilePicture:     randomProfilePicture(store),
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
			Email:     "ligia@test.com",
			BirthDate: ptr(time.Date(1950, 1, 10, 0, 0, 0, 0, time.Local)),
			Profile: &model.UserProfile{
				ProfilePicture: randomProfilePicture(store),
			},
		},
		{
			FirstName: "Aline",
			LastName:  "Test",
			Email:     "aline@test.com",
			BirthDate: ptr(time.Date(1950, 10, 5, 0, 0, 0, 0, time.Local)),
			Profile: &model.UserProfile{
				ProfilePicture: randomProfilePicture(store),
			},
		},
		{
			FirstName: "Renato",
			LastName:  "Test",
			Email:     "renato@test.com",
			BirthDate: ptr(time.Date(1950, 3, 20, 0, 0, 0, 0, time.Local)),
			Profile: &model.UserProfile{
				ProfilePicture: randomProfilePicture(store),
			},
		},
		{
			FirstName: "Diego",
			LastName:  "Test",
			Email:     "diego@test.com",
			BirthDate: ptr(time.Date(1950, 10, 23, 0, 0, 0, 0, time.Local)),
			Profile: &model.UserProfile{
				ProfilePicture: randomProfilePicture(store),
			},
		},
		{
			FirstName: "Luciana",
			LastName:  "Test",
			Email:     "luciana@test.com",
			BirthDate: ptr(time.Date(1950, 5, 4, 0, 0, 0, 0, time.Local)),
			Profile: &model.UserProfile{
				ProfilePicture: randomProfilePicture(store),
			},
		},
		{
			FirstName: "Weslley",
			LastName:  "Test",
			Email:     "weslley@test.com",
			BirthDate: ptr(time.Date(1950, 10, 8, 0, 0, 0, 0, time.Local)),
			Profile: &model.UserProfile{
				ProfilePicture: randomProfilePicture(store),
			},
		},
		{
			FirstName: "Carol",
			LastName:  "Test",
			Email:     "carol@test.com",
			BirthDate: ptr(time.Date(1950, 12, 26, 0, 0, 0, 0, time.Local)),
			Profile: &model.UserProfile{
				ProfilePicture: randomProfilePicture(store),
			},
		},
	}

	users = append(users, mockUsers(store)...)

	for _, user := range users {
		user.FriendCount = 0
		// Seen sometime in the last week, so the friends and members boxes
		// have an order to show.
		user.LastSeenAt = ptr(time.Now().Add(-time.Duration(RandomInt(7*24*60)) * time.Minute))
	}

	return users
}

func GetCommunitiesSeed(owners []*model.User, store storage.ObjectStore) []*model.Community {
	communities := []*model.Community{
		{
			Name:     "Deus me disse: desce e arrasa",
			Category: string(transport.CategoryComputersInternet),
			Profile: &model.CommunityProfile{
				ProfilePicture: randomCommunityPicture(store),
				Description:    ptr("Descrição da comunidade para testes."),
			},
		},
		{
			Name:     "Eu nunca terminei uma borracha",
			Category: string(transport.CategoryIndividuals),
			Profile: &model.CommunityProfile{
				ProfilePicture: randomCommunityPicture(store),
				Description: ptr(`Você desde que é gente fica intigado com o desaparecimento das suas borrachas?
									Não entende por que de uma hora para outra sua borracha resolve passear?
									Seja porque você a perde logo no começo, ou porque ela some do nada, ou também às vezes porque seu amigo pega "emprestado"
									Agora você sabe que você não é o único no mundo que nunca terminou uma borracha!`),
			},
		},
		{
			Name:     "Tocava campainha e corria",
			Category: string(transport.CategoryOther),
			Profile: &model.CommunityProfile{
				ProfilePicture: randomCommunityPicture(store),
				Description:    ptr("Descrição da comunidade para testes."),
			},
		},
		{
			Name:     "Só mais 5 minutinhos",
			Category: string(transport.CategoryTravel),
			Profile: &model.CommunityProfile{
				ProfilePicture: randomCommunityPicture(store),
				Description:    ptr("Descrição da comunidade para testes."),
			},
		},
		{
			Name:     "Não sei individualizar duplas",
			Category: string(transport.CategoryPetsAnimals),
			Profile: &model.CommunityProfile{
				ProfilePicture: randomCommunityPicture(store),
				Description: ptr(`Imagino como se os integrantes da dupla estivessem sempre juntos, o tempo todo, como se eles fossem uma única pessoa.
									A propósito, você não acha que tanto o Chitãozinho e Xororó como o Chitãozinho e Xororó saíram bem na foto?`),
			},
		},
		{
			Name:     "Tive medo da Gina dos palitos",
			Private:  true,
			Category: string(transport.CategoryArtsEntertainment),
			Profile: &model.CommunityProfile{
				ProfilePicture: randomCommunityPicture(store),
				Description:    ptr("Descrição da comunidade para testes."),
			},
		},
		{
			Name:     "Eu odeio acordar cedo",
			Category: string(transport.CategoryAlumniSchools),
			Profile: &model.CommunityProfile{
				ProfilePicture: randomCommunityPicture(store),
				Description:    ptr("Descrição da comunidade para testes."),
			},
		},
		{
			Name:     "Queria sorvete, mas era feijão",
			Category: string(transport.CategoryRomanceRelationships),
			Profile: &model.CommunityProfile{
				ProfilePicture: randomCommunityPicture(store),
				Description: ptr(`sua mãe, tia ou avó, eh daquelas que bota o resto daquele "delicioso" feijão do mês passado nos potes de sorvete?
    								pois eh, se você já foi enganado naquela tarde quente de verão, foi louco achando q era sorvete, e abriu e viu aquela coisa nojenta e marrom com cheiro forte de feijão!
    								com certeza você já passou por essa situação, pode até ter sido outra gororoba, mas fatalmente você já caiu na pegadinha do pote de sorvete... algumas vezes, já te tinha separado o potinho, colher e cobertura pro sorvete!`),
			},
		},
		{
			Name:     "Eu colho flores e o Tony Ramos",
			Private:  true,
			Category: string(transport.CategoryHealthWellnessFitness),
			Profile: &model.CommunityProfile{
				ProfilePicture: randomCommunityPicture(store),
				Description:    ptr("EU COLHO FLORES E O TONY RAMOS KKKK!!!!!."),
			},
		},
	}

	// The main user owns the first 21 mock communities, one more than a page
	// of "minhas comunidades".
	var mainUser *model.User
	for _, owner := range owners {
		if owner.Email == MainUserEmail() {
			mainUser = owner
		}
	}
	mocks := mockCommunities(store)
	communities = append(communities, mocks...)

	for i, community := range communities {
		owner := RandomChoice(owners)
		if mock := i - (len(communities) - len(mocks)); mock >= 0 && mock < 21 && mainUser != nil {
			owner = mainUser
		}
		community.Owner = owner
		community.OwnerID = owner.ID
		community.CreatedBy = owner
		community.CreatedByID = owner.ID
		community.MemberCount = 0
	}
	return communities
}

func GetTopicsSeed(
	communities []*model.Community,
	users []*model.User,
) []*model.Topic {
	seedSize := 800
	topics := make([]*model.Topic, 0, seedSize)

	for i := range seedSize {
		community := RandomChoice(communities)
		author := RandomChoice(users)

		topic := &model.Topic{
			CommunityID: community.ID,
			Community:   community,
			AuthorID:    author.ID,
			Author:      author,
			Title:       fmt.Sprintf("Tópico de teste %v", i),
			Text:        "Esse eu Lorem enim mollit ipsum dolore. Consectetur aute officia incididunt occaecat nulla nostrud amet dolor adipisicing sint amet laboris est. Est non deserunt nisi ex.",
		}
		topics = append(topics, topic)
	}

	return topics
}

// GetTopicRepliesSeed returns replies to random topics, in the same
// community. The topics must already be saved: replies point to their ids,
// and MySQL enforces that foreign key.
func GetTopicRepliesSeed(
	topics []*model.Topic,
	users []*model.User,
) []*model.Topic {
	seedSize := len(topics)
	replies := make([]*model.Topic, 0, seedSize)

	for i := range seedSize {
		parent := RandomChoice(topics)
		author := RandomChoice(users)

		reply := &model.Topic{
			CommunityID:   parent.CommunityID,
			AuthorID:      author.ID,
			ParentTopicID: &parent.ID,
			Title:         fmt.Sprintf("Re: %s (%v)", parent.Title, i),
			Text:          "Concordo plenamente. Lorem ipsum dolor sit amet, consectetur adipisicing elit.",
		}
		replies = append(replies, reply)
	}

	return replies
}
