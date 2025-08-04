package product

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"time"
	"url_shortness/internal/domain/entities/Table"
	"url_shortness/internal/infra/database/query"
	"url_shortness/pkg/logger"
)

func GetUser(pool *pgxpool.Pool, rdb *redis.Client, login string) (Table.Account, error) {
	cacheKey := fmt.Sprintf("user:%s", login)
	ctx := context.Background()

	var account Table.Account
	val, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		if err = json.Unmarshal([]byte(val), &account); err == nil {
			return account, nil
		}

		rdb.Del(ctx, cacheKey)
	}

	account, err = query.GetAccount(pool, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return account, fmt.Errorf("Пользователь не найден")
		}
		logger.Get().Error("DB error: " + err.Error())
		return account, fmt.Errorf("Ошибка сервера")
	}

	if accountJSON, err := json.Marshal(account); err == nil {
		rdb.Set(ctx, cacheKey, accountJSON, 15*time.Minute)
	} else {
		logger.Get().Warn("Cache marshal error: " + err.Error())
	}

	return account, nil
}
