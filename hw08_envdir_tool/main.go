package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("go-envdir envDir command <arg1>...")
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatalf("error reading envDir: %s", os.Args[1])
	}

	exitCode := RunCmd(os.Args[2:], env)
	os.Exit(exitCode)
}
