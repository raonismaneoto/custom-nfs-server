package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		panic("no command has been provided")
	}

	switch command := args[0]; command {
	case "mount":
		log.Println("exec mount command")

	case "save":
		log.Println("exec save command")
	case "read":
		log.Println("exec read command")
	case "chpem":
		log.Println("exec chpem command")
	default:
		fmt.Printf("Usage: \nexecutable <command> <args> [options]")
	}
}
