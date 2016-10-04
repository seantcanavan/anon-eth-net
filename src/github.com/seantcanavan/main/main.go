package main

import (
	"os"
	"strings"
	"time"
	"config"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("One and only argument accepted: the file path to config.json to initialize the program")
		os.Exit(1)
	}

	if os.Args[1] == "h" || os.Args[1] == "help" || os.Args[1] == "?" {
		fmt.Println("One and only argument accepted: the file path to config.json to initialize the program")
	}

	if _, err := os.Stat(os.Args[1]); err == nil {
  		cfg := config.GetConfigFromFile(os.Args[1])
	} else {
		os.Exit(1)
	}

	go waitForUpdates()

	if(cfg.MineEther) {
		go mine()		
	}
}

func waitForUpdates() {
	time.Sleep()
}

func mine() {

}