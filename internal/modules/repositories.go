package modules

import (
	uR "app/internal/modules/user/repository"
	pR "app/internal/modules/pet/repository"
	sR "app/internal/modules/store/repository"
	"database/sql"
)

type Repository struct {
	User uR.UserRepositoryer
	Pet  pR.PetRepositoryer
	Store sR.StoreRepositoryer
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User: uR.NewUserRepository(db),
		Pet:  pR.NewPetRepository(db),
		Store: sR.NewStoreRepository(db),
	}
}
