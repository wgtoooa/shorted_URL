package query

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTableURL(pool *pgxpool.Pool) (err error) {
	query := `
	Create table if not exists URL(
		id serial primary key,
		full_url varchar not null,
		short_url varchar not null,
		account_id int not null references Account(id),
		created_at timestamp default now()
	    
	)`
	_, err = pool.Exec(context.Background(), query)
	return
}
