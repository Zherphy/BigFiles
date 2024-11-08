package cmd

import (
	"bytes"
	"context"
	"errors"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/wish/git"
	"github.com/metalogical/BigFiles/server/ssh/sshutils"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"io"
	"os/exec"
	"time"
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

// AuthenticateResponse is the git-lfs-authenticate JSON response object.
type AuthenticateResponse struct {
	Header    map[string]string `json:"header"`
	Href      string            `json:"href"`
	ExpiresIn time.Duration     `json:"expires_in"`
	ExpiresAt time.Time         `json:"expires_at"`
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
	pk1 := sshutils.PublicKeyFromContext(ctx)
	log.Printf("处理git-lfs-authenticate /n")
	log.Printf(pk1.Type())
	//给gitee发送ssh认证
	cfg := &ssh.ClientConfig{
		User: "git",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(),
		},
		HostKeyCallback: ssh.FixedHostKey(pk1),
	}
	client, err := ssh.Dial("tcp", "gitee.com:22", cfg)
	if err != nil {
		log.Printf("failed to dial: %v", err)
	}
	defer func(client *ssh.Client) {
		err := client.Close()
		if err != nil {
			log.Printf("failed to close ssh client: %v", err)
		}
	}(client)

	//创建一个SSH会话
	session, err := client.NewSession()
	if err != nil {
		log.Printf("failed to create session: %v", err)
	}
	defer func(session *ssh.Session) {
		err := session.Close()
		if err != nil {
			log.Printf("failed to close session: %v", err)
		}
	}(session)

	// 设置会话的标准输出和标准错误输出为字节缓冲区，以便获取命令执行结果
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	// 执行指定的指令
	err = session.Run("git-lfs-authenticate wj00037/lfs-test.git download")
	if err != nil {
		log.Printf("命令执行出错: %v", err)
		log.Printf("标准错误输出: %s", stderrBuf.String())
	}
	// 输出命令执行的标准输出结果
	log.Printf("标准输出: %s", stdoutBuf.String())
	return nil
}
