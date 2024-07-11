package app

import (
	"fmt"
	"os"

	"github.com/jasonli0226/ssh-connection-manager/pkg/domain"
	"golang.org/x/crypto/ssh"
)

type SSHManager struct {
	repo domain.ConnectionRepository
}

func NewSSHManager(repo domain.ConnectionRepository) *SSHManager {
	return &SSHManager{repo: repo}
}

func (m *SSHManager) AddConnection(alias, host, user, password string, port int) error {
	if alias == "" || host == "" || user == "" || port <= 0 || port > 65535 {
		return NewAppError("Invalid connection details", nil)
	}

	conn := domain.Connection{
		Alias:    alias,
		Host:     host,
		User:     user,
		Password: password,
		Port:     port,
	}

	if err := m.repo.Add(conn); err != nil {
		return NewAppError("Failed to add connection", err)
	}

	return nil
}

func (m *SSHManager) ListConnections() ([]domain.Connection, error) {
	conns, err := m.repo.List()
	if err != nil {
		return nil, NewAppError("Failed to list connections", err)
	}
	return conns, nil
}

func (m *SSHManager) Connect(alias string) error {
	conn, err := m.repo.Get(alias)
	if err != nil {
		return NewAppError(fmt.Sprintf("Connection with alias '%s' not found", alias), err)
	}

	config := &ssh.ClientConfig{
		User: conn.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(conn.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", conn.Host, conn.Port), config)
	if err != nil {
		return NewAppError("Failed to establish SSH connection", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return NewAppError("Failed to create SSH session", err)
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return NewAppError("Failed to request PTY", err)
	}

	if err := session.Shell(); err != nil {
		return NewAppError("Failed to start shell", err)
	}

	if err := session.Wait(); err != nil {
		return NewAppError("SSH session ended with error", err)
	}

	return nil
}

func (m *SSHManager) DeleteConnection(alias string) error {
	if alias == "" {
		return NewAppError("Invalid alias", nil)
	}

	err := m.repo.Delete(alias)
	if err != nil {
		return NewAppError(fmt.Sprintf("Failed to delete connection with alias '%s'", alias), err)
	}

	return nil
}
