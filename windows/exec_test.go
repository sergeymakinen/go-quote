// +build windows

package windows

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sergeymakinen/go-quote/internal/testutil"
	"golang.org/x/sys/windows/registry"
)

func TestCmd_Quote_Exec(t *testing.T) {
	for _, it := range testutil.InputTests('"') {
		it := it
		t.Run(it.Name, func(t *testing.T) {
			if strings.HasPrefix(it.Name, "bytes:") || strings.Contains(it.Name, `\n`) {
				t.Skipf("Name=%s", it.Name)
			}
			t.Parallel()
			testutil.TestExecOutput(t, it.Input, "cmd.exe", "/c", `"chcp 65001 > NUL && echo `+Cmd.Quote(it.Input)+`"`)
		})
	}
}

func TestPSSingleQuote_Quote_Exec(t *testing.T) {
	for _, it := range testutil.InputTests('\'', '$', '`') {
		it := it
		t.Run(it.Name, func(t *testing.T) {
			if strings.HasPrefix(it.Name, "bytes:") {
				t.Skipf("Name=%s", it.Name)
			}
			t.Parallel()
			testutil.TestExecOutput(t, it.Input, "powershell.exe", "-NoProfile", "-NonInteractive", "-Command", Argv.Quote("$OutputEncoding = [Console]::OutputEncoding = [Text.UTF8Encoding]::UTF8; echo "+PSSingleQuote.Quote(it.Input)))
		})
	}
}

func TestPSDoubleQuote_Quote_Exec(t *testing.T) {
	for _, it := range testutil.InputTests('\'', '$', '`') {
		it := it
		t.Run(it.Name, func(t *testing.T) {
			if strings.HasPrefix(it.Name, "bytes:") {
				t.Skipf("Name=%s", it.Name)
			}
			t.Parallel()
			testutil.TestExecOutput(t, it.Input, "powershell.exe", "-NoProfile", "-NonInteractive", "-Command", Argv.Quote("$OutputEncoding = [Console]::OutputEncoding = [Text.UTF8Encoding]::UTF8; echo "+PSDoubleQuote.Quote(it.Input)))
		})
	}
}

func TestPwshDoubleQuote_Quote_Exec(t *testing.T) {
	if _, err := exec.LookPath("pwsh.exe"); err != nil {
		t.Skip(`no pwsh.exe`)
	}
	for _, it := range testutil.InputTests('\'', '$', '`') {
		it := it
		t.Run(it.Name, func(t *testing.T) {
			if strings.HasPrefix(it.Name, "bytes:") {
				t.Skipf("Name=%s", it.Name)
			}
			t.Parallel()
			testutil.TestExecOutput(t, it.Input, "pwsh.exe", "-NoProfile", "-NonInteractive", "-Command", Argv.Quote("$OutputEncoding = [Console]::OutputEncoding = [Text.UTF8Encoding]::UTF8; echo "+PwshDoubleQuote.Quote(it.Input)))
		})
	}
}

const envName = "GOQUOTETESTENV"

func TestArgv_Quote_Exec(t *testing.T) {
	if k, err := openEnvKey(); err == nil {
		k.DeleteValue(envName)
	}
	for _, it := range testutil.InputTests('"') {
		t.Run(it.Name, func(t *testing.T) {
			if strings.HasPrefix(it.Name, "bytes:") {
				t.Skipf("Name=%s", it.Name)
			}
			_, cmd, err := testutil.Output("setx.exe", envName, Argv.Quote(it.Input))
			if err != nil {
				t.Fatalf("Cmd.Output() = _, %v; want nil\nCmd: %v", err, cmd)
			}
			testutil.TestOutput(t, cmd, it.Input, envValue())
		})
	}
}

var envKey registry.Key

func openEnvKey() (registry.Key, error) {
	if envKey == 0 {
		k, err := registry.OpenKey(registry.CURRENT_USER, "Environment", registry.QUERY_VALUE|registry.SET_VALUE)
		if err != nil {
			return k, err
		}
		envKey = k
	}
	return envKey, nil
}

func envValue() string {
	k, err := openEnvKey()
	if err != nil {
		return ""
	}
	s, _, err := k.GetStringValue(envName)
	if err == nil {
		k.DeleteValue(envName)
	}
	return s
}

const msiProp = "PROPTEST"

func TestMsiexec_Quote_Exec(t *testing.T) {
	if k, err := openMsiKey(); err == nil {
		k.DeleteValue(msiProp)
	}
	wd, err := os.Getwd()
	if err != nil {
		panic("failed to get current directory: " + err.Error())
	}
	msiPath := filepath.Join(wd, "testdata/testutil.msi")
	for _, it := range testutil.InputTests('"') {
		t.Run(it.Name, func(t *testing.T) {
			if strings.HasPrefix(it.Name, "bytes:") {
				t.Skipf("Name=%s", it.Name)
			}
			_, cmd, err := testutil.Output("msiexec.exe", "/quiet", "/i", msiPath, msiProp+"="+Msiexec.Quote(it.Input))
			if err != nil {
				t.Fatalf("Cmd.Output() = _, %v; want nil\nCmd: %v", err, cmd)
			}
			testutil.TestOutput(t, cmd, it.Input, msiValue())
		})
	}
}

var msiKey registry.Key

func openMsiKey() (registry.Key, error) {
	if msiKey == 0 {
		k, err := registry.OpenKey(registry.CURRENT_USER, `Software\goquote`, registry.QUERY_VALUE|registry.SET_VALUE)
		if err != nil {
			return k, err
		}
		msiKey = k
	}
	return msiKey, nil
}

func msiValue() string {
	k, err := openMsiKey()
	if err != nil {
		return ""
	}
	s, _, err := k.GetStringValue(msiProp)
	if err == nil {
		k.DeleteValue(msiProp)
	}
	return s
}
