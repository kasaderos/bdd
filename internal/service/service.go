package service

import (
	"bdd/config"
	"bdd/internal/repository"
)

type Service struct {
}

func New(repo *repository.Postgres, conf config.ServiceConfig) *Service {
	return &Service{}
}
