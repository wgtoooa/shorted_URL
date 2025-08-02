package query

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"url_shortness/internal/repository/database/Table"
)

func CreateTableURL(pool *pgxpool.Pool) (err error) {
	query := `
	Create table if not exists URL(
		id serial primary key,
		full_url varchar not null ,
		short_url varchar not null unique,
		account_id int not null references Account(id),
		created_at timestamp default now(),
	    FOREIGN KEY (account_id) REFERENCES Account(id) ON DELETE CASCADE
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
	order by id DESC 
	limit 5
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

func DeleteURLByShortURLL(pool *pgxpool.Pool, shortURL string) (err error) {
	query := `DELETE FROM url WHERE short_url = $1`
	_, err = pool.Exec(context.Background(), query, shortURL)
	return
}

func PutShortURL(pool *pgxpool.Pool, OldShortURL, NewShortURL string) error {
	query := `update URL set short_url = $1 where short_url = $2`
	cmdTag, err := pool.Exec(context.Background(), query, NewShortURL, OldShortURL)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("Ссылка для изменения не найдена")
	}

	return nil
}

func IsExistsURL(pool *pgxpool.Pool, account_id int, full_url string) (bool, error) {
	var id int
	query := `SELECT id FROM url WHERE full_url = $1 AND account_id = $2`

	err := pool.QueryRow(context.Background(), query, full_url, account_id).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Запись не найдена
			return false, nil
		}

		return false, fmt.Errorf("database query error: %w", err)
	}

	// Запись найдена
	return true, nil
}

func IsDuplicateKeyError(err error) (bool, string) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			field := "unknown"
			switch pgErr.ConstraintName {
			case "url_short_url_key":
				field = "short_url"
			}

			return true, fmt.Sprintf(
				"duplicate key value violates unique constraint '%s' (field: %s). Details: %s",
				pgErr.ConstraintName,
				field,
				pgErr.Detail,
			)
		}
	}
	return false, ""
}
