package controller

import (
	"testing"
	"yorgut/internal/transport"

	"github.com/gin-gonic/gin/binding"
)

func TestCommunityCategoryBinding(t *testing.T) {
	valid := transport.UpdateCommunityRequest{Name: "Comunidade", Category: transport.CategoryMusic}
	if err := binding.Validator.ValidateStruct(valid); err != nil {
		t.Errorf("music: %v", err)
	}

	for _, category := range []transport.CommunityCategory{"", "Música", "musica"} {
		request := transport.UpdateCommunityRequest{Name: "Comunidade", Category: category}
		if err := binding.Validator.ValidateStruct(request); err == nil {
			t.Errorf("%q: got no error", category)
		}
	}
}

func TestNotBlankBinding(t *testing.T) {
	for _, text := range []string{"", " ", "\n\t  \n"} {
		request := transport.CreateScrapRequest{UserRecipientID: 1, Content: text}
		if err := binding.Validator.ValidateStruct(request); err == nil {
			t.Errorf("scrap %q: got no error", text)
		}
		topic := transport.CreateTopicRequest{CommunityID: 1, Title: text, Text: "texto"}
		if err := binding.Validator.ValidateStruct(topic); err == nil {
			t.Errorf("topic title %q: got no error", text)
		}
	}

	valid := []any{
		transport.CreateScrapRequest{UserRecipientID: 1, Content: " oi "},
		transport.CreateTestimonialRequest{UserRecipientID: 1, Content: "amigo"},
		transport.CreateTopicRequest{CommunityID: 1, Title: "título", Text: "texto"},
		transport.ReplyTopicRequest{Text: "resposta"},
	}
	for _, request := range valid {
		if err := binding.Validator.ValidateStruct(request); err != nil {
			t.Errorf("%+v: %v", request, err)
		}
	}
	reply := transport.ReplyTopicRequest{Text: "   "}
	if err := binding.Validator.ValidateStruct(reply); err == nil {
		t.Error("blank reply: got no error")
	}
}
