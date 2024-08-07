package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Order struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	ProductID int    `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

var db *sql.DB

func listOrders(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, user_id, product_id, quantity, status, created_at, updated_at FROM orders")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.ProductID, &order.Quantity, &order.Status, &order.CreatedAt, &order.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		orders = append(orders, order)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO orders (user_id, product_id, quantity, status, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id, created_at, updated_at`
	err := db.QueryRow(query, order.UserID, order.ProductID, order.Quantity, order.Status).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func initDB() {
	var err error
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
	log.Println("Connection string:", connStr)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initDB()

	r := mux.NewRouter()
	r.HandleFunc("/order", listOrders).Methods("GET")
	r.HandleFunc("/order", createOrder).Methods("POST")

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
