package main

import (
	_ "expvar"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
  
  "github.com/190930-UTA-CW-Go/project2-AGDJ/commands"
	"github.com/190930-UTA-CW-Go/project2-AGDJ/opendb"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("client")))
	http.HandleFunc("/login", login)
	http.HandleFunc("/numcontainer", numcontainer)
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
  //username := request.FormValue("username")
	//pass := request.FormValue("pass")
  //fmt.Println("Username: " + username)
	//fmt.Println("Password: " + pass)
	// commands.SignIn("username, password", username, pass)
	//commands.CreateAccount(username, pass)
	//users := commands.QueryAllUsers()
	//fmt.Println(users)
	// id, name, pw := commands.QueryUser("ben")
	// fmt.Println(id, name, pw)
	// commands.CreateRunning(8080, "ben")
	// commands.CreateRunning(9000, "ben")
	// commands.CreateAccount("godfrey", "hello")
	// commands.CreateRunning(9090, "godfrey")
	// containers := commands.QueryAllRunning("")
	// fmt.Println(containers)
	// containers = commands.QueryAllRunning("ben")
	// fmt.Println(containers)
	// commands.DeleteRunning(8080)
	// containers = commands.QueryAllRunning("")
	// fmt.Println(containers)
	user := Loggedin{false}
	uname := request.FormValue("username")
	pass := request.FormValue("pw")
	temp, _ := template.ParseFiles("client/templates/login.html")
	if uname == "akhv" {
		if pass == "password" {
			user.Signedin = true
		} else {
			user.Signedin = false
		}
	}
	//here we pass in the user which is a Loggedin Struct which holds only one value of
	//boolean this will help the html template workout which html to show
	//you can see the usecase of the template in the login.html which handles that
	fmt.Println(temp.Execute(response, user))

}

//should be a nice display page welcoming the user into webserver asking how many
//alpine images they would like to run
func numcontainer(response http.ResponseWriter, request *http.Request) {
	numcon := Numcont{false}
	numcust := request.FormValue("numcontainer")
	temp1, _ := template.ParseFiles("client/templates/numcontainer.html")
	numcust1, _ := strconv.Atoi(numcust)
	if numcust1 > 0 {
		numcon.Numcon = true
	} else {
		numcon.Numcon = false
	}

	fmt.Println(temp1.Execute(response, numcon))
}
