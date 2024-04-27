package cmd

import (
	"fmt"

	"sub-store-manager-cli/docker"
	"sub-store-manager-cli/lib"
	"sub-store-manager-cli/vars"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update a sub-store docker container",
	Run: func(cmd *cobra.Command, args []string) {
		updateContainer()
	},
}

func init() {
	updateCmd.Flags().StringVarP(&inputVersion, "version", "v", "", "The target version to update")
	updateCmd.Flags().StringVarP(&inputName, "name", "n", "", "The target sub-store container name to update")
}

func updateContainer() {
	name := inputName
	if name == "" {
		name = vars.DockerNameBE
	}

	oldContainer, isExist := docker.GetContainerByName(name)
	if !isExist {
		lib.PrintError("The container does not exist.", nil)
	}

	c := docker.Container{
		Name:          name,
		ImageName:     oldContainer.ImageName,
		ContainerType: oldContainer.ContainerType,
		HostPort:      oldContainer.HostPort,
		Version:       inputVersion,
	}

	// 检查指定版本
	if inputVersion == "" {
		c.SetLatestVersion()
		fmt.Println("No version specified, using the latest version")
	} else {
		isValid := c.CheckVersionValid()
		if !isValid {
			lib.PrintError("The version is not valid.", nil)
		}
	}

	// 如果当前运行的版本就是目标版本，则不更新
	if oldContainer.Version == c.Version {
		lib.PrintInfo("The current version is the target version, no need to update.")
		return
	}

	// 获取旧容器的端口信息
	if p, err := oldContainer.GetPortInfo(); err != nil {
		lib.PrintError("Failed to get port info:", err)
	} else {
		c.HostPort = p.Public
		if p.HostIP == "127.0.0.1" {
			c.Private = true
		} else {
			c.Private = false
		}
	}

	c.SetDockerfile("")
	c.CreateImage()

	// 删除旧容器, 启动新容器
	oldContainer.Stop()
	oldContainer.Delete()
	oldContainer.DeleteImage()
	c.StartImage()
}
