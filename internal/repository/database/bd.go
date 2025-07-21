package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"log"
	"time"
	"url_shortness/internal/repository/database/query"
	"url_shortness/pkg/logger"
)

var DB *pgxpool.Pool

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

func BDinit(config DataBaseConfig) {

	dsn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s",
		config.User, config.DBName, config.Password, config.Host, config.Port)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Log.Error(err.Error())
	}
	log.Println(dsn)
	//// Устанавливаем параметры пула
	//poolConfig.MaxConns = config.MaxConns
	//poolConfig.MinConns = config.MinConns
	//poolConfig.MaxConnLifetime = config.MaxConnLifetime
	//poolConfig.MaxConnIdleTime = config.MaxConnIdleTime
	//poolConfig.ConnConfig.ConnectTimeout = config.ConnectTimeout TODO: Разобраться

	// Создаем пул соединений
	DB, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		logger.Log.Error("failed to connect to PostgreSQL", zap.Error(err))
		return
	}

	// Проверяем соединение
	conn, err := DB.Acquire(context.Background())
	if err != nil {
		log.Println("failed to acquire connection from pool")
		panic(err)
	}
	conn.Release()

	log.Println("Successfully connected to PostgreSQL")

	err = query.CreatedTableAccount(DB) // create table accounts if exists
	if err != nil {
		logger.Log.Error("failed create table account", zap.Error(err))
		return
	}
	err = query.CreateTableURL(DB) // create table URL if exists
	if err != nil {
		logger.Log.Error("failed create table URL ", zap.Error(err))
		return
	}

}
