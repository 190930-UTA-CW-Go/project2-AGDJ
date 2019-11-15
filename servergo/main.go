package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/190930-UTA-CW-Go/project2-AGDJ/servergo/aptprog"
	"github.com/190930-UTA-CW-Go/project2-AGDJ/servergo/cpumem"
	"github.com/190930-UTA-CW-Go/project2-AGDJ/servergo/cpuusage"
	"github.com/190930-UTA-CW-Go/project2-AGDJ/servergo/lscpu"
	"github.com/190930-UTA-CW-Go/project2-AGDJ/servergo/sysinfo"
)

//ButlerInfoStruct will be used to pass vital butler client information
type ButlerInfoStruct struct {
	Sysinfo  sysinfo.SysInfo   `json:"SYSINFO"`
	Lscpu    lscpu.LSCPU       `json:"LSCPU"`
	CPUUsage cpuusage.CPUUsage `json:"CPUUSAGE"`
	Cpumem   cpumem.CPUTOP     `json:"CPUMEM"`
}

//Apps will pass the apps.
type Apps struct {
	Apps []aptprog.AptProgsStruct `json:"APPS"`
}

// Super =
type Super struct {
	Machines []ButlerInfoStruct
	Count    []int
}

var clients []string = []string{"52.176.60.129", "40.69.155.213"}
var superHolder Super = getSuperHolder()

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
	http.HandleFunc("/welcome", welcome)
	http.HandleFunc("/register", register)
	http.HandleFunc("/open", open)
	http.HandleFunc("/installapps", installApps)
	http.HandleFunc("/typedprogs", typedprogs)
	http.ListenAndServe(":8081", nil)
}

func welcome(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("server/templates/welcome.html")
	if err != nil {
		log.Println(err)
	}
	temp.Execute(w, nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("METHOD:", r.Method)
	switch r.Method {
	case "POST":
		r.ParseForm()
		var username string = fmt.Sprint(r.Form["username"][0])
		var password string = fmt.Sprint(r.Form["password"][0])
		fmt.Println(username)
		fmt.Println(password)
	}
	http.Redirect(w, r, "/welcome", http.StatusSeeOther)
}

func open(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("server/templates/workerinfo.html")
	if err != nil {
		log.Println("this code sucks")
	}
	holder := getWorkerInfo()

	/////////////////////

	var count int = 5
	var manyButlers []ButlerInfoStruct
	for i := 0; i < count; i++ {
		manyButlers = append(manyButlers, holder)
	}

	var numMachines []int
	for i := 1; i <= len(manyButlers); i++ {
		numMachines = append(numMachines, i)
	}
	superMan := Super{manyButlers, numMachines}

	fmt.Println(holder.Lscpu.Architecture)
	log.Println(temp.Execute(w, superMan))
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

// selected gets inputs passed in form and downloads programs
func typedprogs(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("server/templates/typedprogs.html")
	if err != nil {
		log.Println("Error in typedprogs function template")
	}
	query := r.FormValue("pname")
	query = strings.Replace(query, ",", "", -1)
	programs := strings.Fields(query)
	progs := make([]aptprog.AptProgsStruct, 0)
	for i := 0; i < len(programs); i++ {
		progs = append(progs, aptprog.AptProgsStruct{Name: programs[i]})
	}
	fmt.Println(progs)
	infoByte, _ := json.Marshal(progs)
	url := "http://localhost:8080/searchinstall"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(infoByte))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	temp.Execute(w, progs)
}

///////////////////// API FUNCTIONS ///////////////////////////////////

//InstallOnClients installs selected applications on all the client machines
func InstallOnClients(install []aptprog.AptProgsStruct) {
	for _, value := range clients {
		PostProgramsToInstall(install, value)
	}
}

//PostProgramsToInstall will send program list of things to be installed
func PostProgramsToInstall(install []aptprog.AptProgsStruct, ip string) {

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
    
func getApps() Apps {
	fmt.Println("start application getting worker info")
	var infoHolder Apps
	response, err := http.Get("http://" + clients[0] + ":8080/apps")
	if err != nil {
		fmt.Printf("HTTP request failed with err %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal(data, &infoHolder)
	}
	return infoHolder
}

func getSuperHolder() Super {
	var holder Super
	for key, val := range clients {
		holder.Machines = append(holder.Machines, getWorkerInfo(val))
		holder.Count = append(holder.Count, key+1)
	}
	return holder
}
