package seed

import (
	"errors"
	"fmt"
	"math/rand"
	"path"
	"strings"
	"time"
	"yorgut/internal/model"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomChoice[T any](arr []T) T {
	if len(arr) == 0 {
		panic("RandomChoice called with empty slice")
	}
	return arr[rng.Intn(len(arr))]
}

func RandomBoolean() bool {
	return rng.Intn(2) == 1
}

func GetUsersSeed() []*model.User {
	return []*model.User{
		model.NewUser(
			"Caique",
			"Test",
			"caique@test.com",
			time.Date(1950, 12, 9, 0, 0, 0, 0, time.Local),
		),
		model.NewUser(
			"Lígia",
			"Test",
			"ligia@test.com",
			time.Date(1950, 1, 10, 0, 0, 0, 0, time.Local),
		),
		model.NewUser(
			"Aline",
			"Test",
			"aline@test.com",
			time.Date(1950, 10, 5, 0, 0, 0, 0, time.Local),
		),
		model.NewUser(
			"Renato",
			"Test",
			"renato@test.com",
			time.Date(1950, 3, 20, 0, 0, 0, 0, time.Local),
		),
		model.NewUser(
			"Diego",
			"Test",
			"diego@test.com",
			time.Date(1950, 10, 23, 0, 0, 0, 0, time.Local),
		),
		model.NewUser(
			"Luciana",
			"Test",
			"luciana@test.com",
			time.Date(1950, 5, 4, 0, 0, 0, 0, time.Local),
		),
		model.NewUser(
			"Weslley",
			"Test",
			"weslley@test.com",
			time.Date(1950, 10, 8, 0, 0, 0, 0, time.Local),
		),
		model.NewUser(
			"Carol",
			"Test",
			"carol@test.com",
			time.Date(1950, 12, 26, 0, 0, 0, 0, time.Local),
		),
	}
}

type CommunitySeed struct {
	Name          string
	ImageSeedPath string
	Category      string
}

func GetCommunitiesSeed() []CommunitySeed {
	return []CommunitySeed{
		{
			Name:          "Deus me disse: desce e arrasa",
			ImageSeedPath: "arrasa.png",
			Category:      "Computadores e Internet",
		},
		{
			Name:          "Eu nunca terminei uma borracha",
			ImageSeedPath: "borracha.png",
			Category:      "Pessoas",
		},
		{
			Name:          "Tocava campainha e corria",
			ImageSeedPath: "campainha.png",
			Category:      "Outros",
		},
		{
			Name:          "Só mais 5 minutinhos",
			ImageSeedPath: "cinco.png",
			Category:      "Viagens",
		},
		{
			Name:          "Não sei individualizar duplas",
			ImageSeedPath: "duplas.png",
			Category:      "Animais e bichos de estimacão",
		},
		{
			Name:          "Tive medo da Gina dos palitos",
			ImageSeedPath: "gina.png",
			Category:      "Artes e Entretenimento",
		},
		{
			Name:          "Eu odeio acordar cedo",
			ImageSeedPath: "odeio_acordar_cedo.png",
			Category:      "Ex-alunos e Escolas",
		},
		{
			Name:          "Queria sorvete, mas era feijão",
			ImageSeedPath: "sorvete.png",
			Category:      "Romances e Relacionamentos",
		},
		{
			Name:          "Eu colho flores e o Tony Ramos",
			ImageSeedPath: "tony.png",
			Category:      "Saúde, Bem-estar e Fitness",
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
