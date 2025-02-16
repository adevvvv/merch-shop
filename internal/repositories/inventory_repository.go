// internal/repositories/inventory_repository.go
package repositories

import (
	"context"
	"database/sql"
	"shop/internal/models"
)

type InventoryRepository interface {
	GetByUserID(ctx context.Context, userID int) ([]models.InventoryItem, error)
	AddItem(ctx context.Context, userID int, itemType string, quantity int) error
}

type inventoryRepository struct {
	db *sql.DB
}

func (r *inventoryRepository) GetByUserID(ctx context.Context, userID int) ([]models.InventoryItem, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT type, quantity FROM inventory WHERE user_id = $1",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		if err := rows.Scan(&item.Type, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *inventoryRepository) AddItem(ctx context.Context, userID int, itemType string, quantity int) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO inventory(user_id, type, quantity) 
		VALUES($1, $2, $3) 
		ON CONFLICT (user_id, type) 
		DO UPDATE SET quantity = inventory.quantity + $3`,
		userID, itemType, quantity,
	)
	return err
}