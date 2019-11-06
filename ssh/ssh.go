package ssh

import (
	"fmt"
	"os/exec"
)

// Command =
func Command(cmd string, login string, ip string) string {
	results, err := exec.Command("ssh", login+"@"+ip, cmd).Output()
	if err != nil {
		fmt.Println("ERROR: Invalid Command")
	}
	return string(results)
}

// CmdGetInfo =
func CmdGetInfo(login string, ip string) string {
	currentUser := Command("whoami", login, ip)
	currentUser = currentUser[:len(currentUser)-1]
	hostName := Command("hostname", login, ip)

	return currentUser + "@" + string(hostName)
}

// SetupDocker =
func SetupDocker(login string, password string, ip string) {
	fmt.Println(InstallDocker(login, password, ip))
	fmt.Println(StartDocker(login, password, ip))
	fmt.Println(EnableDocker(login, password, ip))
}

// InstallDocker =
func InstallDocker(login string, password string, ip string) string {
	return Command("echo "+password+" | sudo -S apt install docker.io -y", login, ip)
}

// StartDocker =
func StartDocker(login string, password string, ip string) string {
	return Command("echo "+password+" | sudo -S systemctl start docker", login, ip)
}

// EnableDocker =
func EnableDocker(login string, password string, ip string) string {
	return Command("echo "+password+" | sudo -S systemctl enable docker", login, ip)
}
