package cmd

import (
	"context"
	"errors"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/wish/git"
	"github.com/spf13/cobra"
	"io"
	"os/exec"
)

type Service string

// ServiceCommand is used to run a git service command.
type ServiceCommand struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Args   []string

	// Modifier functions
	CmdFunc func(*exec.Cmd)
}

const (
	// LFSAuthenticateService is the LFS authenticate service.
	LFSAuthenticateService = "git-lfs-authenticate"
	// OperationDownload is the operation name for a download request.
	OperationDownload = "download"

	// OperationUpload is the operation name for an upload request.
	OperationUpload = "upload"
)

// GitLFSAuthenticateCommand returns a cobra command for git-lfs-authenticate.
func GitLFSAuthenticateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "git-lfs-authenticate REPO OPERATION",
		Short:  "Git LFS authenticate",
		Args:   cobra.ExactArgs(2),
		Hidden: true,
		RunE:   gitRunE,
	}

	return cmd
}

func gitRunE(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	logger := log.FromContext(ctx)
	service := Service(cmd.Name())
	stdin := cmd.InOrStdin()
	stdout := cmd.OutOrStdout()
	stderr := cmd.ErrOrStderr()
	scmd := ServiceCommand{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
		Args:   args,
	}

	switch service {
	case LFSAuthenticateService:
		if err := service.Handler(ctx, scmd); err != nil {
			logger.Error("failed to handle lfs service", "service", service, "err", err, "repo")
			return git.ErrSystemMalfunction
		}

		return nil
	}

	return errors.New("unsupported git service")
}

func (s Service) Handler(ctx context.Context, cmd ServiceCommand) error {
	return nil
}
