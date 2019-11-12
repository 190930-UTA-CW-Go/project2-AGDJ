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
	Lscpu    lscpu.LSCPU
	CPUUsage cpuusage.CPUUsage `json:"CPUUSAGE"`
	Cpumem   cpumem.CPUTOP     `json:"CPUMEM"`
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

func main() {
	Post("david", "chang")

	//////////////////
	serveAndListen()
}
func getWorkerInfo() ButlerInfoStruct {
	fmt.Println("start application getting worker info")
	response, err := http.Get("http://localhost:8080/getbutlerinfo")
	var infoHolder ButlerInfoStruct
	if err != nil {
		fmt.Printf("HTTp request failed with err %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal(data, &infoHolder)
	}
	return infoHolder
}

func serveAndListen() {
	http.Handle("/", http.FileServer(http.Dir("server")))
	http.HandleFunc("/open", open)
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

// Post =
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
