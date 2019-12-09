package discovery

import (
	"github.com/grandcat/zeroconf"
	"os"
)

type Server struct {
	zserver *zeroconf.Server

	hostname string
	service  string
	port     int
}

func NewServer(service string, port int) (*Server, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	return &Server{
		hostname: hostname,
		service:  service,
		port:     port,
	}, nil
}

func (s *Server) Register() error {
	zserver, err := zeroconf.Register(s.hostname, s.service, "local.", s.port, []string{}, nil)
	if err != nil {
		return err
	}
	s.zserver = zserver
	return nil
}

func (s *Server) Shutdown() {
	if s.zserver != nil {
		s.zserver.Shutdown()
		s.zserver = nil
	}
}
