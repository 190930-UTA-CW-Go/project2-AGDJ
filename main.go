package main

import (
	_ "expvar"
	"fmt"
	"net/http"

	"github.com/gmac220/project2-AGDJ/commands"
	"github.com/gmac220/project2-AGDJ/opendb"
)

func main() {
	opendb.StartDB()
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/login", login)
	http.ListenAndServe(":9000", nil)
}

func login(response http.ResponseWriter, request *http.Request) {
	username := request.FormValue("username")
	pass := request.FormValue("pass")
	fmt.Println("Username: " + username)
	fmt.Println("Password: " + pass)
	// commands.SignIn("username, password", username, pass)
	commands.CreateAccount(username, pass)
	users := commands.QueryAllUsers()
	fmt.Println(users)
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
	//output := []byte("Username: " + username + "Pass: " + pass)

	//response.Write(output)
}
