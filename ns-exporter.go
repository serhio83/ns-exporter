package main

import (
	"context"
	"log"

	"ns-exporter/nsenter"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const defaultDockerAPIVersion = "v1.38"

var (
	contId []string
	pids   []int
)

func main() {
	cli, err := client.NewClientWithOpts(client.WithVersion(defaultDockerAPIVersion))
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		contId = append(contId, container.ID)
	}

	for _, id := range contId {
		stats, err := cli.ContainerInspect(context.Background(), id)
		if err != nil {
			panic(err)
		}
		pids = append(pids, stats.State.Pid)
	}

	for _, pid := range pids{
		nscfg := nsenter.Config{
			Mount:	true,
			Target: pid,
		}

		stdout, _, err := nscfg.Execute("ss -nat")
		if err != nil {
			log.Println(err)
		}
		log.Println(stdout)
	}
}
