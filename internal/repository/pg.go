package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Postgres struct {
}

func NewPostgresRepo(pool *pgxpool.Pool) *Postgres {
	return &Postgres{}
}
