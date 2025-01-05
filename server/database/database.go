package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func Setup() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "client_server.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(db *sql.DB) error {
	statement, err := db.Prepare(
		"CREATE TABLE IF NOT EXISTS cotacao (id INTEGER PRIMARY KEY AUTOINCREMENT, cotacao VARCHAR(10) NULL)",
	)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	return err
}

func InsertCotacao(cotacao string) error {
	db, err := Setup()
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	stmt, err := db.PrepareContext(ctx, "insert into cotacao(cotacao) values(?)")
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Timeout ao persistir cotação no banco de dados")
		}

		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cotacao)
	return err
}
