package main

import (
	"fmt"

	"github.com/190930-UTA-CW-Go/project2-AGDJ/ssh"
)

func main() {
	fmt.Println(ssh.CmdGetInfo(login, ip))
	//ssh.SetupDocker(login, password, ip)
	fmt.Println(ssh.DockerStatus(login, ip))
	fmt.Println(ssh.ListContainers(login, password, ip))
	fmt.Println(ssh.ListImages(login, password, ip))
	//fmt.Println(ssh.TestRun(login, password, ip))
}
