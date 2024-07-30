package main

import (
	"context"
	"log"

	_ "github.com/lib/pq"
)

type Resolver struct{}

type queryResolver struct{ *Resolver }

func (r *queryResolver) ListOrders(ctx context.Context) ([]*Order, error) {
	rows, err := db.QueryContext(ctx, "SELECT id, user_id, product_id, quantity, status, created_at, updated_at FROM orders")
	if err != nil {
		log.Println("Error querying orders:", err)
		return nil, err
	}
	defer rows.Close()

	var orders []*Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.ProductID, &order.Quantity, &order.Status, &order.CreatedAt, &order.UpdatedAt); err != nil {
			log.Println("Error scanning order:", err)
			return nil, err
		}
		orders = append(orders, &order)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, err
	}
	return orders, nil
}

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }
