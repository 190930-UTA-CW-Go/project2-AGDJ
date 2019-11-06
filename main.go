package main

import (
	"database/sql"
	_ "expvar"
	"fmt"
	"net/http"

	"github.com/gmac220/project2-AGDJ/opendb"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/login", login)
	http.ListenAndServe(":9000", nil)
}

func login(response http.ResponseWriter, request *http.Request) {
	username := request.FormValue("username")
	pass := request.FormValue("pass")

	fmt.Println("Username: " + username)
	fmt.Println("Password: " + pass)
	SignIn(username, pass)
	//output := []byte("Username: " + username + "Pass: " + pass)

	//response.Write(output)
}

// SignIn verifies if customer or employee credentials match database
func SignIn(username string, password string) (string, string) {
	var usernamedb, passdb string
	var row *sql.Row

	db := opendb.OpenDB()
	defer db.Close()
	row = db.QueryRow("SELECT username, password FROM users WHERE username = $1", username)
	row.Scan(&usernamedb, &passdb)
	fmt.Println("Logged in with", usernamedb, passdb)
	return usernamedb, passdb
}
