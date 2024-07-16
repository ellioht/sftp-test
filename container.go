package sftptest

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"net"
	"os"
	"path/filepath"
)

const (
	AtmozSftpImage = "atmoz/sftp:latest"
)

type Config struct {
	ImageName string
	MountDir  string
}

type Container struct {
	Id    string
	Port  string
	close func() error
}

func NewContainer(ctx context.Context, cfg *Config) (*Container, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	hostDir := filepath.Join(dir, cfg.MountDir)

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	img, err := cli.ImagePull(ctx, cfg.ImageName, image.PullOptions{})
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(os.Stdout, img)
	if err != nil {
		return nil, err
	}

	containerCfg := &container.Config{
		Image: cfg.ImageName,
		Cmd:   []string{"foo:pass:1001"},
	}

	hostPort, err := findAvailablePort()
	if err != nil {
		return nil, err
	}

	containerHstCfg := &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{
			"22/tcp": {
				{
					HostIP:   "0.0.0.0",
					HostPort: hostPort,
				},
			},
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: hostDir,
				Target: "/home/foo/upload",
			},
		},
	}

	res, err := cli.ContainerCreate(ctx, containerCfg, containerHstCfg, nil, nil, "")
	if err != nil {
		return nil, err
	}

	if err = cli.ContainerStart(ctx, res.ID, container.StartOptions{}); err != nil {
		return nil, err
	}

	cleanup := func() error {
		if err = cli.ContainerStop(ctx, res.ID, container.StopOptions{}); err != nil {
			return err
		}

		if err = cli.ContainerRemove(ctx, res.ID, container.RemoveOptions{}); err != nil {
			return err
		}

		err = img.Close()
		if err != nil {
			return err
		}

		return nil
	}

	return &Container{
		Id:    res.ID,
		Port:  hostPort,
		close: cleanup,
	}, nil
}

func (c *Container) Close() error {
	return c.close()
}

func findAvailablePort() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "", err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "", err
	}
	defer l.Close()

	return fmt.Sprintf("%d", l.Addr().(*net.TCPAddr).Port), nil
}
