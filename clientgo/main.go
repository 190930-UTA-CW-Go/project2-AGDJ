package main

import (
	"encoding/json"
	"flag"
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
func handleRequests() {
	route := mux.NewRouter().StrictSlash(true)
	route.HandleFunc("/getbutlerinfo", getButlerInfo)
	log.Fatal(http.ListenAndServe(":8080", route))
}

func main() {
	handleRequests()
	appSearch := flag.String("search", "", "searches for program")
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

	if *appSearch != "" {
		aptprog.SearchProgHandler(*appSearch)
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
