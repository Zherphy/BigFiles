package ssh

import (
	"context"
	"errors"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
)

var (
	publicKeyCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "soft_serve",
		Subsystem: "ssh",
		Name:      "public_key_auth_total",
		Help:      "The total number of public key auth requests",
	}, []string{"allowed"})
)

type SSHServer struct {
	srv *ssh.Server
	ctx context.Context
}

func NerSSHServer(ctx context.Context) (*SSHServer, error) {
	var err error
	s := &SSHServer{
		ctx: ctx,
	}

	mw := []wish.Middleware{
		CommandMiddleware,
	}

	opts := []ssh.Option{
		ssh.PublicKeyAuth(s.PublicKeyHandler),
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
	s.srv.Addr = ":23231"
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

// PublicKeyAuthHandler handles public key authentication.
func (s *SSHServer) PublicKeyHandler(ctx ssh.Context, pk ssh.PublicKey) (allowed bool) {
	if pk == nil {
		return false
	}

	return true
}
