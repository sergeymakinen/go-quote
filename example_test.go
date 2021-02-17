package quote_test

import (
	"fmt"

	"github.com/sergeymakinen/go-quote/unix"
	"github.com/sergeymakinen/go-quote/windows"
)

func ExampleQuoting_unix() {
	filename := `Long File With 'Single' & "Double" Quotes.txt`
	fmt.Println(unix.SingleQuote.MustQuote(filename))
	// Echoing inside 'sh -c' requires a quoting to make it safe with an arbitrary string
	quoted := unix.SingleQuote.Quote(filename)
	fmt.Println([]string{
		"sh",
		"-c",
		fmt.Sprintf("echo %s | callme", quoted),
	})
	unquoted, _ := unix.SingleQuote.Unquote(quoted)
	fmt.Println(unquoted)
	// Output:
	// true
	// [sh -c echo 'Long File With '"'"'Single'"'"' & "Double" Quotes.txt' | callme]
	// Long File With 'Single' & "Double" Quotes.txt
}

func ExampleQuoting_windows() {
	// If you have to deal with a different from Argv command-line quoting
	// when starting processes on Windows, don't forget to manually create a command-line
	// via the CmdLine SysProcAttr attribute:
	//
	//  cmd := exec.Command(name)
	//  cmd.SysProcAttr = &windows.SysProcAttr{
	//  	CmdLine: strings.Join(args, " "),
	//  }
	filename := `Long File With 'Single' & "Double" Quotes.txt`
	fmt.Println(windows.Argv.MustQuote(filename))
	// Using both Argv and Cmd quoting as callme.exe requires the Argv quoting
	// and its safe usage in cmd.exe requires the Cmd quoting
	quoted := windows.Argv.Quote(filename)
	fmt.Println([]string{
		"cmd.exe",
		"/C",
		fmt.Sprintf("callme.exe %s", windows.Cmd.Quote(quoted)),
	})
	unquoted, _ := windows.Cmd.Unquote(quoted)
	unquoted, _ = windows.Argv.Unquote(unquoted)
	fmt.Println(unquoted)
	// Output:
	// true
	// [cmd.exe /C callme.exe ^"Long^ File^ With^ ^'Single^'^ ^&^ \^"Double\^"^ Quotes.txt^"]
	// Long File With 'Single' & "Double" Quotes.txt
}
