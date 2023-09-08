package checker

import (
	"checker/service"
	"context"
)

type Server struct {
	service service.Service
}

var TurnOnTumbler bool = true

func (s *Server) Run(service service.Service) error {
	for TurnOnTumbler {

	}
	return nil
}

func (s *Server) ShutDown(ctx context.Context) error {
	TurnOnTumbler = false
	return nil
}
