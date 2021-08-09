package app

import (
	"github.com/kylegk/notes/db"
)

type Configuration struct {
	DB   db.DB
}

var Context *Configuration

func Init() {
	c := &Configuration{}
	dbConn, err := db.InitDB(db.Schema)
	if err != nil {
		panic(err)
	}
	c.DB = dbConn
	Context = c
}