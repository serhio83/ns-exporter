package main

import (
"context"
	"log"
	"regexp"
	"strings"

	"ns-exporter/pkg/nsenter"

"github.com/docker/docker/api/types"
"github.com/docker/docker/client"
)

const defaultDockerAPIVersion = "v1.35"

type Container struct {
	Name    string
	Pid     int
	Sockets []string
}

type Containers struct {
	List []Container
}

var re = regexp.MustCompile(`ESTAB.*`)

func getContainers() Containers {
	// init docker client
	cli, err := client.NewClientWithOpts(client.WithVersion(defaultDockerAPIVersion))
	if err != nil {
		log.Println(err)
	}

	// get list of docker containers
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
			Name: 		container.Names[0][1:],
			Pid:  		pid.State.Pid,
			Sockets:	getSockets(pid.State.Pid),
		}
		myList = append(myList, cont)
	}
	return Containers{List: myList}
}


func getSockets(pid int) []string {
	// configure to enter network namespace (Net: true)
	nscfg := nsenter.Config{
		Net:	true,
		Target: pid,
	}

	// get sockets from container namespace
	stdout, _, err := nscfg.Execute("ss", "-nat")
	if err != nil {
		log.Println(err)
	}

	// find all needed sockets in ESTAB state
	str := re.FindAllString(stdout, len(stdout))
	socks := make([]string, 0)
	for _, s := range str{
		result := strings.Fields(s)
		socks = append(socks, result[4])
	}
	return socks
}
