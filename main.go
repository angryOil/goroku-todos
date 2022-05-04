package main

import (
	"net/http"
	"os"
	"todos/app"
)

func main() {
	port := os.Getenv("PORT")
	/*
		깃허브 커밋을 위한 주석
	*/
	m := app.MakeHandler(os.Getenv("DATABASE_URL"))
	defer m.Close()
	err := http.ListenAndServe(":"+port, m)
	if err != nil {
		panic(err)
	}

}
