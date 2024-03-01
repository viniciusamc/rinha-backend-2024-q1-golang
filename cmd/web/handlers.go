package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/gorilla/mux"
)

type RequestTransacao struct {
    Valor     int    `json:"valor"`
    Tipo      string `json:"tipo"`
    Descricao string `json:"descricao"`
}

func transacoes(w http.ResponseWriter, r *http.Request) {
	vars, err := strconv.Atoi(mux.Vars(r)["id"]) 

    if err != nil {
        w.WriteHeader(http.StatusUnprocessableEntity)
        return
    }

    if vars > 5 || vars < 1 {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    var request RequestTransacao

    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        w.WriteHeader(http.StatusUnprocessableEntity)
        return
    }

    if request.Valor < 0 {
        w.WriteHeader(http.StatusUnprocessableEntity)
        return
    }

    if request.Tipo != "d" && request.Tipo != "c" {
        w.WriteHeader(http.StatusUnprocessableEntity)
        return
    }

    if utf8.RuneCountInString(request.Descricao) > 10 || utf8.RuneCountInString(request.Descricao) < 1 {
        w.WriteHeader(http.StatusUnprocessableEntity)
        return
    }

    novoSaldo, novoLimite, err := Transacao(vars, request.Valor, request.Tipo, request.Descricao)

    if err != nil {
        w.WriteHeader(http.StatusUnprocessableEntity)
        return
    }

    response := map[string]int{
        "saldo":  novoSaldo,
        "limite": novoLimite,
    }

    err = json.NewEncoder(w).Encode(response)

    if err != nil {
        fmt.Print("error in enconde")
        return
    }
}

func extrato(w http.ResponseWriter, r *http.Request) {
	vars, err := strconv.Atoi(mux.Vars(r)["id"]) 

    if err != nil {
        w.WriteHeader(http.StatusUnprocessableEntity)
        return
    }

    if vars > 5 || vars < 1 {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    extrato := Extrato(vars)

    err = json.NewEncoder(w).Encode(extrato)

    if err != nil {
        fmt.Print("error in enconde")
        return
    }
}
