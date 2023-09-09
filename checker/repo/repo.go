package repo

import (
	"checker/models"
	"github.com/jmoiron/sqlx"
)

type NNChatGPT interface {
	AddResultNN(add models.Storage) error
	InitTable(columns map[string]string) error
	Query(str string) error
}

type Repository struct {
	NNChatGPT
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		NNChatGPT: NewNnChatGPTRepo(db),
	}
}
