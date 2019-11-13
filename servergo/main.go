package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/190930-UTA-CW-Go/project2-AGDJ/servergo/cpumem"
	"github.com/190930-UTA-CW-Go/project2-AGDJ/servergo/cpuusage"
	"github.com/190930-UTA-CW-Go/project2-AGDJ/servergo/lscpu"
)

//ButlerInfoStruct will be used to pass vital butler client information
type ButlerInfoStruct struct {
	Lscpu    lscpu.LSCPU       `json:"LSCPU"`
	CPUUsage cpuusage.CPUUsage `json:"CPUUSAGE"`
	Cpumem   cpumem.CPUTOP     `json:"CPUMEM"`
	Apps     []AptProgsStruct  `json:"APPS"`
}

// UserInfo =
type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//AptProgsStruct lists all the
type AptProgsStruct struct {
	Name string `json:"APPNAME"`
	Desc string `json:"DESC"`
}

//////////// main function /////////////////
func main() {
	Post("david", "chang")
	//testing if PostProgram will install a single application
	// var wyrd []AptProgsStruct
	// wyrd = append(wyrd, AptProgsStruct{Name: "wyrd", Desc: "Description dont matter"})
	// PostProgramsToInstall(wyrd)

	//////////////////
	serveAndListen()
}

//////////////////// Front End Functions //////////////////
func serveAndListen() {
	http.Handle("/", http.FileServer(http.Dir("server")))
	http.HandleFunc("/open", open)
	http.HandleFunc("/installapps", installApps)
	http.ListenAndServe(":8081", nil)
}

func open(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("server/templates/workerinfo.html")
	if err != nil {
		log.Println("this code sucks")
	}
	holder := getWorkerInfo()
	fmt.Println(holder.Lscpu.Architecture)
	log.Println(temp.Execute(w, holder))
}

func installApps(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("server/templates/installapps.html")
	if err != nil {
		log.Println("error in install Apps")
	}
	holder := getWorkerInfo()
	// fmt.Println(holder.Apps)
	log.Println(temp.Execute(w, holder))
	// temp.Execute(w, holder)
}

///////////////////// API FUNCTIONS ///////////////////////////////////

// Post function was the first triel of server sending post
func Post(username string, password string) {
	infoData := &UserInfo{
		Username: username,
		Password: password}
	infoByte, _ := json.Marshal(infoData)
	url := "http://localhost:8080/userinfo"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(infoByte))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

//PostProgramsToInstall will send program list of things to be installed
func PostProgramsToInstall(install []AptProgsStruct) {
	marshData, err := json.Marshal(install)
	if err != nil {
		log.Println("Marshaling program installation went wrong")
	}
	url := "http://localhost:8080/install"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(marshData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("client sending install list fail")
	}
	defer resp.Body.Close()
	fmt.Println("Sent List")
}

//Gets information from remote server
func getWorkerInfo() ButlerInfoStruct {
	fmt.Println("start application getting worker info")
	response, err := http.Get("http://localhost:8080/getbutlerinfo")
	var infoHolder ButlerInfoStruct
	if err != nil {
		fmt.Printf("HTTP request failed with err %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal(data, &infoHolder)
	}
	return infoHolder
}
