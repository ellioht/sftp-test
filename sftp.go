package sftptest

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"time"
)

type SftpClient struct {
	Client    *sftp.Client
	sshClient *ssh.Client
}

func NewSftpClient(cfg *ssh.ClientConfig) (*SftpClient, error) {
	client, err := ssh.Dial("tcp", "localhost:22", cfg)
	if err != nil {
		return nil, err
	}

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		if err = client.Close(); err != nil {
			return nil, err
		}
		return nil, err
	}

	return &SftpClient{
		Client: sftpClient,
	}, nil
}

func (s *SftpClient) Close() error {
	if err := s.Client.Close(); err != nil {
		return err
	}
	if err := s.sshClient.Close(); err != nil {
		return err
	}
	return nil
}

func CreateDefaultSSHConfig() *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: "foo",
		Auth: []ssh.AuthMethod{
			ssh.Password("pass"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         20 * time.Second,
	}
}
