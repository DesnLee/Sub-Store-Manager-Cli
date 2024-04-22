package docker

import (
	"context"

	"sub-store-manager-cli/lib"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var (
	dcIsInit bool
	dc       *client.Client
	dcCtx    context.Context
)

type Container struct {
	Name            string
	ImageName       string
	Version         string
	HostPort        string
	ContainerType   string
	DockerfileStr   string
	DockerContainer types.Container
	Network         string
	Private         bool
}

type PortInfo struct {
	HostIP  string
	Public  string
	Private string
	Type    string
}

func InitDockerClient() {
	if dcIsInit {
		return
	}

	// 创建 Docker 客户端
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		lib.PrintError("Failed to create Docker client:", err)
	}

	dc = dockerClient
	dcCtx = context.Background()
	dcIsInit = true
}
