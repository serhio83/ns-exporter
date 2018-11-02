package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const defaultDockerAPIVersion = "v1.38"

var (
	contId []string
	pids   []int
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion(defaultDockerAPIVersion))
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%+v", containers)

	for _, container := range containers {
		contId = append(contId, container.ID)
	}

	for _, id := range contId {
		stats, err := cli.ContainerInspect(ctx, id)
		if err != nil {
			panic(err)
		}
		pids = append(pids, stats.State.Pid)
	}
	fmt.Println(pids)
}
