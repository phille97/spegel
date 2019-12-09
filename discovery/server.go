package discovery

import (
	"context"

	"github.com/godbus/dbus/v5"
	"github.com/holoplot/go-avahi"
)

type Server struct {
	service string
	port    uint16
}

func NewServer(service string, port uint16) (*Server, error) {
	return &Server{
		service: service,
		port:    port,
	}, nil
}

func (s *Server) Register(ctx context.Context) error {
	conn, err := dbus.SystemBus()
	if err != nil {
		return err
	}

	a, err := avahi.ServerNew(conn)
	if err != nil {
		return err
	}

	eg, err := a.EntryGroupNew()
	if err != nil {
		return err
	}

	hostname, err := a.GetHostName()
	if err != nil {
		return err
	}

	fqdn, err := a.GetHostNameFqdn()
	if err != nil {
		return err
	}

	err = eg.AddService(avahi.InterfaceUnspec, avahi.ProtoUnspec, 0, hostname, s.service, "local", fqdn, s.port, nil)
	if err != nil {
		return err
	}

	err = eg.Commit()
	if err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}
