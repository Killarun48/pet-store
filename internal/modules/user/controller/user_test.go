package controller

import (
	"app/internal/models"
	"context"
	"database/sql"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestGetByID(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		userService func(ctx context.Context, id string) (models.User, error)
		wantStatus  int
		wantBody    string
	}{
		{
			name:   "valid user ID",
			userID: "1",
			userService: func(ctx context.Context, id string) (models.User, error) {
				return models.User{ID: 1, UserName: "John Doe"}, nil
			},
			wantStatus: 200,
			wantBody:   "{\"id\":1,\"username\":\"John Doe\",\"firstName\":\"\",\"lastName\":\"\",\"email\":\"\",\"password\":\"\",\"phone\":\"\",\"userStatus\":0}\n",
		},
		{
			name:   "invalid user ID (sql.ErrNoRows)",
			userID: "2",
			userService: func(ctx context.Context, id string) (models.User, error) {
				return models.User{}, sql.ErrNoRows
			},
			wantStatus: 404,
			wantBody:   "{\"code\":404,\"type\":\"unknown\",\"message\":\"user not found\"}\n",
		},
		{
			name:   "error from userService.GetByID",
			userID: "3",
			userService: func(ctx context.Context, id string) (models.User, error) {
				return models.User{}, errors.New("some error")
			},
			wantStatus: 400,
			wantBody:   "{\"code\":400,\"type\":\"unknown\",\"message\":\"some error\"}\n",
		},
		{
			name:   "error from json.Marshal",
			userID: "4",
			userService: func(ctx context.Context, id string) (models.User, error) {
				return models.User{}, errors.New("json: error calling Marshal")
			},
			wantStatus: 400,
			wantBody:   "{\"code\":400,\"type\":\"unknown\",\"message\":\"json: error calling Marshal\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/users/"+tt.userID, nil)

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("userID", tt.userID)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

			//userService := &mockUserService{getByID: tt.userService}
			//responder := &mockResponder{}
			//responder := responder.NewResponder()
			/* uc := &UserController{userService: userService, responder: responder}

			uc.GetUserByName(w, r) */

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}

type mockUserService struct {
	getByID func(ctx context.Context, id string) (models.User, error)
	create  func(ctx context.Context, user models.User) error
	delete  func(ctx context.Context, id string) error
	list    func(ctx context.Context) ([]models.User, error)
	update  func(ctx context.Context, user models.User) error
}

func (m *mockUserService) GetByID(ctx context.Context, id string) (models.User, error) {
	return m.getByID(ctx, id)
}

func (m *mockUserService) Create(ctx context.Context, user models.User) error {
	return m.create(ctx, user)
}

func (m *mockUserService) Delete(ctx context.Context, id string) error {
	return m.delete(ctx, id)
}

func (m *mockUserService) List(ctx context.Context) ([]models.User, error) {
	return m.list(ctx)
}

func (m *mockUserService) Update(ctx context.Context, user models.User) error {
	return m.update(ctx, user)
}
