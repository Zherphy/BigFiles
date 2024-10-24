package ssh

import (
	"github.com/charmbracelet/ssh"
	"github.com/metalogical/BigFiles/server/ssh/cmd"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/spf13/cobra"
)

var cliCommandCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "lfs_serve",
	Subsystem: "cli",
	Name:      "commands_total",
	Help:      "Total times each command was called",
}, []string{"command"})

// CommandMiddleware handles git commands and CLI commands.
// This middleware must be run after the ContextMiddleware.
func CommandMiddleware(sh ssh.Handler) ssh.Handler {
	return func(s ssh.Session) {
		_, _, ptyReq := s.Pty()
		if ptyReq {
			sh(s)
			return
		}

		ctx := s.Context()
		args := s.Command()
		cliCommandCounter.WithLabelValues(cmd.CommandName(args)).Inc()
		rootCmd := &cobra.Command{
			Short:        "lfs Serve is a self-hostable Git server for the command line.",
			SilenceUsage: true,
		}
		rootCmd.CompletionOptions.DisableDefaultCmd = true

		rootCmd.SetUsageFunc(cmd.UsageFunc)
		rootCmd.AddCommand(
			cmd.GitLFSAuthenticateCommand(),
		)

		rootCmd.SetArgs(args)
		if len(args) == 0 {
			// otherwise it'll default to os.Args, which is not what we want.
			rootCmd.SetArgs([]string{"--help"})
		}
		rootCmd.SetIn(s)
		rootCmd.SetOut(s)
		rootCmd.SetErr(s.Stderr())
		rootCmd.SetContext(ctx)

		if err := rootCmd.ExecuteContext(ctx); err != nil {
			s.Exit(1) // nolint: errcheck
			return
		}
	}
}
