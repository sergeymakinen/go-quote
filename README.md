# quote

[![tests](https://github.com/sergeymakinen/go-quote/workflows/tests/badge.svg)](https://github.com/sergeymakinen/go-quote/actions?query=workflow%3Atests)
[![Go Reference](https://pkg.go.dev/badge/github.com/sergeymakinen/go-quote.svg)](https://pkg.go.dev/github.com/sergeymakinen/go-quote)
[![Go Report Card](https://goreportcard.com/badge/github.com/sergeymakinen/go-quote)](https://goreportcard.com/report/github.com/sergeymakinen/go-quote)
[![codecov](https://codecov.io/gh/sergeymakinen/go-quote/branch/master/graph/badge.svg)](https://codecov.io/gh/sergeymakinen/go-quote)

Package quote defines interfaces shared by other packages
that quote command-line arguments and variables.

See the documentation for the [unix](https://pkg.go.dev/github.com/sergeymakinen/go-quote/unix) and [windows](https://pkg.go.dev/github.com/sergeymakinen/go-quote/windows) packages for more information.

## Installation

Use go get:

```bash
go get github.com/sergeymakinen/go-quote
```

Then import the package into your own code:

```go
import "github.com/sergeymakinen/go-quote"
```


## Example

### Unix

```go
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
```

### Windows

```go
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
```

## License

BSD 3-Clause
