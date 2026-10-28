package seed

import (
	"fmt"
	"os"
	"strings"
	"time"
	"yorgut/internal/model"
	"yorgut/internal/storage"
	"yorgut/internal/transport"
)

// Extra mock data, so the paginated lists (friends, members, communities,
// forum, replies, testimonials) have more than one page in development.

// MainUserEmail is the developer's own account (SEED_MAIN_USER_EMAIL, so you
// can sign in as it): most mock users befriend it, and it owns enough
// communities to fill more than one page.
func MainUserEmail() string {
	if email := strings.TrimSpace(os.Getenv("SEED_MAIN_USER_EMAIL")); email != "" {
		return email
	}
	return "main@test.com"
}

// BigCommunityName is the community every user joins, so its member list
// pages.
const BigCommunityName = "Eu odeio acordar cedo"

// LongTopicTitle is a forum topic with enough replies to page.
const LongTopicTitle = "Apresente-se aqui! 👋"

var mockFirstNames = []string{
	"Ana", "Bruno", "Camila", "Daniel", "Eduarda", "Felipe", "Gabriela", "Henrique",
	"Isabela", "João", "Karina", "Leonardo", "Mariana", "Nicolas", "Olívia", "Pedro",
	"Quésia", "Rafael", "Sabrina", "Thiago", "Ursula", "Vinícius", "Wanessa", "Xande",
	"Yasmin", "Zeca",
}

var mockCommunityNames = []string{
	"Eu amo pão de queijo", "Odeio segunda-feira", "Eu já chorei vendo desenho",
	"Sou viciado em café", "Eu falo sozinho", "Tenho medo de palhaço",
	"Eu durmo de meia", "Amo cheiro de chuva", "Eu como a borda da pizza",
	"Sou do tempo do disquete", "Eu pulo no sofá", "Odeio acordar com despertador",
	"Eu canto no chuveiro", "Já caí da bicicleta", "Amo tirinhas de jornal",
	"Eu leio o fim do livro antes", "Tenho saudade do Orkut", "Sou fã de novela",
	"Eu juro que vou começar a dieta", "Amo dias nublados", "Eu ainda tenho MSN",
	"Quem nunca colou na prova?", "Eu como miojo cru", "Amo cachorro vira-lata",
	"Eu odeio pernilongo", "Já perdi o ônibus correndo", "Sou da turma do fundão",
	"Eu tenho mania de arrumação", "Amo festa junina", "Eu esqueço o nome das pessoas",
	"Tenho um caderno de receitas", "Sou noturno", "Eu adoro chuva de verão",
}

// mockCategories are all the categories: every one gets some communities.
var mockCategories = transport.CommunityCategories

// emailName drops the accents of the mock names, for their e-mails.
var emailName = strings.NewReplacer("á", "a", "ã", "a", "é", "e", "í", "i", "ó", "o", "ú", "u", "ç", "c")

// mockUsers returns extra users with generated pictures. Every fifth has a
// birthday in the coming weeks, for the birthdays box.
func mockUsers(store storage.ObjectStore) []*model.User {
	today := time.Now()
	users := make([]*model.User, len(mockFirstNames))
	for i, name := range mockFirstNames {
		birth := time.Date(1980+RandomInt(25), time.Month(1+RandomInt(12)), 1+RandomInt(28), 0, 0, 0, 0, time.Local)
		if i%5 == 0 {
			soon := today.AddDate(0, 0, i%30)
			birth = time.Date(1985+RandomInt(15), soon.Month(), soon.Day(), 0, 0, 0, 0, time.Local)
		}
		users[i] = &model.User{
			Email:     fmt.Sprintf("%s@test.com", emailName.Replace(strings.ToLower(name))),
			FirstName: name,
			LastName:  "Mock",
			BirthDate: &birth,
			Profile: &model.UserProfile{
				ProfilePicture: randomProfilePicture(store),
			},
		}
	}
	return users
}

// mockCommunities returns extra communities with generated pictures, without
// owners (GetCommunitiesSeed sets them).
func mockCommunities(store storage.ObjectStore) []*model.Community {
	communities := make([]*model.Community, len(mockCommunityNames))
	for i, name := range mockCommunityNames {
		communities[i] = &model.Community{
			Name:     name,
			Category: string(RandomChoice(mockCategories)),
			Profile: &model.CommunityProfile{
				ProfilePicture: randomCommunityPicture(store),
				Description:    ptr("Comunidade de teste."),
			},
		}
	}
	return communities
}

var testimonialTexts = []string{
	"Amigo para todas as horas! Sempre disposto a ajudar.",
	"Uma das pessoas mais engraçadas que eu conheço. Te adoro!",
	"Conheci na escola e até hoje somos amigos. Pessoa incrível!",
	"Super gente boa, inteligente e parceiro de todas as horas.",
	"Não tenho palavras para descrever essa pessoa maravilhosa.",
	"Melhor companhia para um café e uma boa conversa.",
	"Sempre com um sorriso no rosto. Contagia todo mundo!",
	"Amizade verdadeira é isso. Obrigado por tudo!",
	"Pessoa de coração enorme. Sorte de quem tem por perto.",
	"Parceiro de jogos, de rolês e de conversas até tarde.",
	"Sabe ouvir como ninguém. Um amigo de verdade.",
	"Divertido, leal e sempre presente. Te admiro muito!",
	"Escrevi esse depoimento só para aparecer no seu perfil :)",
	"Lembra daquela viagem? Inesquecível! Saudades.",
	"Você é demais! Nunca mude.",
}

// TestimonialText returns the i-th mock testimonial.
func TestimonialText(i int) string {
	return testimonialTexts[i%len(testimonialTexts)]
}

var scrapTexts = []string{
	"Oiii! Passando pra deixar um beijo!",
	"Saudades de você! Bora marcar alguma coisa?",
	"Feliz sexta-feira!! 🎉",
	"Viu o jogo ontem? Que sufoco!",
	"Me add no MSN!",
	"Adorei as fotos novas!",
	"Lembra de mim? Estudamos juntos na 5ª série!",
	"Bom dia flor do dia! 🌻",
	"Responde meu scrap, sumido(a)!",
	"Hahaha, aquela comunidade é muito boa!",
	"Parabéns pela conquista! Merecido demais.",
	"Passando pra dizer que você é demais!",
	"Vai na festa sábado?",
	"Já ouviu a música nova? Viciei!",
	"Obrigado pelo depoimento, amei!",
}

// ScrapText returns the i-th mock scrap.
func ScrapText(i int) string {
	return scrapTexts[i%len(scrapTexts)]
}
