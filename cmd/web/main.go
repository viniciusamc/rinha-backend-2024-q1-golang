package main

import (
	"context"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func main() {
	r := mux.NewRouter()

    conn, err := pgxpool.New(context.Background(), "postgres://admin:123@db:5432/rinha")
    db = conn 

    if err != nil {
        fmt.Println(err)
    }

    defer db.Close()

	r.HandleFunc("/clientes/{id}/transacoes", transacoes).Methods("POST")
	r.HandleFunc("/clientes/{id}/extrato", extrato).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}
