package windows

import (
	"fmt"
	"strings"

	"github.com/sergeymakinen/go-quote"
)

type msiexec struct{}

func (msiexec) MustQuote(s string) bool {
	return strings.ContainsAny(s, argvUnsafeChars)
}

func (msiexec) Quote(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

func (msiexec) Unquote(s string) (string, error) {
	var (
		buf     strings.Builder
		inQuote bool
	)
	for i := 0; i < len(s); i++ {
		if s[i] == '"' {
			if !inQuote {
				inQuote = true
			} else {
				if i+1 < len(s) && s[i+1] == '"' {
					buf.WriteByte('"')
					i++
				} else {
					inQuote = false
				}
			}
			continue
		}
		if !inQuote {
			return "", &quote.SyntaxError{
				Msg:    fmt.Sprintf("character %#U outside of quoted string", s[i]),
				Offset: i + 1,
			}
		}
		buf.WriteByte(s[i])
	}
	if inQuote {
		return "", &quote.SyntaxError{
			Msg:    "unterminated quoted string",
			Offset: len(s),
		}
	}
	return buf.String(), nil
}

// Msiexec quotes and unquotes strings, surrounded by double quotes ("…")
// containing characters special to the Windows Installer (msiexec.exe).
//
// For example, the following string:
//
//  a b:"c d" 'e''f'  "g\""
//
// Would be quoted as:
//
//  "a b:""c d"" 'e''f'  ""g\"""""
//
// See https://docs.microsoft.com/en-us/windows/win32/msi/command-line-options
// for details.
var Msiexec quote.Quoting = msiexec{}
