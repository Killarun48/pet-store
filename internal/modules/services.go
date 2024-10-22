package modules

import (
	uS "app/internal/modules/user/service"
	pS "app/internal/modules/pet/service"
	sS "app/internal/modules/store/service"
)

type Service struct {
	User uS.UserServicer
	Pet  pS.PetServicer
	Store sS.StoreServicer
}

func NewService(repos *Repository) *Service {
	return &Service{
		User: uS.NewUserService(repos.User),
		Pet:  pS.NewPetService(repos.Pet),
		Store: sS.NewStoreService(repos.Store),
	}
}