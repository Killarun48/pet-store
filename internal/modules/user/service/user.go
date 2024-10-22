package service

import (
	"app/internal/models"
	"context"
	"database/sql"
	"errors"
	"os"

	"github.com/go-chi/jwtauth"
)

type UserServicer interface {
	GetUserByName(ctx context.Context, userName string) (models.User, error)
	UpdateUser(ctx context.Context, userName string, user models.User) (id int, err error)
	DeleteUser(ctx context.Context, userName string) error
	CreateUser(ctx context.Context, user models.User) (id int, err error)
	LoginUser(ctx context.Context, userName string, password string) (token string, err error)
	LogoutUser(ctx context.Context) error
	CreateUsersWithArrayInput(ctx context.Context, users []models.User) error
	CreateUsersWithListInput(ctx context.Context, users []models.User) error
}

type UserRepositoryer interface {
	GetUserByName(ctx context.Context, userName string) (models.User, error)
	UpdateUser(ctx context.Context, userName string, user models.User) (id int, err error)
	DeleteUser(ctx context.Context, userName string) error
	CreateUser(ctx context.Context, user models.User) (id int, err error)
	CreateUsersWithArrayInput(ctx context.Context, users []models.User) error
	CreateUsersWithListInput(ctx context.Context, users []models.User) error
}

type UserService struct {
	userRepository UserRepositoryer
}

func NewUserService(userRepository UserRepositoryer) UserServicer {
	return &UserService{
		userRepository: userRepository,
	}
}

func (u *UserService) GetUserByName(ctx context.Context, userName string) (models.User, error) {
	return u.userRepository.GetUserByName(ctx, userName)
}

func (u *UserService) UpdateUser(ctx context.Context, userName string, user models.User) (int, error) {
	_, err := u.GetUserByName(ctx, userName)
	if err != nil {
		return 0, err
	}

	return u.userRepository.UpdateUser(ctx, userName, user)
}

func (u *UserService) DeleteUser(ctx context.Context, userName string) error {
	user, err := u.GetUserByName(ctx, userName)
	if err != nil {
		return err
	}

	if user.UserStatus == -1 {
		return errors.New("user not found, maybe user already deleted")
	}

	return u.userRepository.DeleteUser(ctx, userName)
}

func (s *UserService) CreateUser(ctx context.Context, user models.User) (int, error) {
	id, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.username" {
			return 0, errors.New("user already exists")
		}
	}

	return id, nil
}

func (s *UserService) LoginUser(ctx context.Context, userName string, password string) (token string, err error) {
	user, err := s.GetUserByName(ctx, userName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("user not found")
		}
		return "", err
	}

	if user.UserStatus == -1 {
		return "", errors.New("user not found, maybe user deleted")
	}

	if user.Password != password {
		return "", errors.New("invalid username/password supplied")
	}

	signKey := os.Getenv("SIGN_KEY")

	tokenAuth := jwtauth.New("HS256", []byte(signKey), nil)

	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"username": userName, "password": password})
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (u *UserService) LogoutUser(ctx context.Context) error {
	panic("implement me")
}

func (u *UserService) CreateUsersWithArrayInput(ctx context.Context, users []models.User) error {
	return u.userRepository.CreateUsersWithArrayInput(ctx, users)
}

func (u *UserService) CreateUsersWithListInput(ctx context.Context, users []models.User) error {
	return u.userRepository.CreateUsersWithListInput(ctx, users)
}
