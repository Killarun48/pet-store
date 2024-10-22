package repository

import (
	"app/internal/models"
	"context"
	"database/sql"
	"errors"
	"log"
	"mime/multipart"

	sq "github.com/Masterminds/squirrel"
)

const (
	petsTable       = "pets"
	photosTable     = "pet_photos"
	categoriesTable = "categories"
	tagsTable       = "tags"
	tagPetsTable    = "tag_pets"
)

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

type PetRepository struct {
	db *sql.DB
}

func NewPetRepository(db *sql.DB) PetRepositoryer {
	return &PetRepository{
		db: db,
	}
}

func (r PetRepository) UploadFile(ctx context.Context, id int, file multipart.File, fileHeader *multipart.FileHeader) error {
	_, err := sq.Insert("pet_photos").
		Columns("pet_id", "photo_url").
		Values(id, fileHeader.Filename).
		RunWith(r.db).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func (r PetRepository) AddPet(ctx context.Context, pet models.Pet) (models.Pet, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Pet{}, err
	}
	defer tx.Rollback()

	// создание/обновление категории
	_, err = sq.Replace(categoriesTable).
		Columns("id", "name").
		Values(pet.Category.ID, pet.Category.Name).
		RunWith(tx).Exec()
	if err != nil {
		return models.Pet{}, err
	}

	// создание питомца
	res, err := sq.Insert(petsTable).
		Columns("name", "category_id", "status").
		Values(pet.Name, pet.Category.ID, pet.Status).
		RunWith(tx).Exec()
	if err != nil {
		return models.Pet{}, err
	}

	// получение id
	id, err := res.LastInsertId()
	if err != nil {
		return models.Pet{}, err
	}

	pet.ID = int(id)

	// добавление фото
	insertBuilder := sq.Insert(photosTable).Columns("pet_id", "photo_url")
	for _, photo := range pet.PhotoUrls {
		insertBuilder = insertBuilder.Values(pet.ID, photo)
	}

	_, err = insertBuilder.RunWith(tx).Exec()
	if err != nil {
		return models.Pet{}, err
	}

	// создание/обновление тегов
	insertBuilder = sq.Replace(tagsTable).Columns("id", "name")
	for _, tag := range pet.Tags {
		insertBuilder = insertBuilder.Values(tag.ID, tag.Name)
	}

	_, err = insertBuilder.RunWith(tx).Exec()
	if err != nil {
		return models.Pet{}, err
	}

	// добавляние принадлежности тегов
	insertBuilder = sq.Insert(tagPetsTable).Columns("pet_id", "tag_id")
	for _, tag := range pet.Tags {
		insertBuilder = insertBuilder.Values(pet.ID, tag.ID)
	}

	_, err = insertBuilder.RunWith(tx).Exec()
	if err != nil {
		return models.Pet{}, err
	}

	// фиксация транзакции
	err = tx.Commit()
	if err != nil {
		return models.Pet{}, err
	}

	return pet, nil
}

