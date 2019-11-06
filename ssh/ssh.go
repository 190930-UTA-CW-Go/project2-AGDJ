package ssh

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

// CmdGetInfo =
func CmdGetInfo(login string, ip string) string {
	result, _ := exec.Command("ssh", login+"@"+ip, "whoami").Output()
	currentUser := string(result)
	currentUser = currentUser[:len(currentUser)-1]

	hostName, _ := exec.Command("ssh", login+"@"+ip, "hostname").Output()

	return currentUser + "@" + string(hostName)
}

// Install =
func Install(login string, password string, ip string) {
	// Ex) "wyrd"
	// Ex) "curl"
	// Ex) "docker.io"
	var app string
	fmt.Print("Type program to install: ")
	fmt.Scan(&app)
	fmt.Println("")

	results, _ := exec.Command("ssh", login+"@"+ip, "echo "+password+" | sudo -S apt install "+app+" -y").Output()
	fmt.Println(string(results))
}

func connect(username string, password string, ip string, port string) {
	// Connect to ssh cient
	config := &ssh.ClientConfig{
		//To resolve "Failed to dial: ssh: must specify HostKeyCallback" error
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		User:            username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}
	client, err := ssh.Dial("tcp", ip+":"+port, config)
	if err != nil {
		fmt.Println("Invalid Login or Password")
		//http.Redirect(response, request, "/login", http.StatusSeeOther)
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

		///////////////////////////////////////////////////
		fmt.Println(GetInfo(in, out))
		///////////////////////////////////////////////////

		// automatically "exit"
		in <- "exit"
		session.Wait()
	}
}

// GetInfo =
func GetInfo(in chan<- string, out <-chan string) string {
	// Get current user and slice off '\n'
	in <- "whoami"
	slice := strings.Fields(<-out)

	// Get host name
	in <- "hostname"
	hostName := <-out

	// Concatenate to resemble Linux
	return slice[0] + "@" + hostName
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
