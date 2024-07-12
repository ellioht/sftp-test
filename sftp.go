package sftptest

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"time"
)

func CreateSftpClient() (*sftp.Client, *ssh.Client, error) {
	cfg := &ssh.ClientConfig{
		User: "foo",
		Auth: []ssh.AuthMethod{
			ssh.Password("pass"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         20 * time.Second,
	}

	client, err := ssh.Dial("tcp", "localhost:22", cfg)
	if err != nil {
		return nil, nil, err
	}

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		if err = client.Close(); err != nil {
			return nil, nil, err
		}
		return nil, nil, err
	}

	return sftpClient, client, nil
}
