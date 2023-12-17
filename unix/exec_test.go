// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package unix

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"

	"github.com/sergeymakinen/go-quote/internal/testutil"
)

func TestSingleQuote_Quote_Exec(t *testing.T) {
	for _, it := range testutil.InputTests('\'', '\t', '\n', ' ', '"') {
		it := it
		t.Run(it.Name, func(t *testing.T) {
			t.Parallel()
			testutil.TestExecOutput(t, it.Input, "sh", "-c", `printf '%s\n' `+SingleQuote.Quote(it.Input))
		})
	}
}

func TestDoubleQuote_Quote_Exec(t *testing.T) {
	for _, it := range testutil.InputTests('"', '\t', '\n', ' ', '$', '\'') {
		it := it
		t.Run(it.Name, func(t *testing.T) {
			t.Parallel()
			testutil.TestExecOutput(t, strings.ReplaceAll(it.Input, "!", `\!`), "/bin/sh", "-c", `printf '%s\n' `+DoubleQuote.Quote(it.Input))
		})
	}
}

func TestANSIC_Quote_Exec(t *testing.T) {
	if ansiCShell == "" {
		t.Skip(`no shell with \uxxxx support`)
	}
	for _, it := range testutil.InputTests('\'', '\t', '\n', ' ', '$', '"') {
		it := it
		t.Run(it.Name, func(t *testing.T) {
			if strings.Contains(it.Name, "bytes") {
				t.Skipf("it.Name=%s", it.Name)
			}
			t.Parallel()
			testutil.TestExecOutput(t, it.Input, ansiCShell, "-c", `printf '%s\n' `+ANSIC.Quote(it.Input))
		})
	}
}

func TestANSIC_QuoteBinary_Exec(t *testing.T) {
	if ansiCShell == "" {
		t.Skip(`no shell with \uxxxx support`)
	}
	for _, it := range testutil.InputTests('\'', '\t', '\n', ' ', '$', '"') {
		it := it
		t.Run(it.Name, func(t *testing.T) {
			t.Parallel()
			testutil.TestExecOutput(t, it.Input, ansiCShell, "-c", `printf '%s\n' `+ANSIC.QuoteBinary([]byte(it.Input)))
		})
	}
}

var ansiCShell string

func init() {
	// macOS uses an ancient Bash, so testing for \uxxxx support
	for _, s := range []string{"bash", "zsh"} {
		out, err := exec.Command(s, "-c", `printf '%s\n' $'\u200D'`).Output()
		if err == nil && !bytes.Contains(out, []byte("$")) && !bytes.Contains(out, []byte(`200D`)) {
			ansiCShell = s
			break
		}
	}
}
