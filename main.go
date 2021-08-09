package main

import (
	"github.com/kylegk/notes/app"
	"github.com/kylegk/notes/router"
)

func main ()  {
	app.Init()
	router.AddRouting()
}