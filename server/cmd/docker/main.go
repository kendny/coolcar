package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"time"
)

func main() {
	c, err := client.NewEnvClient()

	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	resp, err := c.ContainerCreate(ctx, &container.Config{
		Image: "mongo:4",
		ExposedPorts: nat.PortSet{
			"27017/tcp": {},
		},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			"27017/tcp": []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: "0", //"27018", 如果设置为0，自动调空闲端口
				},
			},
		},
	}, nil, "")

	if err != nil {
		panic(err)
	}

	err = c.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}

	// 查端口
	inspRes, err := c.ContainerInspect(ctx, resp.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listening at: %v\n", inspRes.NetworkSettings.Ports["27017/tcp"])
	// 测试完将docker container 删掉
	fmt.Println("container started")
	time.Sleep(10 * time.Second)
	fmt.Println("killing container")
	err = c.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		panic(err)
	}
}
