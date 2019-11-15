package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/190930-UTA-CW-Go/project2-AGDJ/clientgo/aptprog"
	"github.com/190930-UTA-CW-Go/project2-AGDJ/clientgo/cpumem"
	"github.com/190930-UTA-CW-Go/project2-AGDJ/clientgo/cpuusage"
	"github.com/190930-UTA-CW-Go/project2-AGDJ/clientgo/lscpu"
	"github.com/190930-UTA-CW-Go/project2-AGDJ/clientgo/sysinfo"
	"github.com/gorilla/mux"
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

func getApps(w http.ResponseWriter, r *http.Request) {
	defer def()
	aptprog.SearchProgHandler()
	//see maybe file doesnt get created in time
	//check if asynchronous
	holder := Apps{
		Apps: aptprog.GetSearchInfo(),
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(holder)
	if err != nil {
		panic(err)
	}
}

func getButlerInfo(w http.ResponseWriter, r *http.Request) {
	defer def()
	sysinfo.CreateSystemInfoFile2()
	cpumem.CreateTopSnapshot()
	cpuusage.CreateCPUUsage()
	lscpu.CreateLSCPUFILE()
	butlerHolder := ButlerInfoStruct{
		Sysinfo: sysinfo.ReadSysInfo(),
		Lscpu:   lscpu.ReadLSCPUCommand(), CPUUsage: cpuusage.GetCPUUsage(),
		Cpumem: cpumem.GetTopSnapshot(),
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(butlerHolder)
}

//ProgramsToInstall handles post requests of desired applications to be installed.
func ProgramsToInstall(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		UnmarshProgList, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		var progList []aptprog.AptProgsStruct
		json.Unmarshal(UnmarshProgList, &progList)
		//call on batch install
		BatchInstall(progList)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}
}

//ProgramsToUninstall will call batch uninstall and handle the api
func ProgramsToUninstall(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		UnmarshProgList, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		var progList []aptprog.AptProgsStruct
		json.Unmarshal(UnmarshProgList, &progList)
		//call on batch install
		BatchUninstall(progList)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}
}

//BatchUninstall will uninstall a batchof passed in programs
func BatchUninstall(hold []aptprog.AptProgsStruct) {
	for _, k := range hold {
		// aptprog.InstallProgHandler(k.Name)
		fmt.Println(k.Name + "Uninstalling ...")
		aptprog.UninstallProgHandler(k.Name, "")
	}
	fmt.Println("Done")
	log.Println("Successful removal of program/s")
}

//BatchInstall will Install a bathc of passed in programs.
func BatchInstall(hold []aptprog.AptProgsStruct) {
	aptprog.UpgradeProgHandler("")
	for _, k := range hold {
		// aptprog.InstallProgHandler(k.Name)
		fmt.Println(k.Name + "installing ...")
		aptprog.InstallProgHandler(k.Name, "")
	}
	fmt.Println("Done")
	log.Println("Successful installation of program/s")
}

func def() {
	fmt.Println("Program Panicked out!")
	if r := recover(); r != nil {
		fmt.Println("WHy is it breaking though")
	}
}

//////////////////////// HTTP SERVER HERE //////////////////
func handleRequests() {
	route := mux.NewRouter().StrictSlash(true)
	route.HandleFunc("/getbutlerinfo", getButlerInfo)
	route.HandleFunc("/install", ProgramsToInstall)
	route.HandleFunc("/uninstall", ProgramsToUninstall)
	route.HandleFunc("/apps", getApps)
	log.Fatal(http.ListenAndServe(":8080", route))
}

////////////////////// MAIN FUNCTION HERE ///////////////////
func main() {
	handleRequests()
}
