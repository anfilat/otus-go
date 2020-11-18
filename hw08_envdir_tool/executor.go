package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	//nolint:gosec
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Env = changeEnv(env)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
		return selfReturnCode
	}
	return 0
}

func changeEnv(env Environment) []string {
	result := os.Environ()
	for name, v := range env {
		result = filterEnv(result, name)
		if !v.Remove {
			result = append(result, name+"="+v.Value)
		}
	}
	return result
}

func filterEnv(envVars []string, delName string) []string {
	n := 0
	for _, v := range envVars {
		name := strings.SplitN(v, "=", 2)[0]
		if name != delName {
			envVars[n] = v
			n++
		}
	}
	return envVars[:n]
}
