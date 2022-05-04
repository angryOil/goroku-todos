package main

import (
	"joyroku/todos/app"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	/*
		깃허브 커밋을 위한 주석
	*/
	m := app.MakeHandler("./test.db")
	defer m.Close()
	err := http.ListenAndServe(":"+port, m)
	if err != nil {
		panic(err)
	}

}