func (pr PetRepository) UpdatePet(ctx context.Context, pet models.Pet) (models.Pet, error) {
	tx, err := pr.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Pet{}, err
	}
	defer tx.Rollback()

	// создание/обновление категории
	_, err = sq.Replace(categoriesTable).
		Columns("id", "name").
		Values(pet.Category.ID, pet.Category.Name).
		RunWith(tx).Exec()
	if err != nil {
		return models.Pet{}, err
	}

	// обновление питомца
	_, err = sq.Update(petsTable).
		SetMap(map[string]interface{}{
			"name":        pet.Name,
			"category_id": pet.Category.ID,
			"status":      pet.Status,
		}).
		Where(sq.Eq{"id": pet.ID}).
		RunWith(tx).Exec()
	if err != nil {
		return models.Pet{}, err
	}

	// сначала снести фото
	_, err = sq.Delete(photosTable).
		Where(sq.Eq{"pet_id": pet.ID}).
		RunWith(tx).Exec()
	if err != nil {
		return models.Pet{}, err
	}

	// добавление фото
	insertBuilder := sq.Insert(photosTable).Columns("pet_id", "photo_url")
	for _, photo := range pet.PhotoUrls {
		insertBuilder = insertBuilder.Values(pet.ID, photo)
	}

	_, err = insertBuilder.RunWith(tx).Exec()
	if err != nil {
		return models.Pet{}, err
	}

	// создание/обновление тегов
	insertBuilder = sq.Replace(tagsTable).Columns("id", "name")
	for _, tag := range pet.Tags {
		insertBuilder = insertBuilder.Values(tag.ID, tag.Name)
	}

	_, err = insertBuilder.RunWith(tx).Exec()
	if err != nil {
		return models.Pet{}, err
	}

	// сначала удаляем принадлежность тегов
	_, err = sq.Delete(tagPetsTable).
		Where(sq.Eq{"pet_id": pet.ID}).
		RunWith(tx).Exec()
	if err != nil {
		return models.Pet{}, err
	}

	// добавляние принадлежности тегов
	insertBuilder = sq.Insert(tagPetsTable).Columns("pet_id", "tag_id")
	for _, tag := range pet.Tags {
		insertBuilder = insertBuilder.Values(pet.ID, tag.ID)
	}

	_, err = insertBuilder.RunWith(tx).Exec()
	if err != nil {
		return models.Pet{}, err
	}

	// фиксация транзакции
	err = tx.Commit()
	if err != nil {
		return models.Pet{}, err
	}

	return pet, nil
}

func (pr PetRepository) FindPetsByStatus(ctx context.Context, status []string) ([]models.Pet, error) {
	rows, err := sq.Select(
		"pets.id",
		"pets.name",
		"pets.category_id",
		"pets.status",
		"categories.id",
		"categories.name",
		"pet_photos.photo_url",
		"tags.id",
		"tags.name",
	).
		From(petsTable).
		LeftJoin("categories ON pets.category_id = categories.id").
		LeftJoin("pet_photos ON pets.id = pet_photos.pet_id").
		LeftJoin("tag_pets ON pets.id = tag_pets.pet_id").
		LeftJoin("tags ON tag_pets.tag_id = tags.id").
		Where(sq.Eq{"pets.status": status}).
		RunWith(pr.db).
		QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	type PetRow struct {
		ID         sql.NullInt64
		Name       sql.NullString
		CategoryID sql.NullInt64
		Status     sql.NullString
		Category   struct {
			ID   sql.NullInt64
			Name sql.NullString
		}
		PhotoUrls []sql.NullString
		Tags      []struct {
			ID   sql.NullInt64
			Name sql.NullString
		}
	}

	var petRowArr []*PetRow

	for rows.Next() {
		var petRow PetRow
		var photo sql.NullString
		var tag struct {
			ID   sql.NullInt64
			Name sql.NullString
		}

		err = rows.Scan(
			&petRow.ID,
			&petRow.Name,
			&petRow.CategoryID,
			&petRow.Status,
			&petRow.Category.ID,
			&petRow.Category.Name,
			&photo,
			&tag.ID,
			&tag.Name,
		)
		if err != nil {
			log.Println(err)
			continue
		}

		var existingPet *PetRow
		for _, p := range petRowArr {
			if p.ID.Valid && p.ID.Int64 == petRow.ID.Int64 {
				existingPet = p
				break
			}
		}

		if existingPet != nil {
			if photo.Valid {
				existingPet.PhotoUrls = append(existingPet.PhotoUrls, photo)
			}

			if tag.ID.Valid {
				existingPet.Tags = append(existingPet.Tags, tag)
			}
		} else {
			if !photo.Valid {
				petRow.PhotoUrls = []sql.NullString{}
			} else {
				petRow.PhotoUrls = append(petRow.PhotoUrls, photo)
			}

			if !tag.ID.Valid {
				petRow.Tags = []struct {
					ID   sql.NullInt64
					Name sql.NullString
				}{}
			} else {
				petRow.Tags = append(petRow.Tags, tag)
			}

			petRowArr = append(petRowArr, &petRow)
		}
	}

	result := make([]models.Pet, len(petRowArr))
	for i, p := range petRowArr {
		result[i] = models.Pet{
			ID:     int(p.ID.Int64),
			Name:   p.Name.String,
			Status: p.Status.String,
			Category: models.Category{
				ID:   int(p.Category.ID.Int64),
				Name: p.Category.Name.String,
			},
		}

		for _, photo := range p.PhotoUrls {
			result[i].PhotoUrls = append(result[i].PhotoUrls, photo.String)
		}

		for _, tag := range p.Tags {
			result[i].Tags = append(result[i].Tags, models.Tag{
				ID:   int(tag.ID.Int64),
				Name: tag.Name.String,
			})
		}
	}

	return result, nil
}

