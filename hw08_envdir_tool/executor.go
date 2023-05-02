package main

import (
	"errors"
	"os"
	"os/exec"
)

func updateEnvironmentVariables(env Environment) {
	if env == nil {
		return
	}

	for key, value := range env {
		if value.NeedRemove {
			os.Unsetenv(key)
			continue
		}
		os.Setenv(key, value.Value)
	}
}

func getExitError(exitError *exec.ExitError) int {
	if exitError != nil {
		return exitError.ExitCode()
	}
	return 0
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 0
	}

	command := cmd[0]
	args := make([]string, 0)
	if len(cmd) > 1 {
		args = cmd[1:]
	}
	updateEnvironmentVariables(env)

	executable := exec.Command(command, args...)
	executable.Stderr = os.Stderr
	executable.Stdout = os.Stdout
	executable.Stdin = os.Stdin
	err := executable.Run()
	var exitErr *exec.ExitError
	if err != nil && errors.As(err, &exitErr) {
		return getExitError(exitErr)
	}
	return 0
}
