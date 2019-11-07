package main

import (
	_ "expvar"
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("client")))
	http.HandleFunc("/login", login)
	//http.HandleFunc("/numcontainer", numcontainer)
	http.ListenAndServe(":9000", nil)
}

//Loggedin is a structure whic will hold the value to let the user into the server to edit information
type Loggedin struct {
	Signedin bool
}

//login function should verify entered usernamen and password against the database data
//sets the appropriate variables against the
func login(response http.ResponseWriter, request *http.Request) {
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

}
