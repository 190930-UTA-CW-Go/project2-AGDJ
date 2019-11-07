package main

import (
	"github.com/gittingdavid/project-2/potato"
)

func main() {
	login := "_"
	password := "_"
	ip := "_"
	port := "22"

	potato.Connect(login, password, ip, port)

}
