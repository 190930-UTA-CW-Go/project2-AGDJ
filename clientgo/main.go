package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

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
	aptprog.SearchProgHandler()
	holder := Apps{
		Apps: aptprog.GetSearchInfo(),
	}
	json.NewEncoder(w).Encode(holder)
}

func getButlerInfo(w http.ResponseWriter, r *http.Request) {
	sysinfo.CreateSystemInfoFile2()
	cpumem.CreateTopSnapshot()
	cpuusage.CreateCPUUsage()
	lscpu.CreateLSCPUFILE()
	butlerHolder := ButlerInfoStruct{
		Sysinfo: sysinfo.ReadSysInfo(),
		Lscpu:   lscpu.ReadLSCPUCommand(), CPUUsage: cpuusage.GetCPUUsage(),
		Cpumem: cpumem.GetTopSnapshot(),
	}
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
	appSearch := flag.Bool("search", false, "searches for program")
	appInstall := flag.String("install", "", "install programs using apt")
	appUninstall := flag.String("uninstall", "", "uninstall programs using apt")
	appUpgrade := flag.String("upgrade", "", "upgrade programs using apt")
	appKill := flag.String("kill", "", "kills a process based on id")
	appCreateReport := flag.Bool("report", false, "return a report on the system")
	flag.Parse()

	if *appCreateReport {
		systemInfoLoc := os.ExpandEnv("$HOME/lscpuvar.txt")
		cpuUsageLoc := os.ExpandEnv("$HOME/cpupercentage.txt")
		cpumemLoc := os.ExpandEnv("$HOME/cpumem.txt")
		cpumem.CreateTopSnapshot()
		cpuusage.CreateCPUUsage()
		lscpu.CreateLSCPUFILE()
		exec.Command("bash", "-c", "cat "+systemInfoLoc+" "+cpuUsageLoc+" "+cpumemLoc+">> "+os.ExpandEnv("$HOME")+"/superinfo.txt").Run()
	}

	if *appSearch {
		aptprog.SearchProgHandler()
	}

	if *appInstall != "" {
		// aptprog.InstallProgHandler(*appInstall)
		aptprog.InstallProgHandler(*appInstall, "")
	}

	if *appUninstall != "" {
		aptprog.UninstallProgHandler(*appUninstall, "")
	}

	if *appUpgrade != "" {
		aptprog.UpgradeProgHandler("")
	}

	if *appKill != "" {
		aptprog.KillProcessHandler(*appKill, "")
	}

	handleRequests()
}
