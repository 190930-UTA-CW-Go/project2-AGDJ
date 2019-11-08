package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"text/template"
)

//CMDOUT is a struct
type CMDOUT struct {
	CmdOutput []string
}

func main() {
	fmt.Println("Server started")
	http.Handle("/", http.FileServer(http.Dir("client")))
	http.HandleFunc("/indexpage", indexpage)
	//runCommand("exit")
	//fmt.Println(readFile())
	http.ListenAndServe(":8080", nil)
}
func indexpage(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("client/template/cmd.html")
	if err != nil {
		log.Fatal(err)
	}
	hold := r.FormValue("command")
	runCommand(hold)
	CDMOUT2 := readFile()
	fmt.Println(temp.Execute(w, CDMOUT2))
}
func runCommand(str string) {
	//opening of the file
	cmdFile := os.ExpandEnv("$HOME/bashoutput.txt")
	file, _ := os.OpenFile(cmdFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_RDONLY, 0666)
	defer file.Close()
	if str == "exit" {
		file.Write([]byte("Cannot exit this terminal kill instance from pod.com\n"))
	} else if str == "" {
		file.Write([]byte("Welcome Pleb\n"))
	} else {
		cdm, _ := exec.Command("sh", "-c", str).CombinedOutput()
		file.Write(cdm)
	}
}

func readFile() CMDOUT {
	cmdFile := os.ExpandEnv("$HOME/bashoutput.txt")
	file, err := os.Open(cmdFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	hold := CMDOUT{}
	scanner := bufio.NewReader(file)
	var txthold string
	for {
		txthold, err = scanner.ReadString('\n')
		if err != nil {
			break
		}
		hold.CmdOutput = append(hold.CmdOutput, txthold[:len(txthold)-1])
	}
	if err != io.EOF {
		fmt.Printf("failed: %v\n", err)
	}
	return hold
}
