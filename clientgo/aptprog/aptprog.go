package aptprog

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

//AptProgsStruct lists all the
type AptProgsStruct struct {
	Name string `json:"APPNAME"`
	Desc string `json:"DESC"`
}

// SearchProgHandler searches apt for specified program user searches for and puts it inside a text file
func SearchProgHandler() {
	searchFile := os.ExpandEnv("$HOME/searchapps.txt")
	file, err := os.OpenFile(searchFile, os.O_CREATE|os.O_WRONLY|os.O_RDONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	usrSearch := exec.Command("apt", "search", ".")
	searchOutput, stderr := usrSearch.Output()
	if stderr != nil {
		fmt.Println(stderr)
	}
	file.Write(searchOutput)
}

// GetSearchInfo parses the data inside searchapps to get all the program w/ descriptions
func GetSearchInfo() []AptProgsStruct {
	var progs []AptProgsStruct = make([]AptProgsStruct, 1)
	searchFile := os.ExpandEnv("$HOME/searchapps.txt")
	file, err := os.Open(searchFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewReader(file)
	prog := AptProgsStruct{}

	var txthold string
	var count int = 0
	var index int = 0
	var lineSwitcher bool = true
	for err != io.EOF {
		txthold, err = scanner.ReadString('\n')
		if count > 2 {
			if txthold != "\n" {
				if lineSwitcher {
					progs[index].Desc = txthold[2:]
					lineSwitcher = false

				} else {
					progs[index].Name = txthold
					lineSwitcher = true
				}
			} else {
				index++
				prog = AptProgsStruct{Name: "", Desc: ""}
				progs = append(progs, prog)
			}
		}
		count++
	}

	// fmt.Println(progs)
	return progs
}

// InstallProgHandler installs the program that is passed in
func InstallProgHandler(appname string) {
	exec.Command("sudo", "apt", "install", "-y", appname).Run()
}

// UpgradeProgHandler upgrades the program to the latest version in apt
func UpgradeProgHandler(appname string) {
	exec.Command("sudo", "apt", "upgrade", "-y", appname).Run()
}

// UninstallProgHandler removes the program that is passed in
func UninstallProgHandler(appname string) {
	exec.Command("sudo", "apt", "purge", "-y", appname).Run()
}

// KillProcessHandler kills process
func KillProcessHandler(procid string) {
	exec.Command("sudo", "kill", procid).Run()
}
