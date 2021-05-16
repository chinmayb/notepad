package main

import (
	"log"

	"github.com/chinmayb/notepad/cmd/root"
)

func main() {

	command := root.NewRootCommand()
	if err := command.Execute(); err != nil {
		log.Fatal()
	}
}
