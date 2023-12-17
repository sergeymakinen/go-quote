package testutil

import (
	"bytes"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows"
)

func Output(name string, args ...string) ([]byte, []string, error) {
	cmd := exec.Command(name, args...)
	var args2 []string
	for i, arg := range append([]string{cmd.Path}, args...) {
		if i < len(args) {
			args2 = append(args2, windows.EscapeArg(arg))
		} else {
			args2 = append(args2, arg)
		}
	}
	cmd.SysProcAttr = &windows.SysProcAttr{
		CmdLine: strings.Join(args2, " "),
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, args2, err
	}
	out = bytes.TrimSuffix(out, []byte("\r\n"))
	return out, args2, nil
}
