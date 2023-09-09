package repo

import (
	"checker/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ChatGPTRepo struct {
	db *sqlx.DB
}

func NewNnChatGPTRepo(db *sqlx.DB) *ChatGPTRepo {
	return &ChatGPTRepo{db: db}
}

func (r *ChatGPTRepo) InitTable(columns map[string]string) error {
	var query string = fmt.Sprintf("CREATE TABLE %s (sentence varchar(300) NOT NULL, real_answer varchar(50) NOT NULL, current_answer varchar(50) NOT NULL", nnAnswer)
	for column := range columns {
		query = query + ", " + columns[column] + " NUMERIC(10,2) DEFAULT 0"
	}
	query = query + ")"

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (r *ChatGPTRepo) AddResultNN(add models.Storage) error {
	query := fmt.Sprintf("INSERT INTO %s (sentence, answer) VALUES ($1, $2)", nnAnswer)

	if _, err := r.db.Exec(query, add.Sentence, add.Answer); err != nil {
		return err
	}

	return nil
}

func (r *ChatGPTRepo) Query(str string) error {
	_, err := r.db.Exec(str)
	return err
}
