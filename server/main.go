package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/chanudinho/fullcycle-desafio-client-server/server/database"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	db, err := database.Setup()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = database.Migrate(db)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/cotacao", getCotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

func getCotacaoHandler(w http.ResponseWriter, r *http.Request) {
	cotacao, err := getCotacao()
	if err != nil {
		http.Error(w, "Erro ao buscar cotação", http.StatusInternalServerError)
		return
	}

	err = database.InsertCotacao(cotacao.Bid)
	if err != nil {
		http.Error(w, "Erro ao inserir cotação no banco de dados", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cotacao)
}

func getCotacao() (Cotacao, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	url := "https://economia.awesomeapi.com.br/json/USD-BRL"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Cotacao{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("request canceled by client")
		}

		return Cotacao{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Cotacao{}, err
	}

	var cotacao []Cotacao
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		return Cotacao{}, err
	}

	return cotacao[0], nil
}
