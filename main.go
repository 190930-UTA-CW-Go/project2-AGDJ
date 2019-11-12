package main

import (
	_ "expvar"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/190930-UTA-CW-Go/project2-AGDJ/commands"
	"github.com/190930-UTA-CW-Go/project2-AGDJ/opendb"
	"github.com/190930-UTA-CW-Go/project2-AGDJ/sshsetup"
)

func main() {
	opendb.StartDB()
	http.Handle("/", http.FileServer(http.Dir("client")))
	//http.HandleFunc("/", welcome)
	http.HandleFunc("/index1", login)
	http.HandleFunc("/numcontainer2", numcontainer)
	http.ListenAndServe(":9000", nil)
}

//Loggedin is a structure whic will hold the value to let the user into the server to edit information
type Loggedin struct {
	Signedin bool
}

//Numcont is a structure for number of containers
type Numcont struct {
	Numcon bool
}

//login function should verify entered usernamen and password against the database data
//sets the appropriate variables against the
func login(response http.ResponseWriter, request *http.Request) {
	user := Loggedin{false}
	uname := request.FormValue("username")
	pass := request.FormValue("pw")
	commands.CreateAccount("godfrey", "hello")
	temp, _ := template.ParseFiles("client/login.html")
	// fmt.Println("form value", uname, pass)
	_, unamedb, passdb := commands.SignIn(uname)
	if uname == unamedb {
		if pass == passdb {
			user.Signedin = true
		} else {
			user.Signedin = false
		}
	}
	//fmt.Println(temp.Execute(response, user))
	temp.Execute(response, user)
}

//should be a nice display page welcoming the user into webserver asking how many
//alpine images they would like to run
func numcontainer(response http.ResponseWriter, request *http.Request) {
	numcon := Numcont{false}
	numcust := request.FormValue("numcontainer2")
	temp1, _ := template.ParseFiles("client/templates/numcontainer2.html")
	numcust1, _ := strconv.Atoi(numcust)
	if numcust1 > 0 {
		numcon.Numcon = true
		login := "_"
		password := "_"
		ip := "_"
		port := "22"

		sshsetup.Connect(login, password, ip, port)
	} else {
		numcon.Numcon = false
	}

	fmt.Println(temp1.Execute(response, numcon))
}
