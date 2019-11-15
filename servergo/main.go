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
	"github.com/190930-UTA-CW-Go/project2-AGDJ/servergo/server"
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

// Super holds slice of client data and slice used for count
type Super struct {
	Machines []ButlerInfoStruct
	Count    []int
}

var clients []string = []string{"localhost", "localhost"}
var superHolder Super = getSuperHolder()

//////////// main function /////////////////
func main() {

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
	http.HandleFunc("/signin", enter)
	//http.HandleFunc("/open", open)
	http.HandleFunc("/installapps", installApps)
	http.HandleFunc("/typedprogs", typedprogs)
	http.ListenAndServe(":8081", nil)
}

func enter(w http.ResponseWriter, r *http.Request) {
	// temp, err := template.ParseFiles("index.html")
	// if err != nil {
	// 	log.Println("uuuupppsss")
	// }
	var username = r.FormValue("username")
	var pass = r.FormValue("pw")

	fmt.Println(username)
	fmt.Println(pass)
	trigger := server.SignIn(username, pass)
	if trigger == true {
		http.Redirect(w, r, "/open", http.StatusSeeOther)

	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	// log.Println(temp.Execute(w, server.SignIn))
}

func open(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("server/templates/workerinfo.html")
	if err != nil {
		log.Println("this code sucks")
	}
	/////////////////////
	fmt.Println(superHolder.Machines[0].Lscpu)
	log.Println(temp.Execute(w, superHolder))
}

func installApps(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("server/templates/installapps.html")
	if err != nil {
		log.Println("error in install Apps")
	}
	holder := superHolder.Machines[0]
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
	//populating selected program into []AptProgStruct
	query := r.FormValue("pname")
	query = strings.Replace(query, ",", "", -1)
	programs := strings.Fields(query)
	progs := make([]aptprog.AptProgsStruct, 0)
	for i := 0; i < len(programs); i++ {
		progs = append(progs, aptprog.AptProgsStruct{Name: programs[i]})
	}
	fmt.Println(progs)
	//installation here
	InstallOnClients(progs)
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
	url := "http://" + ip + ":8080/install"
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

//Gets information from client server
func getWorkerInfo(ip string) ButlerInfoStruct {
	fmt.Println("start application getting worker info")
	var infoHolder ButlerInfoStruct
	response, err := http.Get("http://" + ip + ":8080/getbutlerinfo")
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
