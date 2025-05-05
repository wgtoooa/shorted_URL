package query

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"url_shortness/internal/repository/Table"
)

func CreateTableURL(pool *pgxpool.Pool) (err error) {
	query := `
	Create table if not exists URL(
		id serial primary key,
		full_url varchar not null unique,
		short_url varchar not null unique,
		account_id int not null references Account(id),
		created_at timestamp default now()
	    )`
	_, err = pool.Exec(context.Background(), query)
	return
}

func GetURLS(pool *pgxpool.Pool, login string) ([]Table.URL, error) {
	query := `
	SELECT u.id, u.full_url, u.short_url, u.created_at
	FROM url u
	JOIN account a ON u.account_id = a.id
	WHERE a.login = $1
	`
	rows, err := pool.Query(context.Background(), query, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []Table.URL
	for rows.Next() {
		var u Table.URL
		err = rows.Scan(&u.Id, &u.Full_url, &u.Short_url, &u.Created_at)
		if err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func CreateURL(pool *pgxpool.Pool, full_url, short_url string, account_id int) (err error) {
	query := `insert into URL(full_url,short_url,account_id) values ($1,$2,$3)`
	_, err = pool.Exec(context.Background(), query, full_url, short_url, account_id)
	return
}

func GetURLByShortURL(pool *pgxpool.Pool, shortURL string) (fullURL string, err error) {
	query := `SELECT full_url FROM url WHERE short_url = $1`
	err = pool.QueryRow(context.Background(), query, shortURL).Scan(&fullURL)
	return
}
