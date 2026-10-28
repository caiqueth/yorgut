package service

import (
	"errors"
	"yorgut/internal/model"
	"yorgut/internal/repository"
	"yorgut/internal/transport"
)

var (
	ErrTestimonialNotFound  = newError(ErrNotFound, "O depoimento não existe.")
	ErrTestimonialForbidden = newError(ErrForbidden, "Este depoimento não é seu.")
	ErrTestimonialNotFriend = newError(ErrForbidden, "Só amigos podem escrever depoimentos.")
)

type TestimonialService struct {
	testimonialRepository repository.TestimonialRepository
	userRepository        repository.UserRepository
	friendshipRepository  repository.FriendshipRepository
}

func NewTestimonialService(
	testimonialRepository repository.TestimonialRepository,
	userRepository repository.UserRepository,
	friendshipRepository repository.FriendshipRepository,
) *TestimonialService {
	return &TestimonialService{
		testimonialRepository: testimonialRepository,
		userRepository:        userRepository,
		friendshipRepository:  friendshipRepository,
	}
}

// CreateTestimonial stores a pending testimonial; it only shows up publicly
// after the recipient approves it. Only friends write testimonials, never
// for themselves.
func (t *TestimonialService) CreateTestimonial(authorId uint, recipientId uint, content string) error {
	friends, err := areFriends(t.friendshipRepository, authorId, recipientId)
	if err != nil {
		return err
	}
	if !friends {
		return ErrTestimonialNotFriend
	}

	author, err := t.userRepository.FindById(authorId)
	if err != nil {
		return fromRepository(err, "o autor")
	}

	recipient, err := t.userRepository.FindById(recipientId)
	if err != nil {
		return fromRepository(err, "o destinatário")
	}

	testimonialModel := model.NewTestimonial(author, recipient, content)
	err = t.testimonialRepository.Create(testimonialModel)
	if err != nil {
		return err
	}

	return nil
}

// FindTestimonialsOfUser returns one page of the approved testimonials, newest
// first. Pending ones are never paged and only shown to the recipient.
func (t *TestimonialService) FindTestimonialsOfUser(viewerId uint, userId uint, page int, pageSize int) (*transport.FindTestimonialsOfUserResponse, error) {
	if err := requireAccess(t.userRepository, t.friendshipRepository, viewerId, userId, sectionTestimonials); err != nil {
		return nil, err
	}

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

	canWrite, err := areFriends(t.friendshipRepository, viewerId, userId)
	if err != nil {
		return nil, err
	}

	response := transport.FindTestimonialsOfUserResponse{}
	response.FromTestimonials(approved, pending, canWrite, page, pageSize, count)

	return &response, nil
}

// ApproveTestimonial is only allowed for the recipient.
func (t *TestimonialService) ApproveTestimonial(viewerId uint, testimonialId uint) error {
	testimonial, err := t.testimonialRepository.FindById(testimonialId)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrTestimonialNotFound
	}
	if err != nil {
		return err
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
	if errors.Is(err, repository.ErrNotFound) {
		return ErrTestimonialNotFound
	}
	if err != nil {
		return err
	}

	if testimonial.UserRecipientID != viewerId && testimonial.UserAuthorID != viewerId {
		return ErrTestimonialForbidden
	}

	return t.testimonialRepository.Delete(testimonialId)
}
