package repository

import (
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/url"
	"os"
	"strconv"
)

// настройки соединения
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DbName   string
	Timeout  int
}

// функция для создания строки подключения
func NewPoolConfig() (*pgxpool.Config, error) {
	config := &Config{
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("PORT"),
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
		DbName:   os.Getenv("DB_NAME"),
	}
	config.Timeout, _ = strconv.Atoi(os.Getenv("TIMEOUT"))

	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=%d",
		"postgres",
		url.QueryEscape(config.Username),
		url.QueryEscape(config.Password),
		config.Host,
		config.Port,
		config.DbName,
		config.Timeout)

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}
	return poolConfig, nil
}
