package server

import (
	"bdoip/networkAdapter"
	"bdoip/repository"
)

type Server struct {
	Adapter    networkAdapter.NetworkAdapter
	Repository repository.Repository
}

func CreateServer() (s *Server, err error) {

	return
}

func (s *Server) Start() error {

	return nil
}
