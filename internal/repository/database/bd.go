package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"url_shortness/internal/logger"
	"url_shortness/internal/repository/database/query"
)

type DataBaseConfig struct {
	User            string
	DBName          string
	Password        string
	Host            string
	Port            string
	SSLMode         string
	MaxConns        int32         // Максимальное количество соединений в пуле
	MinConns        int32         // Минимальное количество соединений в пуле
	MaxConnLifetime time.Duration // Максимальное время жизни соединения
	MaxConnIdleTime time.Duration // Максимальное время бездействия соединения
	ConnectTimeout  time.Duration // Таймаут подключения
}

func InitDB(config DataBaseConfig) (*pgxpool.Pool, error) {
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

	logger.Get().Info("Successfully connected to PostgreSQL")

	// Создаем таблицы
	if err = query.CreatedTableAccount(dbPool); err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("failed to create account table: %w", err)
	}

	if err = query.CreateTableURL(dbPool); err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("failed to create URL table: %w", err)
	}

	return dbPool, nil
}
