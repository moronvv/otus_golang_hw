package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("too few arguments")
	}
	dirPath := os.Args[1]
	cmd := os.Args[2:]

	envs, err := ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(RunCmd(cmd, envs))
}
