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

func getButlerInfo(w http.ResponseWriter, r *http.Request) {
	cpumem.CreateTopSnapshot()
	cpuusage.CreateCPUUsage()
	lscpu.CreateLSCPUFILE()
	butlerHolder := ButlerInfoStruct{
		Lscpu: lscpu.ReadLSCPUCommand(), CPUUsage: cpuusage.GetCPUUsage(),
		Cpumem: cpumem.GetTopSnapshot(),
	}
	json.NewEncoder(w).Encode(butlerHolder)
}

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

func handleRequests() {
	route := mux.NewRouter().StrictSlash(true)
	route.HandleFunc("/getbutlerinfo", getButlerInfo)
	route.HandleFunc("/userinfo", userInfo)
	log.Fatal(http.ListenAndServe(":8080", route))
}

func main() {
	handleRequests()
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
}
