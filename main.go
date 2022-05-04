package main

import (
	"joyroku/todos/app"
	"net/http"
)

func main() {
	/*
		깃허브 커밋을 위한 주석
	*/
	m := app.MakeHandler("./test.db")
	defer m.Close()
	err := http.ListenAndServe(":3000", m)
	if err != nil {
		panic(err)
	}

}
