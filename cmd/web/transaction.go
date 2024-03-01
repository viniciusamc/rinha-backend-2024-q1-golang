package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Cliente struct {
	Saldo  int
	Limite int
}

type TransacaoList struct {
	Valor       int       `json:"valor"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	RealizadaEm time.Time `json:"realizada_em"`
}

type Extract struct {
	Saldo      Account         `json:"saldo"`
	Transacoes []TransacaoList `json:"ultimas_transacoes"`
}

type Account struct {
	Saldo       int       `json:"total"`
	RealizadaEm time.Time `json:"realizada_em"`
	Limite      int       `json:"limite"`
}

func Transacao(id int, valor int, tipo string, descricao string) (saldoCliente int, limiteCliente int, err error) {
	dbL, err := db.Begin(context.Background())
	if err != nil {
		fmt.Print(err, "db error")
	}

	defer dbL.Rollback(context.Background())

	// dbL.Exec(context.Background(), "SELECT pg_advisory_xact_lock($1)", string(id))

	var cliente Cliente

	err = dbL.QueryRow(context.Background(), "SELECT saldo, limite FROM clientes WHERE id = $1", id).Scan(&cliente.Saldo, &cliente.Limite)

	if err != nil {
		fmt.Print(err, "db error")
	}

	if tipo == "d" {
		cliente.Saldo = cliente.Saldo - valor
	} else {
		cliente.Saldo = cliente.Saldo + valor
	}

	if tipo == "d" && (cliente.Limite+cliente.Saldo) < 0 {
		return 0, 0, errors.New("422")
	}

	_, err = dbL.Exec(context.Background(), "INSERT INTO transacoes(id_cliente,valor,tipo,descricao, realizada_em) VALUES($1,$2,$3,$4,$5)", id, valor, tipo, descricao, time.Now())
	if err != nil {
		fmt.Print(err, "db error")
	}

	_, err = dbL.Exec(context.Background(), "UPDATE clientes SET saldo = $1 WHERE id = $2", cliente.Saldo, id)
	if err != nil {
		fmt.Print(err, "db error")
	}

	if err != nil {
		fmt.Print(nil)
		return 0, 0, errors.New("422")
	}

	dbL.Commit(context.Background())

	return cliente.Saldo, cliente.Limite, nil
}

func Extrato(id int) Extract {
	dbL, err := db.Begin(context.Background())
	if err != nil {
		fmt.Print(err, "db error")
	}

	defer dbL.Rollback(context.Background())

	results, err := dbL.Query(context.Background(), "SELECT valor, tipo, descricao, realizada_em FROM transacoes WHERE id_cliente = $1 ORDER BY id DESC LIMIT 10", id)
	if err != nil {
		fmt.Print(nil)
	}

	var transacoes []TransacaoList

	for results.Next() {
		var transacao TransacaoList
		err = results.Scan(&transacao.Valor, &transacao.Tipo, &transacao.Descricao, &transacao.RealizadaEm)
		if err != nil {
			fmt.Print(err)
		}
		transacoes = append(transacoes, transacao)
	}

	var AccountDetails Account

	err = dbL.QueryRow(context.Background(), "SELECT saldo, limite FROM clientes WHERE id = $1", id).Scan(&AccountDetails.Saldo, &AccountDetails.Limite)
	if err != nil {
		fmt.Print(err)
	}

	AccountDetails.RealizadaEm = time.Now()

	var extract Extract

	extract.Transacoes = transacoes
	extract.Saldo = AccountDetails

	return extract
}
