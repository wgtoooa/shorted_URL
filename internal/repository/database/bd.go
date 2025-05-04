package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

var DB *pgxpool.Pool

func BDinit(dsn string) {
	var err error
	DB, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Println("failed to connected to PostgresSQL")
		panic(err)
	}
	log.Println("Successfully connected to PostgresSQL")
	return
}
