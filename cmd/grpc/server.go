package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	pb "github.com/pansani/order-list-go/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var db *sql.DB

type server struct {
	pb.UnimplementedOrderServiceServer
}

func (s *server) ListOrders(ctx context.Context, in *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	rows, err := db.Query("SELECT id, user_id, product_id, quantity, status, created_at, updated_at FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*pb.Order
	for rows.Next() {
		var order pb.Order
		if err := rows.Scan(&order.Id, &order.UserId, &order.ProductId, &order.Quantity, &order.Status, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	return &pb.ListOrdersResponse{Orders: orders}, nil
}

func (s *server) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	now := time.Now().Format(time.RFC3339)
	var orderID int32

	err := db.QueryRowContext(ctx, `
		INSERT INTO orders (user_id, product_id, quantity, status, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id`,
		in.UserId, in.ProductId, in.Quantity, in.Status, now, now).Scan(&orderID)

	if err != nil {
		return nil, err
	}

	order := &pb.Order{
		Id:        orderID,
		UserId:    in.UserId,
		ProductId: in.ProductId,
		Quantity:  in.Quantity,
		Status:    in.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return &pb.CreateOrderResponse{Order: order}, nil
}

func initDB() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initDB()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, &server{})

	reflection.Register(s)

	log.Println("gRPC server started on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
