package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Get all environment variables (system envs without unsetted + user envs).
func getSystemAndUserEnvs(envs Environment) []string {
	envsToRemove := map[string]string{}
	for envName, envVal := range envs {
		if envVal.NeedRemove {
			envsToRemove[envName] = envVal.Value
		}
	}

	finalEnvs := []string{}
	for _, env := range os.Environ() {
		split := strings.Split(env, "=")
		envName := split[0]

		if _, exists := envsToRemove[envName]; !exists {
			finalEnvs = append(finalEnvs, env)
		}
	}

	for envName, envVal := range envs {
		if !envVal.NeedRemove {
			finalEnvs = append(finalEnvs, fmt.Sprintf("%s=%s", envName, envVal.Value))
		}
	}

	return finalEnvs
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, envs Environment) (returnCode int) {
	if len(cmd) < 1 {
		return 1
	}
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	command.Env = getSystemAndUserEnvs(envs)

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout

	if err := command.Run(); err != nil {
		return 1
	}

	return 0
}
