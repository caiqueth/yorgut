package seed

import (
	"math/rand/v2"
	"time"
	"yorgut/internal/model"
)

func chooseRandomOwner(owners []*model.User) *model.User {
	return owners[rand.IntN(len(owners))]
}

func GetUsersSeed() []*model.User {
	return []*model.User{
		model.NewUser(
			"Caique",
			"Test",
			"caique@test.com",
			time.Date(1950, 12, 9, 0, 0, 0, 0, time.UTC),
		),
		model.NewUser(
			"Lígia",
			"Test",
			"ligia@test.com",
			time.Date(1950, 1, 10, 0, 0, 0, 0, time.UTC),
		),
		model.NewUser(
			"Aline",
			"Test",
			"aline@test.com",
			time.Date(1950, 10, 5, 0, 0, 0, 0, time.UTC),
		),
		model.NewUser(
			"Renato",
			"Test",
			"renato@test.com",
			time.Date(1950, 3, 20, 0, 0, 0, 0, time.UTC),
		),
		model.NewUser(
			"Diego",
			"Test",
			"diego@test.com",
			time.Date(1950, 10, 23, 0, 0, 0, 0, time.UTC),
		),
		model.NewUser(
			"Luciana",
			"Test",
			"luciana@test.com",
			time.Date(1950, 5, 4, 0, 0, 0, 0, time.UTC),
		),
		model.NewUser(
			"Weslley",
			"Test",
			"weslley@test.com",
			time.Date(1950, 10, 8, 0, 0, 0, 0, time.UTC),
		),
		model.NewUser(
			"Carol",
			"Test",
			"carol@test.com",
			time.Date(1950, 12, 26, 0, 0, 0, 0, time.UTC),
		),
	}
}

func GetCommunitiesSeed(owners []*model.User) []*model.Community {
	return []*model.Community{
		model.NewCommunity(
			"Deus me disse: desce e arrasa",
			"Descricao",
			"/storage/uploads/community_pictures/arrasa.png",
			chooseRandomOwner(owners),
		),
		model.NewCommunity(
			"Eu nunca terminei uma borracha",
			"Descricao",
			"/storage/uploads/community_pictures/borracha.png",
			chooseRandomOwner(owners),
		),
		model.NewCommunity(
			"Tocava campainha e corria",
			"Descricao",
			"/storage/uploads/community_pictures/campainha.png",
			chooseRandomOwner(owners),
		),
		model.NewCommunity(
			"Só mais 5 minutinhos",
			"Descricao",
			"/storage/uploads/community_pictures/cinco.png",
			chooseRandomOwner(owners),
		),
		model.NewCommunity(
			"Não sei individualizar duplas",
			"Descricao",
			"/storage/uploads/community_pictures/duplas.png",
			chooseRandomOwner(owners),
		),
		model.NewCommunity(
			"Tive medo da Gina dos palitos",
			"Descricao",
			"/storage/uploads/community_pictures/gina.png",
			chooseRandomOwner(owners),
		),
		model.NewCommunity(
			"Eu odeio acordar cedo",
			"Descricao",
			"/storage/uploads/community_pictures/odeio_acordar_cedo.png",
			chooseRandomOwner(owners),
		),
		model.NewCommunity(
			"Queria sorvete, mas era feijão",
			"Descricao",
			"/storage/uploads/community_pictures/sorvete.png",
			chooseRandomOwner(owners),
		),
		model.NewCommunity(
			"Eu colho flores e o Tony Ramos",
			"Descricao",
			"/storage/uploads/community_pictures/tony.png",
			chooseRandomOwner(owners),
		),
	}
}
