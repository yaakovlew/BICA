package service

import (
	"checker/models"
	"testing"
)

func TestParseString(t *testing.T) {
	str := "грусть - 0.904\nзлость - 0.347\nрадость - 0.008\nотсутствие эмоций - 0.005"
	rightResult := []models.ToAddInDB{
		models.ToAddInDB{
			Feeling: "грусть",
			Rate:    "0.904",
		},
		models.ToAddInDB{
			Feeling: "злость",
			Rate:    "0.347",
		},
		models.ToAddInDB{
			Feeling: "радость",
			Rate:    "0.008",
		},
		models.ToAddInDB{
			Feeling: "отсутствие эмоций",
			Rate:    "0.005",
		},
	}
	res := ParseString(str)
	for i, _ := range rightResult {
		if rightResult[i] != res[i] {
			t.Errorf("Not correct data")
		}
	}
}
