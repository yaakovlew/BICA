package repo

import (
	"checker/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type nnRepo struct {
	db *sqlx.DB
}

func NewNnRepo(db *sqlx.DB) *nnRepo {
	return &nnRepo{db: db}
}

func (r *nnRepo) AddResultNN(add models.Storage) error {
	query := fmt.Sprintf("INSERT INTO %s (sentence, answer) VALUES ($1, $2)", nnAnswer)

	if _, err := r.db.Exec(query, add.Sentence, add.Answer); err != nil {
		return err
	}

	return nil
}
