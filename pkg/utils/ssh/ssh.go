package ssh

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	remoteAddr string
	username   string
	sshSigner  ssh.Signer
}

func NewSSHClient(remoteAddr string, username string, sshKeyFile string) (*SSHClient, error) {
	if strings.HasPrefix(sshKeyFile, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		sshKeyFile = strings.Replace(sshKeyFile, "~", homeDir, 1)
	}

	sshKeyData, err := os.ReadFile(sshKeyFile)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(sshKeyData)
	if err != nil {
		return nil, err
	}

	s := &SSHClient{
		remoteAddr: remoteAddr,
		username:   username,
		sshSigner:  signer,
	}

	return s, nil
}

func (s *SSHClient) RemoveAddr() string {
	return s.remoteAddr
}

// AuthorizedKey returns the public key in authorized_keys format
func (s *SSHClient) AuthorizedKey() string {
	return string(ssh.MarshalAuthorizedKey(s.sshSigner.PublicKey()))
}

// Exec executes a command and returns stdout or not nil error
func (s *SSHClient) Exec(ctx context.Context, cmd string) (string, error) {
	stdout, stderr, exitCode, err := s.exec(ctx, cmd)
	if err != nil || exitCode != 0 {
		return "", fmt.Errorf("command %s exit %d stderr %s", cmd, exitCode, stderr)
	}

	return stdout, nil
}

// Run runs commands in sequence, if any command fails, it will return error
func (s *SSHClient) Run(ctx context.Context, cmds ...string) error {
	for _, cmd := range cmds {
		_, stderr, exitCode, err := s.exec(ctx, cmd)
		if err != nil {
			return fmt.Errorf("command %s error: %w", cmd, err)
		}
		if exitCode != 0 {
			return fmt.Errorf("command %s exit %d stderr %s", cmd, exitCode, stderr)
		}
	}

	return nil
}

func (s *SSHClient) exec(_ context.Context, cmd string) (stdout string, stderr string, exitCode int, err error) {
	sshConn, err := s.newConn()
	if err != nil {
		return "", "", -1, err
	}
	sess, err := sshConn.NewSession()
	if err != nil {
		return "", "", -1, err
	}

	stdoutBuf := bytes.NewBuffer(nil)
	stderrBuf := bytes.NewBuffer(nil)
	sess.Stdout, sess.Stderr = stdoutBuf, stderrBuf
	cmd = fmt.Sprintf(`sh --login -c "%s"`, cmd)
	err = sess.Start(cmd)
	if err != nil {
		return "", "", -1, err
	}

	err = sess.Wait()
	if err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
			exitCode = exitErr.ExitStatus()
		} else {
			return "", "", -1, err
		}
	}

	return stdoutBuf.String(), stderrBuf.String(), exitCode, nil
}

func (s *SSHClient) newConn() (*ssh.Client, error) {
	return ssh.Dial("tcp", s.remoteAddr, &ssh.ClientConfig{
		User: s.username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(s.sshSigner),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
}

func (s *SSHClient) WriteFile(content io.Reader, dstFilePath string) error {
	sshConn, err := s.newConn()
	if err != nil {
		return fmt.Errorf("%s@%s failed to create ssh connection: %w", s.username, s.remoteAddr, err)
	}
	defer sshConn.Close()

	sftpClient, err := sftp.NewClient(sshConn)
	if err != nil {
		return fmt.Errorf("%s@%s failed to create sftp client: %w", s.username, s.remoteAddr, err)
	}
	defer sftpClient.Close()

	file, err := sftpClient.Create(dstFilePath)
	if err != nil {
		return fmt.Errorf("%s@%s failed to create file %s: %w", s.username, s.remoteAddr, dstFilePath, err)
	}
	defer file.Close()

	_, err = file.ReadFrom(content)
	if err != nil {
		return fmt.Errorf("%s@%s failed to write file %s: %w", s.username, s.remoteAddr, dstFilePath, err)
	}
	return nil
}
