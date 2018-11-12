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
	name         string
	pid          int
	tcpConnEstab []sock
	tcpConnTW    []sock
}

type Containers struct {
	contList []Container
}

var reConnEst = regexp.MustCompile(`ESTAB.*`)
var reConnTW = regexp.MustCompile(`TIME-WAIT.*`)

// connects to docker API & gets containers names
func (conts *Containers) getContainers() {
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
			name: container.Names[0][1:],
			pid:  pid.State.Pid,
		}
		cont.getTcpConn()
		myList = append(myList, cont)
	}

	conts.contList = myList
}

// collects tcp connections inside network namespace of docker container
func (cont *Container) getTcpConn() {
	// configure to enter network namespace (Net: true)
	nscfg := nsenter.Config{
		Net:    true,
		Target: cont.pid,
	}

	// get sockets from container namespace
	stdout, _, _ := nscfg.Execute("ss", "-4nat")

	// find all needed sockets in ESTAB & TIME-WAIT state by regexps compiled above
	connEst := reConnEst.FindAllString(stdout, len(stdout))
	connTW := reConnTW.FindAllString(stdout, len(stdout))
	cont.tcpConnEstab = parseSocketInfo(connEst)
	cont.tcpConnTW = parseSocketInfo(connTW)
}

// parses stdout & collects IPs + ports
func parseSocketInfo(connList []string) []sock {
	var socks []sock

	for _, s := range connList {
		result := strings.Fields(s)
		mySocket := strings.Split(result[4], ":")
		var skt = make(map[string]int)
		port, _ := strconv.Atoi(mySocket[1])
		skt[mySocket[0]] = port
		socks = append(socks, skt)
	}

	return socks
}
