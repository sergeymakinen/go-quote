// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package unix

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"

	"github.com/sergeymakinen/go-quote/internal/test"
)

func TestSingleQuote_Quote_Exec(t *testing.T) {
	for _, it := range test.InputTests('\'', '\t', '\n', ' ', '"') {
		it := it
		t.Run(it.Name, func(t *testing.T) {
			t.Parallel()
			test.TestExecOutput(t, it.Input, "sh", "-c", `printf '%s\n' `+SingleQuote.Quote(it.Input))
		})
	}
}

func TestDoubleQuote_Quote_Exec(t *testing.T) {
	for _, it := range test.InputTests('"', '\t', '\n', ' ', '$', '\'') {
		it := it
		t.Run(it.Name, func(t *testing.T) {
			t.Parallel()
			test.TestExecOutput(t, strings.ReplaceAll(it.Input, "!", `\!`), "/bin/sh", "-c", `printf '%s\n' `+DoubleQuote.Quote(it.Input))
		})
	}
}

func TestANSIC_Quote_Exec(t *testing.T) {
	// macOS uses an ancient Bash, so testing for \uxxxx support
	var shell string
	for _, s := range []string{"sh", "bash", "zsh"} {
		out, err := exec.Command(s, "-c", `printf '%s\n' $'\u200D'`).Output()
		if err == nil && !bytes.Contains(out, []byte(`200D`)) {
			shell = s
			break
		}
	}
	if shell == "" {
		t.Skip(`no shell with \uxxxx support`)
	}
	for _, it := range test.InputTests('\'', '\t', '\n', ' ', '$', '"') {
		it := it
		t.Run(it.Name, func(t *testing.T) {
			if strings.Contains(it.Name, "bytes") {
				t.Skipf("it.Name=%s", it.Name)
			}
			t.Parallel()
			test.TestExecOutput(t, it.Input, shell, "-c", `printf '%s\n' `+ANSIC.Quote(it.Input))
		})
	}
}

func TestANSIC_QuoteBinary_Exec(t *testing.T) {
	for _, it := range test.InputTests('\'', '\t', '\n', ' ', '$', '"') {
		it := it
		t.Run(it.Name, func(t *testing.T) {
			t.Parallel()
			test.TestExecOutput(t, it.Input, "sh", "-c", `printf '%s\n' `+ANSIC.QuoteBinary([]byte(it.Input)))
		})
	}
}
