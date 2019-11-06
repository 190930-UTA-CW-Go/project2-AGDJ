package main

import (
	"fmt"

	"github.com/190930-UTA-CW-Go/project2-AGDJ/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	//login, password, ip := login()

	login := "agent1"
	password := "secretagent"
	ip := "192.168.56.102"
	fmt.Println(ssh.CmdGetInfo(login, ip))
	ssh.Install(login, password, ip)

}

func login() (login string, password string, ip string) {
	fmt.Print("Login: ")
	fmt.Scan(&login)
	fmt.Print("Password: ")
	result, _ := terminal.ReadPassword(0)
	password = string(result)
	fmt.Println()
	fmt.Print("IP Address: ")
	fmt.Scan(&ip)
	fmt.Println()

	return
}
