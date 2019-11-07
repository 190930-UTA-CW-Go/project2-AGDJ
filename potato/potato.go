package potato

import (
	"fmt"
	"io"
	"log"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

// Commands =
func Commands(in chan<- string, out <-chan string, password string) {
	GetInfo(in, out)
	//SetupDocker(in, out, password)
	DockerStatus(in, out)
	ListContainers(in, out, password)
	ListImages(in, out, password)
}

// GetInfo =
func GetInfo(in chan<- string, out <-chan string) {
	// Get current user and slice off '\n'
	in <- "whoami"
	slice := strings.Fields(<-out)

	// Get host name
	in <- "hostname"
	hostName := <-out

	// Concatenate to resemble Linux
	fmt.Println(slice[0] + "@" + hostName)
}

// SetupDocker =
func SetupDocker(in chan<- string, out <-chan string, password string) {
	InstallDocker(in, out, password)
	StartDocker(in, out, password)
	EnableDocker(in, out, password)
	PullAlpine(in, out, password)
}

// InstallDocker =
func InstallDocker(in chan<- string, out <-chan string, password string) {
	in <- "echo " + password + " | sudo -S apt install docker.io -y"
	fmt.Println(<-out)
}

// StartDocker =
func StartDocker(in chan<- string, out <-chan string, password string) {
	in <- "echo " + password + " | sudo -S systemctl start docker"
	fmt.Println(<-out)
}

// EnableDocker =
func EnableDocker(in chan<- string, out <-chan string, password string) {
	in <- "echo " + password + " | sudo -S systemctl enable docker"
	fmt.Println(<-out)
}

// PullAlpine =
func PullAlpine(in chan<- string, out <-chan string, password string) {
	in <- "echo " + password + " | sudo -S docker pull alpine"
	fmt.Println(<-out)
}

// DockerStatus =
func DockerStatus(in chan<- string, out <-chan string) {
	in <- "systemctl is-active docker"
	fmt.Println(<-out)
}

// ListContainers =
func ListContainers(in chan<- string, out <-chan string, password string) {
	in <- "echo " + password + " | sudo -S docker container ls -a"
	fmt.Println(<-out)
}

// ListImages =
func ListImages(in chan<- string, out <-chan string, password string) {
	in <- "echo " + password + " | sudo -S docker image ls -a"
	fmt.Println(<-out)
}

/*
// TestRun =
func TestRun(login string, password string, ip string) string {
	return Command("echo "+password+" | sudo -S docker run alpine ls -l", login, ip)
}
*/

// Connect =
func Connect(login string, password string, ip string, port string) {
	// Connect to ssh cient
	config := &ssh.ClientConfig{
		//To resolve "Failed to dial: ssh: must specify HostKeyCallback" error
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		User:            login,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}
	client, err := ssh.Dial("tcp", ip+":"+port, config)
	if err != nil {
		fmt.Println("Invalid Login or Password")
	} else {

		session, err := client.NewSession()
		if err != nil {
			panic("Session Failed: " + err.Error())
		}
		defer session.Close()

		modes := ssh.TerminalModes{
			ssh.ECHO:          0,     // disable echoing
			ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
			ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
		}

		if err := session.RequestPty("xterm", 80, 100, modes); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Successfully Connected!")
		fmt.Println()

		w, err := session.StdinPipe()
		if err != nil {
			panic(err)
		}
		r, err := session.StdoutPipe()
		if err != nil {
			panic(err)
		}
		in, out := MuxShell(w, r)
		if err := session.Start("/bin/sh"); err != nil {
			log.Fatal(err)
		}

		// Ignore the shell output
		<-out

		////////////////////////////
		Commands(in, out, password)
		////////////////////////////

		// automatically "exit"
		in <- "exit"
		session.Wait()
	}
}

// MuxShell =
func MuxShell(w io.Writer, r io.Reader) (chan<- string, <-chan string) {
	in := make(chan string, 1)
	out := make(chan string, 1)
	var wg sync.WaitGroup
	wg.Add(1) //for the shell itself
	go func() {
		for cmd := range in {
			wg.Add(1)
			w.Write([]byte(cmd + "\n"))
			wg.Wait()
		}
	}()
	go func() {
		var (
			buf [65 * 1024]byte
			t   int
		)
		for {
			n, err := r.Read(buf[t:])
			if err != nil {
				close(in)
				close(out)
				return
			}
			t += n
			if buf[t-2] == '$' { //assuming the $PS1 == 'sh-4.3$ '
				out <- string(buf[:t])
				t = 0
				wg.Done()
			}
		}
	}()
	return in, out
}
