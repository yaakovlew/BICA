package service

import "checker/repo"

type NNChatGPT interface {
	SendTOChatGPT(str, answer string) error
	ParseCSVFile(path string)
	//SendTONN(string) error
}

type Service struct {
	NNChatGPT
}

func NewService(repo *repo.Repository) *Service {
	return &Service{
		NNChatGPT: NewNNChatGPTService(repo),
	}
}
