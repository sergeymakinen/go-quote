package windows

import (
	"fmt"
	"testing"

	"github.com/sergeymakinen/go-quote/internal/testutil"
)

func Example() {
	// If you have to deal with a different from Argv command-line quoting
	// when starting processes on Windows, don't forget to manually create a command-line
	// via the CmdLine SysProcAttr attribute:
	//
	//  cmd := exec.Command(name)
	//  cmd.SysProcAttr = &windows.SysProcAttr{
	//  	CmdLine: strings.Join(args, " "),
	//  }
	filename := `Long File With 'Single' & "Double" Quotes.txt`
	fmt.Println(Argv.MustQuote(filename))
	// Using both Argv and Cmd quoting as callme.exe requires the Argv quoting
	// and its safe usage in cmd.exe requires the Cmd quoting
	quoted := Argv.Quote(filename)
	fmt.Println([]string{
		"cmd.exe",
		"/C",
		fmt.Sprintf("callme.exe %s", Cmd.Quote(quoted)),
	})
	unquoted, _ := Cmd.Unquote(quoted)
	unquoted, _ = Argv.Unquote(unquoted)
	fmt.Println(unquoted)
	// Output:
	// true
	// [cmd.exe /C callme.exe ^"Long^ File^ With^ ^'Single^'^ ^&^ \^"Double\^"^ Quotes.txt^"]
	// Long File With 'Single' & "Double" Quotes.txt
}

func TestCmd_Quote_Unquote(t *testing.T) {
	tests := []struct {
		Name, Input, Output string
	}{
		{
			Name:   "empty string",
			Input:  "",
			Output: "",
		},
		{
			Name:   "special char escaping",
			Input:  "!\"&'+,;<=>[]^`{}~",
			Output: "^!^\"^&^'^+^,^;^<^=^>^[^]^^^`^{^}^~",
		},
	}
	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			quoted := Cmd.Quote(td.Input)
			testutil.TestDiff(t, "Cmd.Quote() ", td.Output, quoted)
			unquoted, err := Cmd.Unquote(quoted)
			if err != nil {
				t.Fatalf("Cmd.Unquote() = _, %v; want nil", err)
			}
			testutil.TestDiff(t, "Cmd.Unquote()", td.Input, unquoted)
		})
	}
}

func TestCmd_Quote_Unquote_InputTests(t *testing.T) {
	for _, it := range testutil.InputTests('"') {
		t.Run(it.Name, func(t *testing.T) {
			quoted := Cmd.Quote(it.Input)
			unquoted, err := Cmd.Unquote(quoted)
			if err != nil {
				t.Fatalf("Cmd.Unquote() = _, %v; want nil", err)
			}
			testutil.TestDiff(t, "Cmd.Unquote()", it.Input, unquoted)
		})
	}
}
