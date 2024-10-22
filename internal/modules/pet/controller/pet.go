package controller

import (
	"app/internal/infrastructure/responder"
	"app/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

type PetControllerer interface {
	UploadFile(w http.ResponseWriter, r *http.Request)
	AddPet(w http.ResponseWriter, r *http.Request)
	UpdatePet(w http.ResponseWriter, r *http.Request)
	FindPetsByStatus(w http.ResponseWriter, r *http.Request)
	FindPetsByTags(w http.ResponseWriter, r *http.Request)
	GetPetById(w http.ResponseWriter, r *http.Request)
	UpdatePetWithForm(w http.ResponseWriter, r *http.Request)
	DeletePet(w http.ResponseWriter, r *http.Request)
}

type PetServicer interface {
	UploadFile(ctx context.Context, id int, file multipart.File, fileHeader *multipart.FileHeader) error
	AddPet(ctx context.Context, pet models.Pet) (models.Pet, error)
	UpdatePet(ctx context.Context, pet models.Pet) (models.Pet, error)
	FindPetsByStatus(ctx context.Context, status []string) ([]models.Pet, error)
	FindPetsByTags(ctx context.Context, tags []string) ([]models.Pet, error)
	GetPetById(ctx context.Context, id int) (models.Pet, error)
	UpdatePetWithForm(ctx context.Context, id int, name string, status string) error
	DeletePet(ctx context.Context, id int) error
}

type PetController struct {
	petService PetServicer
	responder  responder.Responder
}

func NewPetController(petService PetServicer, responder responder.Responder) PetControllerer {
	return &PetController{
		petService: petService,
		responder:  responder,
	}
}

//	@id			1uploadFile
//	@x-sort		1
//	@Security	ApiKeyAuth
//	@Summary	uploads an image
//	@Tags		pet
//	@Accept		multipart/form-data
//	@Produce	json
//	@Param		petId				path		int		true	"ID of pet to update"
//	@Param		additionalMetadata	formData	string	false	"Additional data to pass to server"
//	@Param		file				formData	file	true	"file to upload"
//	@Success	200					{object}	responder.Response
//	@Router		/pet/{petId}/uploadImage [post]
func (c *PetController) UploadFile(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "petId"))
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	additionalMetadata := r.FormValue("additionalMetadata")

	file, header, err := r.FormFile("file")
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}
	defer file.Close()

	err = c.petService.UploadFile(context.Background(), id, file, header)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	respStr := fmt.Sprintf("additionalMetadata: %s\nFile uploaded to ./%s, %d bytes", additionalMetadata, header.Filename, header.Size)
	c.responder.Success(w, respStr)
}

//	@id			2addPet
//	@x-sort		2
//	@Security	ApiKeyAuth
//	@Summary	Add a new pet to the store
//	@Tags		pet
//	@Accept		json
//	@Produce	json
//	@Param		object	body		models.Pet	true	"Pet object that needs to be added to the store"
//	@Success	200		{object}	models.Pet
//	@Router		/pet [post]
func (c *PetController) AddPet(w http.ResponseWriter, r *http.Request) {
	var pet models.Pet
	err := json.NewDecoder(r.Body).Decode(&pet)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	createPet, err := c.petService.AddPet(context.Background(), pet)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	jsonResp, err := json.MarshalIndent(createPet, "", "  ")
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	fmt.Fprintln(w, string(jsonResp))
}

//	@id			3updatePet
//	@x-sort		3
//	@Security	ApiKeyAuth
//	@Summary	Update an existing pet
//	@Tags		pet
//	@Accept		json
//	@Produce	json
//	@Param		object	body		models.Pet	true	"Pet object that needs to be added to the store"
//	@Success	200		{object}	models.Pet
//	@Router		/pet [put]
func (c *PetController) UpdatePet(w http.ResponseWriter, r *http.Request) {
	var pet models.Pet
	err := json.NewDecoder(r.Body).Decode(&pet)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	updatePet, err := c.petService.UpdatePet(context.Background(), pet)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	jsonResp, err := json.MarshalIndent(updatePet, "", "  ")
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	fmt.Fprintln(w, string(jsonResp))
}

