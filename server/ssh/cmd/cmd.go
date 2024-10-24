package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
	"text/template"
	"unicode"
)

var templateFuncs = template.FuncMap{
	"trim":                    strings.TrimSpace,
	"trimRightSpace":          trimRightSpace,
	"trimTrailingWhitespaces": trimRightSpace,
	"rpad":                    rpad,
	"gt":                      cobra.Gt,
	"eq":                      cobra.Eq,
}

// UsageFunc is a function that can be used as a cobra.Command's
// UsageFunc to render the help output.
func UsageFunc(c *cobra.Command) error {
	//ctx := c.Context()
	//cfg := config.FromContext(ctx)
	hostname := "0.0.0.0"
	port := "23"
	//url, err := url.Parse(cfg.SSH.PublicURL)
	//if err == nil {
	//	hostname = url.Hostname()
	//	port = url.Port()
	//}

	sshCmd := "ssh"
	if port != "" && port != "22" {
		sshCmd += " -p " + port
	}

	sshCmd += " " + hostname
	t := template.New("usage")
	t.Funcs(templateFuncs)
	template.Must(t.Parse(c.UsageTemplate()))
	return t.Execute(c.OutOrStderr(), struct {
		*cobra.Command
		SSHCommand string
	}{
		Command:    c,
		SSHCommand: sshCmd,
	})
}

// CommandName returns the name of the command from the args.
func CommandName(args []string) string {
	if len(args) == 0 {
		return ""
	}
	return args[0]
}

func trimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

// rpad adds padding to the right of a string.
func rpad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(template, s)
}
