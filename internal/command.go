package internal

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
)

func RunCmd(command string, args []string, envs []string) ([]byte, error) {

	args = append(args, "--json")

	cmd := exec.Command(command, args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, envs...)

	var cmdOut, cmdErr bytes.Buffer

	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr

	if err := cmd.Run(); err != nil {
		return nil, errors.New(cmdErr.String())
	}

	return cmdOut.Bytes(), nil
}