//	@id				4findPetsByStatus
//	@x-sort			4
//	@Security		ApiKeyAuth
//	@Summary		Finds Pets by status
//	@Description	Multiple status values can be provided with comma separated strings
//	@Tags			pet
//	@Accept			json
//	@Produce		json
//	@Param			status	query		[]string	true	"Status values that need to be considered for filter"	Enums(available, pending, sold)
//	@Success		200		{object}	[]models.Pet
//	@Router			/pet/findByStatus [get]
func (c *PetController) FindPetsByStatus(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	statusSlice := strings.Split(status, ",")
	pets, err := c.petService.FindPetsByStatus(context.Background(), statusSlice)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	jsonResp, err := json.MarshalIndent(pets, "", "  ")
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	fmt.Fprintln(w, string(jsonResp))
}

//	@id			5findPetsByTags
//	@x-sort		5
//	@Security	ApiKeyAuth
//	@Summary	Finds Pets by tags
//	@Deprecated
//	@Tags		pet
//	@Accept		json
//	@Produce	json
//	@Param		tags	query		[]string	true	"Tags to filter by"	collectionFormat(csv)
//	@Success	200		{object}	[]models.Pet
//	@Router		/pet/findByTags [get]
func (c *PetController) FindPetsByTags(w http.ResponseWriter, r *http.Request) {
	c.responder.ErrorBadRequest(w, errors.New("deprecated"))
}

//	@id			6getPetById
//	@x-sort		6
//	@Security	ApiKeyAuth
//	@Summary	Find pet by ID
//	@Tags		pet
//	@Accept		json
//	@Produce	json
//	@Param		petId	path		int	true	"ID of pet to return"
//	@Success	200		{object}	models.Pet
//	@Router		/pet/{petId} [get]
func (c *PetController) GetPetById(w http.ResponseWriter, r *http.Request) {
	petID := chi.URLParam(r, "petId")
	id, err := strconv.Atoi(petID)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	pet, err := c.petService.GetPetById(context.Background(), id)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	jsonResp, err := json.MarshalIndent(pet, "", "  ")
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	fmt.Fprintln(w, string(jsonResp))
}

//	@id			7updatePetWithForm
//	@x-sort		7
//	@Security	ApiKeyAuth
//	@Summary	Updates a pet in the store with form data
//	@Tags		pet
//	@Accept		application/x-www-form-urlencoded
//	@Produce	json
//	@Param		petId	path		int		true	"ID of pet that needs to be updated"
//	@Param		name	formData	string	false	"Updated name of the pet"
//	@Param		status	formData	string	false	"Updated status of the pet"
//	@Success	200		{object}	responder.Response
//	@Router		/pet/{petId} [post]
func (c *PetController) UpdatePetWithForm(w http.ResponseWriter, r *http.Request) {
	petID := chi.URLParam(r, "petId")
	id, err := strconv.Atoi(petID)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	//multipart/form-data
	/* err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	} */

	name := r.FormValue("name")
	status := r.FormValue("status")

	err = c.petService.UpdatePetWithForm(context.Background(), id, name, status)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	c.responder.Success(w, fmt.Sprint(id))
}

//	@id			8deletePet
//	@x-sort		8
//	@Security	ApiKeyAuth
//	@Summary	Deletes a pet
//	@Tags		pet
//	@Accept		json
//	@Produce	json
//	@Param		petId	path		int	true	"Pet id to delete"
//	@Success	200		{object}	responder.Response
//	@Router		/pet/{petId} [delete]
func (c *PetController) DeletePet(w http.ResponseWriter, r *http.Request) {
	petID := chi.URLParam(r, "petId")
	id, err := strconv.Atoi(petID)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	err = c.petService.DeletePet(context.Background(), id)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	c.responder.Success(w, fmt.Sprint(id))
}
