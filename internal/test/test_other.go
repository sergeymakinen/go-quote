// +build !windows

package test

import (
	"bytes"
	"os/exec"
)

func Output(name string, args ...string) ([]byte, []string, error) {
	out, err := exec.Command(name, args...).Output()
	arg := append([]string{name}, args...)
	if err != nil {
		return nil, arg, err
	}
	out = bytes.TrimSuffix(out, []byte("\n"))
	return out, arg, nil
}
