package controller

import (
	"app/internal/infrastructure/responder"
	"app/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type StoreControllerer interface {
	GetInventory(w http.ResponseWriter, r *http.Request)
	PlaceOrder(w http.ResponseWriter, r *http.Request)
	GetOrderById(w http.ResponseWriter, r *http.Request)
	DeleteOrder(w http.ResponseWriter, r *http.Request)
}

type StoreServicer interface {
	GetInventory(ctx context.Context) (map[string]int, error)
	PlaceOrder(ctx context.Context, order models.Order) (models.Order, error)
	GetOrderById(ctx context.Context, id int) (models.Order, error)
	DeleteOrder(ctx context.Context, id int) error
}

type StoreController struct {
	storeService StoreServicer
	responder    responder.Responder
}

func NewStoreController(storeService StoreServicer, responder responder.Responder) StoreControllerer {
	return &StoreController{
		storeService: storeService,
		responder:    responder,
	}
}

//	@id				1getInventory
//	@Security		ApiKeyAuth
//	@Summary		Returns pet inventories by status
//	@Description	Returns a map of status codes to quantities
//	@Tags			store
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]int
//	@Router			/store/inventory [get]
func (sc StoreController) GetInventory(w http.ResponseWriter, r *http.Request) {
	inventory, err := sc.storeService.GetInventory(context.Background())
	if err != nil {
		sc.responder.ErrorBadRequest(w, err)
		return
	}

	jsonResp, err := json.MarshalIndent(inventory, "", "  ")
	if err != nil {
		sc.responder.ErrorBadRequest(w, err)
		return
	}

	fmt.Fprintln(w, string(jsonResp))
}

//	@id			2placeOrder
//	@Summary	Place an order for a pet
//	@Tags		store
//	@Accept		json
//	@Produce	json
//	@Param		object	body		models.Order	true	"order placed for purchasing the pet"
//	@Success	200		{object}	models.Order
//	@Router		/store/order [post]
func (sc StoreController) PlaceOrder(w http.ResponseWriter, r *http.Request) {

	var order models.Order

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		sc.responder.ErrorBadRequest(w, err)
		return
	}

	createOrder, err := sc.storeService.PlaceOrder(context.Background(), order)
	if err != nil {
		sc.responder.ErrorBadRequest(w, err)
		return
	}

	jsonResp, err := json.MarshalIndent(createOrder, "", "  ")
	if err != nil {
		sc.responder.ErrorBadRequest(w, err)
		return
	}

	fmt.Fprintln(w, string(jsonResp))
}

//	@id				3getOrderById
//	@Summary		Find purchase order by ID
//	@Description	For valid response try integer IDs with value >= 1 and <= 10. Other values will generated exceptions
//	@Tags			store
//	@Accept			json
//	@Produce		json
//	@Param			orderId	path		int	true	"ID of pet that needs to be fetched"
//	@Success		200		{object}	models.Order
//	@Router			/store/order/{orderId} [get]
func (sc StoreController) GetOrderById(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderId")
	id, err := strconv.Atoi(orderID)
	if err != nil {
		sc.responder.ErrorBadRequest(w, err)
		return
	}

	order, err := sc.storeService.GetOrderById(context.Background(), id)
	if err != nil {
		sc.responder.ErrorBadRequest(w, err)
		return
	}

	jsonResp, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		sc.responder.ErrorBadRequest(w, err)
		return
	}

	fmt.Fprintln(w, string(jsonResp))
}

//	@id				4deleteOrder
//	@Summary		Delete purchase order by ID
//	@Description	For valid response try integer IDs with positive integer value. Negative or non-integer values will generate API errors
//	@Tags			store
//	@Accept			json
//	@Produce		json
//	@Param			orderId	path		int	true	"ID of pet that needs to be deleted"
//	@Success		200		{object}	responder.Response
//	@Router			/store/order/{orderId} [delete]
func (sc StoreController) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderId")
	id, err := strconv.Atoi(orderID)
	if err != nil {
		sc.responder.ErrorBadRequest(w, err)
		return
	}	

	err = sc.storeService.DeleteOrder(context.Background(), id)
	if err != nil {
		sc.responder.ErrorBadRequest(w, err)
		return
	}

	sc.responder.Success(w, fmt.Sprint(id))
}