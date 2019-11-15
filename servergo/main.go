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
	Applications []aptprog.AptProgsStruct `json:"APPS"`
	Downloads    []string                 `json:"DOWNLOADED"`
}

// Super =
type Super struct {
	Machines []ButlerInfoStruct
	Count    []int
}

// var clients []string = []string{"52.176.60.129", "40.69.155.213"}
var clients []string = []string{"localhost"}
var superHolder Super = getSuperHolder()
var appsHolder Apps = getApps()

//////////// main function /////////////////
func main() {
	// Post("david", "chang")
	server.StartDB()

	//////////////////
	serveAndListen()
}

//////////////////// Front End Functions //////////////////
func serveAndListen() {
	http.Handle("/", http.FileServer(http.Dir("server")))
	http.HandleFunc("/welcome", welcome)
	http.HandleFunc("/register", register)
	//http.HandleFunc("/open", open)
	http.HandleFunc("/signin", enter)
	http.HandleFunc("/installapps", installApps)
	http.HandleFunc("/typedprogs", typedprogs)
	http.HandleFunc("/uninstall", uninstall)
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
		server.CreateAccount(username, password)
	}
	http.Redirect(w, r, "/welcome", http.StatusSeeOther)
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
		http.Redirect(w, r, "/welcome", http.StatusSeeOther)

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
	holder := appsHolder
	// fmt.Println(holder.Apps)
	//log.Println(temp.Execute(w, holder))
	temp.Execute(w, holder)
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
		server.AddInstalled(programs[i])
	}
	fmt.Println(progs)
	//installation here
	InstallOnClients(progs)
	installedData := server.QueryAllInstalled()
	fmt.Println("this is installed data:", installedData)
	tableData := Apps{Applications: progs, Downloads: installedData}
	fmt.Println("tabledata downloads:", tableData.Downloads)
	fmt.Println("CALLED DATA")
	temp.Execute(w, tableData)
}

// selected gets inputs passed in form and downloads programs
func uninstall(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("server/templates/uninstalled.html")
	if err != nil {
		log.Println("Error in uninstalled function template")
	}
	//populating selected program into []AptProgStruct
	query := r.FormValue("pname")
	query = strings.Replace(query, ",", "", -1)
	programs := strings.Fields(query)
	progs := make([]aptprog.AptProgsStruct, 0)
	for i := 0; i < len(programs); i++ {
		progs = append(progs, aptprog.AptProgsStruct{Name: programs[i]})
		server.DeleteInstalled(programs[i])
	}
	fmt.Println(progs)
	//installation here
	UninstallOnClients(progs)
	installedData := server.QueryAllInstalled()
	fmt.Println("this is uninstalled data:", installedData)
	tableData := Apps{Applications: progs, Downloads: installedData}
	fmt.Println("tabledata downloads:", tableData.Downloads)
	fmt.Println("CALLED DATA")
	temp.Execute(w, tableData)
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

//InstallOnClients installs selected applications on all the client machines
func UninstallOnClients(install []aptprog.AptProgsStruct) {
	for _, value := range clients {
		PostProgramsToUninstall(install, value)
	}
}

//PostProgramsToInstall will send program list of things to be installed
func PostProgramsToUninstall(install []aptprog.AptProgsStruct, ip string) {
	marshData, err := json.Marshal(install)
	if err != nil {
		log.Println("Marshaling program installation went wrong")
	}
	url := "http://" + ip + ":8080/uninstall"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(marshData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("client sending uninstall list fail")
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
	fmt.Println("start application getting applications info")
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
