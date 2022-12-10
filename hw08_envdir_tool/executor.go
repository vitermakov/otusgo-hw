package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// при ошибках, возникающих до запуска подпроцесса вернем код = -1, так как 0, 1 заняты
	returnCode = -1

	if len(cmd) == 0 {
		log.Println("cmd not set")
		return
	}

	for key, value := range env {
		// если переменная уже существует сбросим ее.
		_, exists := os.LookupEnv(key)
		if exists {
			if err := os.Unsetenv(key); err != nil {
				log.Printf("error unset env var %s\n", key)
				return
			}
		}
		if value.NeedRemove {
			continue
		}
		if err := os.Setenv(key, value.Value); err != nil {
			log.Printf("error set env var %s\n", key)
			return
		}
	}

	execCmd := exec.Command(cmd[0], cmd[1:]...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stderr = os.Stderr

	if err := execCmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			returnCode = exitErr.ExitCode()
		} else {
			log.Println(err.Error())
		}
	} else {
		returnCode = 0
	}

	return
}
