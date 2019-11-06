package main

import (
	_ "expvar"
	"fmt"
	"net/http"
)

//Login structure of login
type Login struct {
	username string
	pass     string
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("server")))
	http.HandleFunc("/login", login)
	http.HandleFunc("/numcontainer", numcontainer)
	http.ListenAndServe(":9000", nil)
}

func login(response http.ResponseWriter, request *http.Request) {
	var l Login

	l.username = request.FormValue("username")
	l.pass = request.FormValue("pass")

	fmt.Println(l)

}

func numcontainer(response http.ResponseWriter, request *http.Request) {

	numcontainer := request.FormValue("numcontainer")
	string := "string"
	fmt.Println(numcontainer)
	response = http.ResponseWriter("numcontainer.html")

}
