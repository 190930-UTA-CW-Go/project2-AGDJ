package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	appSearch := flag.String("search", "", "searches for program")
	appInstall := flag.String("install", "", "install programs using apt")
	appUninstall := flag.String("uninstall", "", "uninstall programs using apt")
	appUpgrade := flag.String("upgrade", "", "upgrade programs using apt")
	appKill := flag.String("kill", "", "kills a process based on id")
	flag.Parse()

	if *appSearch != "" {
		SearchProgHandler(*appSearch)
	}

	if *appInstall != "" {
		InstallProgHandler(*appInstall)
	}

	if *appUninstall != "" {
		UninstallProgHandler(*appUninstall)
	}

	if *appUpgrade != "" {
		UpgradeProgHandler(*appUpgrade)
	}

	if *appKill != "" {
		KillProcessHandler(*appKill)
	}
}

// SearchProgHandler searches apt for specified program user searches for and puts it inside a text file
func SearchProgHandler(appname string) {
	searchFile := os.ExpandEnv("$HOME/searchapps.txt")
	file, err := os.OpenFile(searchFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_RDONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	usrSearch := exec.Command("apt", "search", appname)
	searchOutput, stderr := usrSearch.Output()
	if stderr != nil {
		fmt.Println(stderr)
	}
	file.Write(searchOutput)
}

// InstallProgHandler installs the program that is passed in
func InstallProgHandler(appname string) {
	exec.Command("sudo", "apt", "install", "-y", appname).Run()
}

// UpgradeProgHandler upgrades the program to the latest version in apt
func UpgradeProgHandler(appname string) {
	exec.Command("sudo", "apt", "upgrade", "-y", appname).Run()
}

// UninstallProgHandler removes the program that is passed in
func UninstallProgHandler(appname string) {
	exec.Command("sudo", "apt", "purge", "-y", appname).Run()
}

// KillProcessHandler kills process
func KillProcessHandler(procid string) {
	exec.Command("sudo", "kill", procid).Run()
}
