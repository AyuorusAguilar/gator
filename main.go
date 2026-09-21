package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/AyuorusAguilar/gator/internal/commands"
	"github.com/AyuorusAguilar/gator/internal/config"
	"github.com/AyuorusAguilar/gator/internal/database"
	"github.com/AyuorusAguilar/gator/internal/state"
	_ "github.com/lib/pq"
)

func main() {
	conf, err := config.Read() /* Conf */
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	// fmt.Printf("Initial conf: %v\n", *conf)

	s := state.State{Cfg: conf} /* State */

	c := commands.NewCommander() /* Command Struct */
	c.Register("login", commands.HandlerLogin)
	c.Register("register", commands.HandlerRegister)
	c.Register("reset", commands.HandlerReset)
	c.Register("users", commands.HandlerUsers)

	db, err := sql.Open("postgres", conf.DbURL)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	s.Db = database.New(db)
	argu := os.Args[1:] /* Getting arguments from bash */
	if len(argu) < 1 {
		fmt.Printf("No command recived\n")
		os.Exit(1)
	}
	// fmt.Printf("Argu value:\n\t%v\n", argu)

	newCommand := commands.Command{Name: argu[0], Arguments: argu[1:]}
	err = c.Run(&s, newCommand)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	/* var lastConf *config.Config
	lastConf, err = config.Read()
	if err != nil {
		os.Exit(1)
	}
	fmt.Printf("Last conf:\n\t%v\n", lastConf) */
	os.Exit(0)
}
