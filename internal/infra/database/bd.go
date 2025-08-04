package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	query2 "url_shortness/internal/infra/database/query"
	"url_shortness/pkg/Config"
	"url_shortness/pkg/logger"
)

func InitDB(config config.DataBaseConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s",
		config.User, config.DBName, config.Password, config.Host, config.Port)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Проверяем соединение
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		dbPool.Close() // Закрываем пул при ошибке
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}
	conn.Release()

	logger.Get().Info("Successfully connected to PostgresQL")
	// Создаем таблицы
	if err = query2.CreatedTableAccount(dbPool); err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("failed to create account table: %w", err)
	}

	if err = query2.CreateTableURL(dbPool); err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("failed to create URL table: %w", err)
	}
	logger.Get().Info("Successfully create table")
	return dbPool, nil
}
