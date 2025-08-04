package query

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"url_shortness/internal/domain/entities/Table"
	"url_shortness/pkg/logger"
)

func CreatedTableAccount(pool *pgxpool.Pool) (err error) {
	query := `
	Create table if not exists Account (
		id serial primary key,
		login varchar not null unique,
		password varchar not null,
		created_at timestamp default now(),	    
	    countURL int default 0
	)`

	_, err = pool.Exec(context.Background(), query)
	if err != nil {
		logger.Get().Error("failed create table account", zap.Error(err))
	}
	return
}

func CreateAccount(pool *pgxpool.Pool, login string, password string) (err error) {
	query := `	insert into Account(login,password) values($1,$2)`
	_, err = pool.Exec(context.Background(), query, login, password)
	return
}

func GetAccount(pool *pgxpool.Pool, login string) (user Table.Account, err error) {
	query := `select id,login,password,countURL,created_at from Account where login = $1`
	err = pool.QueryRow(context.Background(), query, login).
		Scan(&user.Id, &user.Login, &user.Password, &user.CountURL, &user.CreatedAt)

	return
}

func AccountExists(pool *pgxpool.Pool, login string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM account WHERE login = $1)`

	err := pool.QueryRow(context.Background(), query, login).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check account existence: %w", err)
	}

	return exists, nil
}
