package main

import (
	"context"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Resolver struct{}

type queryResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }

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

func (r *mutationResolver) CreateOrder(ctx context.Context, userID int, productID int, quantity int, status string) (*Order, error) {
	now := time.Now().Format(time.RFC3339)
	var orderID int

	err := db.QueryRowContext(ctx, `
		INSERT INTO orders (user_id, product_id, quantity, status, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id`,
		userID, productID, quantity, status, now, now).Scan(&orderID)

	if err != nil {
		log.Println("Error inserting order:", err)
		return nil, err
	}

	return &Order{
		ID:        orderID,
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (r *Resolver) Query() QueryResolver       { return &queryResolver{r} }
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }
