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
	"github.com/gorilla/mux"
)

//ButlerInfoStruct will be used to pass vital butler client information
type ButlerInfoStruct struct {
	Lscpu    lscpu.LSCPU              `json:"LSCPU"`
	CPUUsage cpuusage.CPUUsage        `json:"CPUUSAGE"`
	Cpumem   cpumem.CPUTOP            `json:"CPUMEM"`
	Apps     []aptprog.AptProgsStruct `json:"APPS"`
}

// UserInfo =
type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func getButlerInfo(w http.ResponseWriter, r *http.Request) {
	cpumem.CreateTopSnapshot()
	cpuusage.CreateCPUUsage()
	lscpu.CreateLSCPUFILE()
	aptprog.SearchProgHandler()
	butlerHolder := ButlerInfoStruct{
		Lscpu: lscpu.ReadLSCPUCommand(), CPUUsage: cpuusage.GetCPUUsage(),
		Cpumem: cpumem.GetTopSnapshot(), Apps: aptprog.GetSearchInfo(),
	}
	json.NewEncoder(w).Encode(butlerHolder)
}

//UserInfo function handles incoming post request for list of programs to be installed.
func userInfo(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//case "GET":
	case "POST":
		result, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		user := UserInfo{}
		json.Unmarshal(result, &user)
		//fmt.Println(user)
		fmt.Println(user.Username)
		fmt.Println(user.Password)
		w.Write([]byte("Received a POST request\n"))
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}
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

//BatchInstall will Install a bathc of passed in programs.
func BatchInstall(hold []aptprog.AptProgsStruct) {
	for _, k := range hold {
		aptprog.InstallProgHandler(k.Name)
	}
	log.Println("Succsessful installation of programs")
}

//////////////////////// HTTP SERVER HERE //////////////////
func handleRequests() {
	route := mux.NewRouter().StrictSlash(true)
	route.HandleFunc("/getbutlerinfo", getButlerInfo)
	route.HandleFunc("/userinfo", userInfo)
	route.HandleFunc("/install", ProgramsToInstall)
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
		aptprog.InstallProgHandler(*appInstall)
	}

	if *appUninstall != "" {
		aptprog.UninstallProgHandler(*appUninstall)
	}

	if *appUpgrade != "" {
		aptprog.UpgradeProgHandler(*appUpgrade)
	}

	if *appKill != "" {
		aptprog.KillProcessHandler(*appKill)
	}

	handleRequests()
}
