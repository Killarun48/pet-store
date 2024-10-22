package main

import (
	"app/internal/infrastructure/responder"
	"app/internal/models"
	"app/internal/modules/user/controller"
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockUserService struct {
	getUserByName             func(ctx context.Context, userName string) (models.User, error)
	updateUser                func(ctx context.Context, userName string, user models.User) (id int, err error)
	deleteUser                func(ctx context.Context, userName string) error
	createUser                func(ctx context.Context, user models.User) (id int, err error)
	loginUser                 func(ctx context.Context, userName string, password string) (token string, err error)
	logoutUser                func(ctx context.Context) error
	createUsersWithArrayInput func(ctx context.Context, users []models.User) error
	createUsersWithListInput  func(ctx context.Context, users []models.User) error
}

func (m *mockUserService) GetUserByName(ctx context.Context, userName string) (models.User, error) {
	return m.getUserByName(ctx, userName)
}

func (m *mockUserService) UpdateUser(ctx context.Context, userName string, user models.User) (id int, err error) {
	return m.updateUser(ctx, userName, user)
}

func (m *mockUserService) DeleteUser(ctx context.Context, userName string) error {
	return m.deleteUser(ctx, userName)
}

func (m *mockUserService) CreateUser(ctx context.Context, user models.User) (id int, err error) {
	return m.createUser(ctx, user)
}

func (m *mockUserService) LoginUser(ctx context.Context, userName string, password string) (token string, err error) {
	return m.loginUser(ctx, userName, password)
}

func (m *mockUserService) LogoutUser(ctx context.Context) error {
	return m.logoutUser(ctx)
}

func (m *mockUserService) CreateUsersWithArrayInput(ctx context.Context, users []models.User) error {
	return m.createUsersWithArrayInput(ctx, users)
}

func (m *mockUserService) CreateUsersWithListInput(ctx context.Context, users []models.User) error {
	return m.createUsersWithListInput(ctx, users)
}

func TestUser(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		userService func(ctx context.Context, userName string) (models.User, error)
		wantStatus  int
		wantBody    string
	}{
		{
			name: "valid body",
			body: `{"name": "Rodney William Whitaker","birthDate": "1931-06-12"}`,
			userService: func(ctx context.Context, userName string) (models.User, error) {
				return models.User{ID: 1, UserName: "Rodney William Whitaker"}, nil
			},
			wantStatus: 200,
			wantBody:   "{\"code\":200,\"success\":true,\"message\":\"автор создан\"}\n",
		},
		{
			name: "error create author",
			body: `{"name": "Rodney William Whitaker","birthDate": "1931-06-12"}`,
			userService: func(ctx context.Context, userName string) (models.User, error) {
				return models.User{}, fmt.Errorf("error")
			},
			wantStatus: 400,
			wantBody:   "{\"code\":400,\"success\":false,\"message\":\"error\"}\n",
		},
		{
			name: "invalid body",
			body: `{"name" "Rodney William Whitaker"}`,
			userService: func(ctx context.Context, userName string) (models.User, error) {
				return models.User{}, nil
			},
			wantStatus: 200,
			wantBody:   "{\"code\":400,\"success\":false,\"message\":\"invalid character '\\\"' after object key\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api", strings.NewReader(tt.body))
			r.Header.Set("Content-Type", "application/json")

			userService := &mockUserService{
				getUserByName: tt.userService,
				//createAuthor: tt.authorService,
			}
			responder := responder.NewResponder()
			ac := controller.NewUserController(userService, responder)
			ac.GetUserByName(w, r)

			assert.Equal(t, tt.wantStatus, w.Code)
			//assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}

func TestServer(t *testing.T) {
	pathDB = "./test.db"
	s := NewServer(":8088")

	go s.Serve()

	s.Stop()
	time.Sleep(2 * time.Second)

	os.Remove("./test.db")
}

func TestUser2(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		userService func(ctx context.Context, userName string) (models.User, error)
		wantStatus  int
		wantBody    string
	}{
		{
			name: "valid body",
			body: `{"name": "Rodney William Whitaker","birthDate": "1931-06-12"}`,
			userService: func(ctx context.Context, userName string) (models.User, error) {
				return models.User{ID: 1, UserName: "Rodney William Whitaker"}, nil
			},
			wantStatus: 200,
			wantBody:   "{\"code\":200,\"success\":true,\"message\":\"автор создан\"}\n",
		},
		{
			name: "error create author",
			body: `{"name": "Rodney William Whitaker","birthDate": "1931-06-12"}`,
			userService: func(ctx context.Context, userName string) (models.User, error) {
				return models.User{}, fmt.Errorf("error")
			},
			wantStatus: 400,
			wantBody:   "{\"code\":400,\"success\":false,\"message\":\"error\"}\n",
		},
		{
			name: "invalid body",
			body: `{"name" "Rodney William Whitaker"}`,
			userService: func(ctx context.Context, userName string) (models.User, error) {
				return models.User{}, nil
			},
			wantStatus: 200,
			wantBody:   "{\"code\":400,\"success\":false,\"message\":\"invalid character '\\\"' after object key\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api", strings.NewReader(tt.body))
			r.Header.Set("Content-Type", "application/json")

			userService := &mockUserService{
				getUserByName: tt.userService,
				//createAuthor: tt.authorService,
			}
			responder := responder.NewResponder()
			ac := controller.NewUserController(userService, responder)
			ac.GetUserByName(w, r)

			assert.Equal(t, tt.wantStatus, w.Code)
			//assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}