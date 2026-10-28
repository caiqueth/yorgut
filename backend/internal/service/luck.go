package service

import (
	"fmt"
	"hash/fnv"
	"time"
)

// fortunes are the home page's "sorte de hoje".
var fortunes = []string{
	"Boas notícias chegarão por e-mail.",
	"Um amigo antigo vai te mandar um recado.",
	"Hoje é um bom dia para entrar numa comunidade nova.",
	"Alguém está visitando seu perfil mais do que você imagina.",
	"Seu próximo depoimento vai te fazer sorrir.",
	"A paciência é a chave da sua felicidade hoje.",
	"Uma surpresa agradável está a caminho.",
	"Você vai reencontrar alguém que não vê há muito tempo.",
	"Seus amigos acham você mais legal do que você pensa.",
	"Hoje é dia de responder aqueles recados atrasados.",
	"Uma conversa inesperada vai mudar o seu dia.",
	"Sua criatividade está em alta. Aproveite!",
	"Não deixe para amanhã o recado que você pode mandar hoje.",
	"Um convite interessante vai aparecer em breve.",
	"Coisas boas acontecem para quem atualiza o perfil.",
	"A sorte sorri para quem é fã dos amigos.",
	"Hoje você vai aprender algo novo com alguém próximo.",
	"Seu bom humor vai contagiar quem estiver por perto.",
	"Uma foto antiga vai trazer boas lembranças.",
	"Grandes amizades começam com um simples scrap.",
}

// TodaysLuck picks the user's fortune for the day: the same all day, and
// usually different from one day to the next and between users.
func TodaysLuck(userId uint, now time.Time) string {
	hash := fnv.New32a()
	fmt.Fprintf(hash, "%d:%s", userId, dayOf(now).Format("2006-01-02"))
	return fortunes[hash.Sum32()%uint32(len(fortunes))]
}
