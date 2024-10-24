package sshutils

import (
	"context"

	"github.com/charmbracelet/ssh"
	gossh "golang.org/x/crypto/ssh"
)

// PublicKeyFromContext returns the public key from the context.
func PublicKeyFromContext(ctx context.Context) gossh.PublicKey {
	if pk, ok := ctx.Value(ssh.ContextKeyPublicKey).(gossh.PublicKey); ok {
		return pk
	}
	return nil
}

// ContextKeySession is the context key for the SSH session.
var ContextKeySession = &struct{ string }{"session"}

// SessionFromContext returns the SSH session from the context.
func SessionFromContext(ctx context.Context) ssh.Session {
	if s, ok := ctx.Value(ContextKeySession).(ssh.Session); ok {
		return s
	}
	return nil
}
