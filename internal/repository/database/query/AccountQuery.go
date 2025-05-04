package query

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"url_shortness/internal/repository/Table"
)

func CreatedTableAccount(pool *pgxpool.Pool) (err error) {
	query := `
	Create table if not exists Account (
		id serial primary key,
		login varchar not null unique,
		password varchar not null,
		status varchar default 'online' check(status in ('online','offline')),
		created_at timestamp default now()	    
	)`

	_, err = pool.Exec(context.Background(), query)
	return
}

func CreateAccount(pool *pgxpool.Pool, login string, password string) (err error) {
	query := `	insert into Account(login,password) values($1,$2)`
	_, err = pool.Exec(context.Background(), query, login, password)
	return
}

func GetAccount(pool *pgxpool.Pool, login string) (user Table.Account, err error) {
	query := `select id,login,password,status,created_at,countURL from Account where login = $1`
	err = pool.QueryRow(context.Background(), query, login).
		Scan(&user.Id, &user.Login, &user.Password, &user.Status, &user.CreatedAt, &user.CountURL)

	return
}

func AccountExists(pool *pgxpool.Pool, login string) (bool, error) {
	user, err := GetAccount(pool, login)
	if user == (Table.Account{}) {
		return false, err
	}
	return true, err
}
