package seed

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"path"
	"strings"
	"time"
	"yorgut/internal/model"
)

func RandomChoice[T any](arr []T) T {
	max := len(arr)
	return arr[rand.IntN(max)]
}

func RandomBoolean() bool {
	return rand.IntN(2) == 1
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

type CommunitySeed struct {
	Name          string
	ImageSeedPath string
}

func GetCommunitiesSeed() []CommunitySeed {
	return []CommunitySeed{
		{
			Name:          "Deus me disse: desce e arrasa",
			ImageSeedPath: "arrasa.png",
		},
		{
			Name:          "Eu nunca terminei uma borracha",
			ImageSeedPath: "borracha.png",
		},
		{
			Name:          "Tocava campainha e corria",
			ImageSeedPath: "campainha.png",
		},
		{
			Name:          "Só mais 5 minutinhos",
			ImageSeedPath: "cinco.png",
		},
		{
			Name:          "Não sei individualizar duplas",
			ImageSeedPath: "duplas.png",
		},
		{
			Name:          "Tive medo da Gina dos palitos",
			ImageSeedPath: "gina.png",
		},
		{
			Name:          "Eu odeio acordar cedo",
			ImageSeedPath: "odeio_acordar_cedo.png",
		},
		{
			Name:          "Queria sorvete, mas era feijão",
			ImageSeedPath: "sorvete.png",
		},
		{
			Name:          "Eu colho flores e o Tony Ramos",
			ImageSeedPath: "tony.png",
		},
	}
}

func GetProfilePictureOf(user *model.User) string {
	return path.Join("internal", "seed", "profile_pictures", fmt.Sprintf("%s.jpg", strings.ToLower(user.FirstName)))
}

func GetCommunityPictureOf(community *model.Community) (string, error) {
	communitySeeds := GetCommunitiesSeed()
	for _, cs := range communitySeeds {
		if cs.Name == community.Name {
			return path.Join("internal", "seed", "community_pictures", cs.ImageSeedPath), nil
		}
	}
	return "", errors.New("could not find seed")
}
