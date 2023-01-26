package repository

import "github.com/jackc/pgx/v4/pgxpool"

type Instance struct {
	Db *pgxpool.Pool
}
