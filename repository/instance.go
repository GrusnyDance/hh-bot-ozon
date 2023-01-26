package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
	"strings"
	"time"
)

type Instance struct {
	Db *pgxpool.Pool
}

type Vacancy struct {
	Id      int
	Created time.Time
	Name    string
	Link    string
}

func (i *Instance) GetVac(mapVac *map[string]string) (*map[string]string, bool) {
	mapRet := make(map[string]string)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*2))
	defer cancel()

	for key, val := range *mapVac {
		row := i.Db.QueryRow(ctx, "SELECT id FROM vacancies WHERE name = $1;", key)
		// Query возвращает структуру pgx.Row

		vac := &Vacancy{}
		if err := row.Scan(&vac.Id); err == pgx.ErrNoRows {
			mapRet[key] = val
			i.InsertVac(key, val)
		}
	}

	if len(mapRet) > 0 {
		return &mapRet, true
	} else {
		return nil, false
	}
}

func (i *Instance) InsertVac(name string, url string) {
	fmt.Println("I am trying to insert")
	trimFactor := os.Getenv("OZON_TRIM")
	url = strings.TrimLeft(url, trimFactor)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*2))
	defer cancel()

	_, err := i.Db.Exec(ctx, "INSERT INTO vacancies (created_at, name, link) VALUES ($1, $2, $3);",
		time.Now(), name, url)
	if err != nil {
		fmt.Println(err)
		fmt.Println(name, url)
	}
}
