package repo

import (
	"checker/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ChatGPTRepo struct {
	db *sqlx.DB
}

func NewChatGPTRepo(db *sqlx.DB) *ChatGPTRepo {
	return &ChatGPTRepo{db: db}
}

func (r *ChatGPTRepo) AddResultChatGPT(add models.Storage) error {
	query := fmt.Sprintf("INSERT INTO %s (sentence, answer) VALUES ($1, $2)", chatGPTAnswer)

	if _, err := r.db.Exec(query, add.Sentence, add.Answer); err != nil {
		return err
	}

	return nil
}
