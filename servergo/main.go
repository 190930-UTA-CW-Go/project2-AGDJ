package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type apps struct {
	Name string
	Desc string
}

func main() {
	Post("david", "chang")

	//////////////////

	fmt.Println("start application getting worker info")
	response, err := http.Get("http://localhost:8080/getbutlerinfo")
	var infoHolder ButlerInfoStruct
	if err != nil {
		fmt.Printf("HTTp request failed with err %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal(data, &infoHolder)
		fmt.Println(infoHolder.Lscpu)
	}
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