func (pr PetRepository) FindPetsByTags(ctx context.Context, tags []string) ([]models.Pet, error) {
	return nil, nil
}

func (r PetRepository) GetPetById(ctx context.Context, id int) (models.Pet, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM pets WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		return models.Pet{}, err
	}

	if !exists {
		return models.Pet{}, errors.New("pet not found")
	}

	var pet models.Pet

	rows, err := sq.Select(
		"pets.id",
		"pets.name",
		"pets.category_id",
		"pets.status",
		"categories.id",
		"categories.name",
		"pet_photos.photo_url",
		"tags.id",
		"tags.name",
	).
		From(petsTable).
		LeftJoin("categories ON pets.category_id = categories.id").
		LeftJoin("pet_photos ON pets.id = pet_photos.pet_id").
		LeftJoin("tag_pets ON pets.id = tag_pets.pet_id").
		LeftJoin("tags ON tag_pets.tag_id = tags.id").
		Where(sq.Eq{"pets.id": id}).
		RunWith(r.db).
		QueryContext(ctx)
	if err != nil {
		log.Println(err)
		return models.Pet{}, err
	}

	type PetRow struct {
		ID         sql.NullInt64
		Name       sql.NullString
		CategoryID sql.NullInt64
		Status     sql.NullString
		Category   struct {
			ID   sql.NullInt64
			Name sql.NullString
		}
		PhotoUrls []sql.NullString
		Tags      []struct {
			ID   sql.NullInt64
			Name sql.NullString
		}
	}

	var petRow PetRow

	for rows.Next() {
		var photo sql.NullString
		var tag struct {
			ID   sql.NullInt64
			Name sql.NullString
		}

		err = rows.Scan(
			&petRow.ID,
			&petRow.Name,
			&petRow.CategoryID,
			&petRow.Status,
			&petRow.Category.ID,
			&petRow.Category.Name,
			&photo,
			&tag.ID,
			&tag.Name,
		)
		if err != nil {
			log.Println(err)
			return models.Pet{}, err
		}

		petRow.PhotoUrls = append(petRow.PhotoUrls, photo)
		petRow.Tags = append(petRow.Tags, tag)
	}

	pet.ID = int(petRow.ID.Int64)
	pet.Name = petRow.Name.String
	pet.Category.ID = int(petRow.Category.ID.Int64)
	pet.Category.Name = petRow.Category.Name.String
	pet.Status = petRow.Status.String

	for _, photo := range petRow.PhotoUrls {
		pet.PhotoUrls = append(pet.PhotoUrls, photo.String)
	}

	for _, tag := range petRow.Tags {
		t := models.Tag{
			ID:   int(tag.ID.Int64),
			Name: tag.Name.String,
		}

		pet.Tags = append(pet.Tags, t)
	}

	return pet, nil
}

func (r PetRepository) UpdatePetWithForm(ctx context.Context, id int, name string, status string) error {
	updateMap := sq.Eq{}

	if name == "" && status == "" {
		return nil
	}

	if name != "" {
		updateMap["name"] = name
	}

	if status != "" {
		updateMap["status"] = status
	}

	_, err := sq.Update(petsTable).
		SetMap(updateMap).
		Where(sq.Eq{"id": id}).
		RunWith(r.db).ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r PetRepository) DeletePet(ctx context.Context, id int) error {
	_, err := sq.Update(petsTable).
		SetMap(map[string]interface{}{
			"status": "deleted",
		}).
		Where(sq.Eq{"id": id}).
		RunWith(r.db).
		ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}
