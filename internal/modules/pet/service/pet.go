package service

import (
	"app/internal/models"
	"context"
	"mime/multipart"
)

type PetServicer interface {
	UploadFile(ctx context.Context, id int, file multipart.File, fileHeader *multipart.FileHeader) error
	AddPet(ctx context.Context, body models.Pet) (models.Pet, error)
	UpdatePet(ctx context.Context, body models.Pet) (models.Pet, error)
	FindPetsByStatus(ctx context.Context, status []string) ([]models.Pet, error)
	FindPetsByTags(ctx context.Context, tags []string) ([]models.Pet, error)
	GetPetById(ctx context.Context, id int) (models.Pet, error)
	UpdatePetWithForm(ctx context.Context, id int, name string, status string) error
	DeletePet(ctx context.Context, id int) error
}

type PetRepositoryer interface {
	UploadFile(ctx context.Context, id int, file multipart.File, fileHeader *multipart.FileHeader) error
	AddPet(ctx context.Context, body models.Pet) (models.Pet, error)
	UpdatePet(ctx context.Context, body models.Pet) (models.Pet, error)
	FindPetsByStatus(ctx context.Context, status []string) ([]models.Pet, error)
	FindPetsByTags(ctx context.Context, tags []string) ([]models.Pet, error)
	GetPetById(ctx context.Context, id int) (models.Pet, error)
	UpdatePetWithForm(ctx context.Context, id int, name string, status string) error
	DeletePet(ctx context.Context, id int) error
}

type PetService struct {
	petRepository PetRepositoryer
}

func NewPetService(petRepository PetRepositoryer) PetServicer {
	return &PetService{
		petRepository: petRepository,
	}
}

func (s *PetService) UploadFile(ctx context.Context, id int, file multipart.File, fileHeader *multipart.FileHeader) error {
	return s.petRepository.UploadFile(ctx, id, file, fileHeader)
}

func (s *PetService) AddPet(ctx context.Context, pet models.Pet) (models.Pet, error) {
	return s.petRepository.AddPet(ctx, pet)
}

func (s *PetService) UpdatePet(ctx context.Context, pet models.Pet) (models.Pet, error) {
	_, err := s.GetPetById(ctx, pet.ID)
	if err != nil {
		return models.Pet{}, err
	}

	return s.petRepository.UpdatePet(ctx, pet)
}

func (s *PetService) FindPetsByStatus(ctx context.Context, status []string) ([]models.Pet, error) {
	return s.petRepository.FindPetsByStatus(ctx, status)
}

func (s *PetService) FindPetsByTags(ctx context.Context, tags []string) ([]models.Pet, error) {
	return s.petRepository.FindPetsByTags(ctx, tags)
}

func (s *PetService) GetPetById(ctx context.Context, id int) (models.Pet, error) {
	return s.petRepository.GetPetById(ctx, id)
}

func (s *PetService) UpdatePetWithForm(ctx context.Context, id int, name string, status string) error {
	_, err := s.GetPetById(ctx, id)
	if err != nil {
		return err
	}
	
	return s.petRepository.UpdatePetWithForm(ctx, id, name, status)
}

func (s *PetService) DeletePet(ctx context.Context, id int) error {
	_, err := s.GetPetById(ctx, id)
	if err != nil {
		return err
	}
	
	return s.petRepository.DeletePet(ctx, id)
}
