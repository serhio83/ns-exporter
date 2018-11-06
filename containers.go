package main

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const defaultDockerAPIVersion = "v1.35"

type Container struct {
	Name    string
	Pid     int
}

type Containers struct {
	list []Container
}

func getContainers() Containers {
	cli, err := client.NewClientWithOpts(client.WithVersion(defaultDockerAPIVersion))
	if err != nil {
		log.Println(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Println(err)
	}
	var myList []Container

	for _, container := range containers {
		pid, err := cli.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			log.Println(err)
		}
		cont := Container{
			Name: container.Names[0][1:],
			Pid:  pid.State.Pid,
		}
		myList = append(myList, cont)
	}
	return Containers{list: myList}
}

//nscfg := nsenter.Config{
//Mount:	true,
//Target: pid,
//}
//
//stdout, _, err := nscfg.Execute("ss -nat")
//if err != nil {
//log.Println(err)
//}
//log.Println(stdout)