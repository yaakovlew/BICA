package repo

import (
	"checker/models"
	"github.com/jmoiron/sqlx"
)

type NN interface {
	AddResultNN(add models.Storage) error
}

type ChatGPT interface {
	AddResultChatGPT(add models.Storage) error
}

type Repository struct {
	NN
	ChatGPT
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		NN:      NewNnRepo(db),
		ChatGPT: NewChatGPTRepo(db),
	}
}
