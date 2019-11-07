package opendb

import (
	"database/sql"
	"fmt"
	"os/exec"

	_ "github.com/lib/pq"
)

// OpenDB opens postgres database
func OpenDB() *sql.DB {
	datasource := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "postgres", "postgres")
	db, err := sql.Open("postgres", datasource)
	if err != nil {
		panic(err)
	}
	return db
}

// StartDB starts the docker database
func StartDB() {
	fmt.Println("starting db build")
	exec.Command("docker", "stop", "runningproject2").Run()
	exec.Command("docker", "rmi", "project2").Run()
	// exec.Command("cd", "~/go/src/github.com/gmac220/project2-AGDJ/db").Run()
	// exec.Command("cd", "~/go/src/github.com/190930-UTA-CW-Go/project2-AGDJ/db").Run()
	exec.Command("docker", "build", "-t", "project2", "db/.").Run()
	exec.Command("docker", "run", "-p=5432:5432", "-d", "--rm", "--name=runningproject2", "project2").Run()
	fmt.Println("Running container")
	// db := OpenDB()
	// defer db.Close()
	// db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", "godfrey", "hello")
	// fmt.Println("finished inserting values")
}
