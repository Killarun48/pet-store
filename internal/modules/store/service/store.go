package service

import (
	"app/internal/models"
	"context"
	"database/sql"
	"errors"
)

type StoreServicer interface {
	GetInventory(ctx context.Context) (map[string]int, error)
	PlaceOrder(ctx context.Context, order models.Order) (models.Order, error)
	GetOrderById(ctx context.Context, id int) (models.Order, error)
	DeleteOrder(ctx context.Context, id int) error
}

type StoreRepositoryer interface {
	GetInventory(ctx context.Context) (map[string]int, error)
	PlaceOrder(ctx context.Context, order models.Order) (models.Order, error)
	GetOrderById(ctx context.Context, id int) (models.Order, error)
	DeleteOrder(ctx context.Context, id int) error
}

type StoreService struct {
	storeRepository StoreRepositoryer
}

func NewStoreService(storeRepository StoreRepositoryer) StoreServicer {
	return &StoreService{
		storeRepository: storeRepository,
	}
}

func (s StoreService) GetInventory(ctx context.Context) (map[string]int, error) {
	return s.storeRepository.GetInventory(ctx)
}

func (s StoreService) PlaceOrder(ctx context.Context, order models.Order) (models.Order, error) {
	return s.storeRepository.PlaceOrder(ctx, order)
}

func (s StoreService) GetOrderById(ctx context.Context, id int) (models.Order, error) {
	order, err := s.storeRepository.GetOrderById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Order{}, errors.New("order not found")
		}
	}

	return order, nil
}

func (s StoreService) DeleteOrder(ctx context.Context, id int) error {
	order, err := s.GetOrderById(ctx, id)
	if err != nil {
		return err
	}

	if order.Status == "deleted" {
		return errors.New("order already deleted")
	}

	return s.storeRepository.DeleteOrder(ctx, id)
}
