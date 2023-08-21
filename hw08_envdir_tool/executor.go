package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Get all environment variables (system envs without unsetted + user envs).
func getSystemAndUserEnvs(envs Environment) []string {
	finalEnvs := []string{}
	for _, env := range os.Environ() {
		split := strings.Split(env, "=")
		envName := split[0]

		if envVal, exists := envs[envName]; exists {
			// skip env if NeedRemove is true
			if !envVal.NeedRemove {
				finalEnvs = append(finalEnvs, fmt.Sprintf("%s=%s", envName, envVal.Value))
			}
			// delete processed envs
			delete(envs, envName)
		} else {
			finalEnvs = append(finalEnvs, env)
		}
	}

	// append remaining envs
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
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}

		return 1
	}

	return 0
}
