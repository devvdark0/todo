package main

import "github.com/devvdark0/todo/internal/app"

func main() {
	if err := app.InitApp(); err != nil {
		panic(err)
	}
}
