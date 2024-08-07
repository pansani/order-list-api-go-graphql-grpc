package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	var err error
	connStr := "user=user password=password dbname=orders_db host=db port=5432 sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initDB()

	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: &Resolver{}}))

	http.Handle("/graphql", srv)
	http.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))

	log.Println("GraphQL server started on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
