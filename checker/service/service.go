package service

import "checker/repo"

type NNChatGPT interface {
	SendTOChatGPT(string) error
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
