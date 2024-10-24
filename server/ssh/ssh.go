package ssh

import (
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"net"
)

type SSHServer struct {
	srv *ssh.Server
}

func NerSSHServer() (*SSHServer, error) {
	var err error
	s := &SSHServer{}

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
