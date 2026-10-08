package service

import (
	"errors"
	"fmt"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/transport"

	"github.com/go-playground/validator/v10"
)

var (
	ErrTestimonialNotFound  = errors.New("testimonial not found")
	ErrTestimonialForbidden = errors.New("this testimonial does not belong to this user")
	ErrTestimonialNotFriend = errors.New("only friends can write testimonials for this user")
)

type TestimonialService struct {
	testimonialRepository repository.TestimonialRepository
	userRepository        repository.UserRepository
	friendshipRepository  repository.FriendshipRepository
	Validate              *validator.Validate
}

func NewTestimonialService(
	testimonialRepository repository.TestimonialRepository,
	userRepository repository.UserRepository,
	friendshipRepository repository.FriendshipRepository,
	validate *validator.Validate,
) *TestimonialService {
	return &TestimonialService{
		testimonialRepository: testimonialRepository,
		userRepository:        userRepository,
		friendshipRepository:  friendshipRepository,
		Validate:              validate,
	}
}

// CanWrite follows Orkut: only friends write testimonials, never for themselves.
func (t *TestimonialService) CanWrite(authorId uint, recipientId uint) bool {
	if authorId == recipientId {
		return false
	}
	return t.friendshipRepository.UsersAreFriends(authorId, recipientId)
}

// CreateTestimonial stores a pending testimonial; it only shows up publicly
// after the recipient approves it.
func (t *TestimonialService) CreateTestimonial(authorId uint, recipientId uint, content string) error {
	if !t.CanWrite(authorId, recipientId) {
		return ErrTestimonialNotFriend
	}

	author, err := t.userRepository.FindById(authorId)
	if err != nil {
		return fmt.Errorf("author not found %v", authorId)
	}

	recipient, err := t.userRepository.FindById(recipientId)
	if err != nil {
		return fmt.Errorf("recipient not found %v", recipientId)
	}

	testimonialModel := model.NewTestimonial(author, recipient, content)
	err = t.testimonialRepository.Create(testimonialModel)
	if err != nil {
		return err
	}

	return nil
}

// FindTestimonialsOfUser pages through the approved testimonials (page 0 means
// all of them). Pending ones are never paged and only shown to the recipient.
func (t *TestimonialService) FindTestimonialsOfUser(viewerId uint, userId uint, page int, pageSize int) (*transport.FindTestimonialsOfUserResponse, error) {
	approved, count, err := t.testimonialRepository.FindOfUserRecipient(userId, true, page, pageSize)
	if err != nil {
		return nil, err
	}

	pending := []model.Testimonial{}
	if viewerId == userId {
		pending, _, err = t.testimonialRepository.FindOfUserRecipient(userId, false, 0, 0)
		if err != nil {
			return nil, err
		}
	}

	response := transport.FindTestimonialsOfUserResponse{}
	response.FromTestimonials(approved, pending, t.CanWrite(viewerId, userId), page, pageSize, count)

	return &response, nil
}

// ApproveTestimonial is only allowed for the recipient.
func (t *TestimonialService) ApproveTestimonial(viewerId uint, testimonialId uint) error {
	testimonial, err := t.testimonialRepository.FindById(testimonialId)
	if err != nil {
		return ErrTestimonialNotFound
	}

	if testimonial.UserRecipientID != viewerId {
		return ErrTestimonialForbidden
	}

	return t.testimonialRepository.Approve(testimonialId)
}

// DeleteTestimonial lets the recipient reject or remove a testimonial, and the
// author take back one they wrote.
func (t *TestimonialService) DeleteTestimonial(viewerId uint, testimonialId uint) error {
	testimonial, err := t.testimonialRepository.FindById(testimonialId)
	if err != nil {
		return ErrTestimonialNotFound
	}

	if testimonial.UserRecipientID != viewerId && testimonial.UserAuthorID != viewerId {
		return ErrTestimonialForbidden
	}

	return t.testimonialRepository.Delete(testimonialId)
}
