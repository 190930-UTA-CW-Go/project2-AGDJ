package main

import (
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
	Lscpu    lscpu.LSCPU       `json:"LSCPU"`
	CPUUsage cpuusage.CPUUsage `json:"CPUUSAGE"`
	Cpumem   cpumem.CPUTOP     `json:"CPUMEM"`
}

func main() {
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
