package controller

import (
	"app/internal/infrastructure/responder"
	"app/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-chi/chi"
)

type UserControllerer interface {
	GetUserByName(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)
	LogoutUser(w http.ResponseWriter, r *http.Request)
	CreateUsersWithArrayInput(w http.ResponseWriter, r *http.Request)
	CreateUsersWithListInput(w http.ResponseWriter, r *http.Request)
}

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

type UserController struct {
	userService UserServicer
	responder   responder.Responder
}

func NewUserController(userService UserServicer, respond responder.Responder) UserControllerer {
	return &UserController{
		userService: userService,
		responder:   respond,
	}
}

//	@id			2getUserByName
//	@Summary	Get user by user name
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		username	path		string	true	"The name that needs to be fetched. Use admin for testing."
//	@Success	200			{object}	models.User
//	@Router		/user/{username} [get]
func (uc UserController) GetUserByName(w http.ResponseWriter, r *http.Request) {
	userName := chi.URLParam(r, "username")

	user, err := uc.userService.GetUserByName(context.Background(), userName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			uc.responder.ErrorNotFound(w, errors.New("user not found"))
			return
		}

		uc.responder.ErrorBadRequest(w, err)
		return
	}

	jsonResp, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		uc.responder.ErrorBadRequest(w, err)
		return
	}

	fmt.Fprintln(w, string(jsonResp))
}

//	@id				3updateUser
//	@Summary		Updated user
//	@Description	This can only be done by the logged in user.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			username	path		string		true	"name that need to be updated"
//	@Param			object		body		models.User	true	"Updated user object"
//	@Success		200			{object}	responder.Response
//	@Router			/user/{username} [put]
func (uc UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userName := chi.URLParam(r, "username")

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		uc.responder.ErrorBadRequest(w, err)
		return
	}

	id, err := uc.userService.UpdateUser(context.Background(), userName, user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			uc.responder.ErrorNotFound(w, errors.New("user not found"))
			return
		}

		uc.responder.ErrorBadRequest(w, err)
		return
	}
	// TODO вернуть ID

	uc.responder.Success(w, fmt.Sprint(id))
}

//	@id				4deleteUser
//	@Summary		Delete user
//	@Description	This can only be done by the logged in user.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			username	path		string	true	"The name that needs to be deleted"
//	@Success		200			{object}	responder.Response
//	@Router			/user/{username} [delete]
func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userName := chi.URLParam(r, "username")

	err := uc.userService.DeleteUser(context.Background(), userName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			uc.responder.ErrorNotFound(w, errors.New("user not found"))
			return
		}

		uc.responder.ErrorBadRequest(w, err)
		return
	}

	uc.responder.Success(w, userName)
}

//	@id				8createUser
//	@Summary		Create user
//	@Description	This can only be done by the logged in user.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			object	body		models.User	true	"Created user object"
//	@Success		200		{object}	responder.Response
//	@Router			/user [post]
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		uc.responder.ErrorBadRequest(w, err)
		return
	}

	id, err := uc.userService.CreateUser(context.Background(), user)
	if err != nil {
		uc.responder.ErrorBadRequest(w, err)
		return
	}

	uc.responder.Success(w, fmt.Sprint(id))
}

//	@id			5loginUser
//	@Summary	Logs user into the system
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		username	query		string	true	"The user name for login" Default(admin)
//	@Param		password	query		string	true	"The password for login in clear text" Default(admin)
//	@Success	200			{object}	responder.Response
//	@Header		200			{string}	X-Expires-After	"date in UTC when token expires"
//	@Header		200			{int}		X-Rate-Limit	"calls per hour allowed by user"
//	@Router		/user/login [get]
func (uc UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	token, err := uc.userService.LoginUser(context.Background(), userName, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			uc.responder.ErrorNotFound(w, errors.New("invalid username/password supplied"))
			return
		}

		uc.responder.ErrorBadRequest(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		Secure:   false,
	})

	w.Header().Set("X-Expires-After", time.Now().Add(1*time.Hour).Format("Mon Jan 2 15:04:05 UTC 2006"))
	w.Header().Set("X-Rate-Limit", "5000")

	session := gofakeit.IntRange(1234567891234, 1934567891234)
	uc.responder.Success(w, fmt.Sprintf("logged in user session:%d", session))
}

//	@id			6logoutUser
//	@Summary	Logs out current logged in user session
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	responder.Response
//	@Router		/user/logout [get]
func (uc UserController) LogoutUser(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   false,
	})

	uc.responder.Success(w, "ok")
}

//	@id			7createUsersWithArrayInput
//	@Summary	Creates list of users with given input array
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		object	body		[]models.User	true	"List of user object"
//	@Success	200		{object}	responder.Response
//	@Router		/user/createWithArray [post]
func (uc UserController) CreateUsersWithArrayInput(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	err := json.NewDecoder(r.Body).Decode(&users)
	if err != nil {
		uc.responder.ErrorBadRequest(w, err)
		return
	}

	err = uc.userService.CreateUsersWithArrayInput(context.Background(), users)
	if err != nil {
		uc.responder.ErrorBadRequest(w, err)
		return
	}

	uc.responder.Success(w, "ok")
}

//	@id			1createUsersWithListInput
//	@Summary	Creates list of users with given input array
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		object	body		[]models.User	true	"List of user object"
//	@Success	200		{object}	responder.Response
//	@Router		/user/createWithList [post]
func (uc UserController) CreateUsersWithListInput(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	err := json.NewDecoder(r.Body).Decode(&users)
	if err != nil {
		uc.responder.ErrorBadRequest(w, err)
		return
	}

	err = uc.userService.CreateUsersWithListInput(context.Background(), users)
	if err != nil {
		uc.responder.ErrorBadRequest(w, err)
		return
	}

	uc.responder.Success(w, "ok")
}
