package service

import (
	"errors"
	"path"
	"yorgut/internal/dto"
	"yorgut/internal/helpers"
	"yorgut/internal/model"
	"yorgut/internal/repository"

	"github.com/go-playground/validator/v10"
)

type CommunityService struct {
	communityRepository repository.CommunityRepository
	userRepository      repository.UserRepository
	Validate            *validator.Validate
}

func (c *CommunityService) Create(request dto.CreateCommunityRequest) error {
	owner, err := c.userRepository.FindById(request.CreatedByID)
	if err != nil {
		return errors.New("could not find specified community owner")
	}

	picture, err := helpers.FileToBytes(request.Picture)
	if err != nil {
		return err
	}

	imageHasValidExtension := helpers.CheckValidImage(picture)
	if !imageHasValidExtension {
		return errors.New("invalid image format")
	}

	storedFileName, err := helpers.WriteFileToLocalStorage(path.Join(helpers.GetUploadsPath(), "community_pictures"), picture)
	if err != nil {
		return err
	}

	communityModel := model.NewCommunity(request.Name, request.Description, storedFileName, owner)
	err = c.communityRepository.Create(communityModel)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommunityService) Delete(id uint) error {
	err := c.communityRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommunityService) FindById(id uint) (dto.FindCommunityByIdResponse, error) {
	community, err := c.communityRepository.FindById(id)
	if err != nil {
		return dto.FindCommunityByIdResponse{}, errors.New("could not find specified community id")
	}

	r := dto.FindCommunityByIdResponse{
		Name:        community.Name,
		Description: community.Description,
		CreatedBy:   community.CreatedByID,
	}
	return r, nil
}

func (c *CommunityService) FindMany() (*dto.Response, error) {
	communities, err := c.communityRepository.FindMany()
	if err != nil {
		return nil, err
	}

	communityDtos := make([]dto.CommunityDto, len(communities))

	for i, community := range communities {
		communityDtos[i] = dto.CommunityDto{
			Name:        community.Name,
			Description: community.Description,
			CreatedBy:   community.CreatedByID,
			PictureUrl:  community.PictureUrl,
		}
	}

	return &dto.Response{
		Data: map[string][]dto.CommunityDto{"communities": communityDtos},
	}, nil
}

func (c *CommunityService) Update(user dto.UpdateCommunityRequest) (dto.Response, error) {
	panic("unimplemented")
}

func NewCommunityService(
	communityRepository repository.CommunityRepository,
	userRepository repository.UserRepository,
	validate *validator.Validate,
) *CommunityService {
	return &CommunityService{
		communityRepository: communityRepository,
		userRepository:      userRepository,
		Validate:            validate,
	}
}
