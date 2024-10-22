package repository

import (
	"app/internal/models"
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type StoreRepositoryer interface {
	GetInventory(ctx context.Context) (map[string]int, error)
	PlaceOrder(ctx context.Context, order models.Order) (models.Order, error)
	GetOrderById(ctx context.Context, id int) (models.Order, error)
	DeleteOrder(ctx context.Context, id int) error
}

type StoreRepository struct {
	db *sql.DB
}

func NewStoreRepository(db *sql.DB) StoreRepositoryer {
	return &StoreRepository{
		db: db,
	}
}

func (r StoreRepository) GetInventory(ctx context.Context) (map[string]int, error) {
	var available, pending, sold sql.NullInt64

	err := sq.Select(
		"COUNT(id) FILTER(WHERE status = 'available')",
		"COUNT(id) FILTER(WHERE status = 'pending')",
		"COUNT(id) FILTER(WHERE status = 'sold')",
		).
		From("pets").
		RunWith(r.db).
		ScanContext(ctx, &available, &pending, &sold)
	if err != nil {
		return nil, err
	}

	inventory := map[string]int{
		"available": int(available.Int64),
		"pending":   int(pending.Int64),
		"sold":      int(sold.Int64),
	}

	return inventory, nil
}	

func (r StoreRepository) PlaceOrder(ctx context.Context, order models.Order) (models.Order, error) {
	res, err := sq.Insert("orders").
		Columns("pet_id", "quantity", "ship_date", "status", "complete").
		Values(order.PetID, order.Quantity, order.ShipDate, order.Status, order.Complete).
		RunWith(r.db).
		ExecContext(ctx)
	if err != nil {
		return models.Order{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.Order{}, err
	}

	order.ID = int(id)

	return order, nil
}

func (r StoreRepository) GetOrderById(ctx context.Context, id int) (models.Order, error) {
	var order models.Order
	err := sq.Select(
		"id", 
		"pet_id",
		"quantity",
		"ship_date",
		"status",
		"complete",
		).
		From("orders").
		Where(sq.Eq{"id": id}).
		RunWith(r.db).
		ScanContext(ctx, &order.ID, &order.PetID, &order.Quantity, &order.ShipDate, &order.Status, &order.Complete)
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (r StoreRepository) DeleteOrder(ctx context.Context, id int) error {
	_, err := sq.Update("orders").
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
