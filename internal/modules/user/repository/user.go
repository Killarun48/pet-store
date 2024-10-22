package repository

import (
	"app/internal/models"
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

const (
	usersTable = "users"
)

type UserRepositoryer interface {
	GetUserByName(ctx context.Context, userName string) (models.User, error)
	UpdateUser(ctx context.Context, userName string, user models.User) (id int, err error)
	DeleteUser(ctx context.Context, userName string) error
	CreateUser(ctx context.Context, user models.User) (id int, err error)
	CreateUsersWithArrayInput(ctx context.Context, users []models.User) error
	CreateUsersWithListInput(ctx context.Context, users []models.User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepositoryer {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUserByName(ctx context.Context, userName string) (models.User, error) {
	var user models.User

	type userRow struct {
		ID         sql.NullInt64
		UserName   sql.NullString
		FirstName  sql.NullString
		LastName   sql.NullString
		Email      sql.NullString
		Password   sql.NullString
		Phone      sql.NullString
		UserStatus sql.NullInt64
	}

	var u userRow

	err := sq.Select(
		"id",
		"username",
		"first_name",
		"last_name",
		"email",
		"password",
		"phone",
		"user_status",
	).
		From(usersTable).
		/* Where(sq.And{
			sq.Eq{"username": userName},
			sq.NotEq{"user_status": -1},
		}). */
		Where(sq.Eq{"username": userName}).
		RunWith(r.db).
		QueryRowContext(ctx).
		Scan(
			&u.ID,
			&u.UserName,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.Password,
			&u.Phone,
			&u.UserStatus,
		)
	if err != nil {
		return models.User{}, err
	}

	user.ID = int(u.ID.Int64)
	user.UserName = u.UserName.String
	user.FirstName = u.FirstName.String
	user.LastName = u.LastName.String
	user.Email = u.Email.String
	user.Password = u.Password.String
	user.Phone = u.Phone.String
	user.UserStatus = int(u.UserStatus.Int64)

	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, userName string, user models.User) (int, error) {
	res, err := sq.Update(usersTable).
		SetMap(map[string]interface{}{
			"username":    user.UserName,
			"first_name":  user.FirstName,
			"last_name":   user.LastName,
			"email":       user.Email,
			"phone":       user.Phone,
			"user_status": user.UserStatus,
		}).
		Where(sq.Eq{"username": userName}).
		RunWith(r.db).
		ExecContext(ctx)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *userRepository) DeleteUser(ctx context.Context, userName string) error {
	_, err := sq.Update(usersTable).
		SetMap(map[string]interface{}{
			"user_status": -1,
		}).
		Where(sq.Eq{"username": userName}).
		RunWith(r.db).
		ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) CreateUser(ctx context.Context, user models.User) (int, error) {
	res, err := sq.Insert(usersTable).
		Columns(
			"username",
			"first_name",
			"last_name",
			"email",
			"password",
			"phone",
			"user_status",
		).
		Values(
			user.UserName,
			user.FirstName,
			user.LastName,
			user.Email,
			user.Password,
			user.Phone,
			user.UserStatus,
		).
		RunWith(u.db).
		ExecContext(ctx)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (u *userRepository) CreateUsersWithArrayInput(ctx context.Context, users []models.User) error {
	insertBuilder := sq.Insert(usersTable).Columns(
		"username",
		"first_name",
		"last_name",
		"email",
		"password",
		"phone",
		"user_status",
	)

	for _, user := range users {
		insertBuilder = insertBuilder.Values(
			user.UserName,
			user.FirstName,
			user.LastName,
			user.Email,
			user.Password,
			user.Phone,
			user.UserStatus,
		)
	}

	_, err := insertBuilder.RunWith(u.db).ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) CreateUsersWithListInput(ctx context.Context, users []models.User) error {
	return u.CreateUsersWithArrayInput(ctx, users)
}
