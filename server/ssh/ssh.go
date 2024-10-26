package ssh

import (
	"context"
	"errors"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
)

type SSHServer struct {
	srv *ssh.Server
	ctx context.Context
}

func NerSSHServer() (*SSHServer, error) {
	var err error
	ctx := context.Background()
	s := &SSHServer{
		ctx: ctx,
	}

	mw := []wish.Middleware{
		CommandMiddleware,
	}

	opts := []ssh.Option{
		wish.WithMiddleware(mw...),
	}
	s.srv, err = wish.NewServer(opts...)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// ListenAndServe starts the SSH server.
func (s *SSHServer) ListenAndServe() error {
	return s.srv.ListenAndServe()
}

// Serve starts the SSH server on the given net.Listener.
func (s *SSHServer) Serve(l net.Listener) error {
	return s.srv.Serve(l)
}

// Close closes the SSH server.
func (s *SSHServer) Close() error {
	return s.srv.Close()
}

func (s *SSHServer) Start() error {
	errg, _ := errgroup.WithContext(s.ctx)
	errg.Go(func() error {
		log.Println("ssh server on 0.0.0.0:22 ...")
		if err := s.ListenAndServe(); !errors.Is(err, ssh.ErrServerClosed) {
			return err
		}
		return nil
	})
	return nil
}
