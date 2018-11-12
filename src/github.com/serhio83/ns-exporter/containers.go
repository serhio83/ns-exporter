package main

import (
	"context"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/serhio83/ns-exporter/nsenter"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const defaultDockerAPIVersion = "v1.35"

type sock map[string]int

type Container struct {
	Name    string
	Pid     int
	Sockets []sock
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
			Name:    container.Names[0][1:],
			Pid:     pid.State.Pid,
			Sockets: getSockets(pid.State.Pid),
		}
		myList = append(myList, cont)
	}
	return Containers{List: myList}
}

func getSockets(pid int) []sock {
	// configure to enter network namespace (Net: true)
	nscfg := nsenter.Config{
		Net:    true,
		Target: pid,
	}

	// get sockets from container namespace
	stdout, _, _ := nscfg.Execute("ss", "-nat")

	// find all needed sockets in ESTAB state by regex compiled above
	str := re.FindAllString(stdout, len(stdout))
	var socks []sock

	for _, s := range str {
		result := strings.Fields(s)
		mySocket := strings.Split(result[4], ":")
		var skt = make(map[string]int)
		port, _ := strconv.Atoi(mySocket[1])
		skt[mySocket[0]] = port
		socks = append(socks, skt)
	}
	return socks
}
