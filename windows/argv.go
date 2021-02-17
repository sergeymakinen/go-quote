// Package windows contains quoting interfaces for Windows shells and programs.
package windows

import (
	"fmt"
	"strings"

	"github.com/sergeymakinen/go-quote"
)

const argvUnsafeChars = "\t \""

type argv struct{}

func (argv) MustQuote(s string) bool {
	return strings.ContainsAny(s, argvUnsafeChars)
}

func (argv) Quote(s string) string {
	var (
		buf     strings.Builder
		slashes int
	)
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"':
			for slashes++; slashes > 0; slashes-- {
				buf.WriteByte('\\')
			}
			buf.WriteByte(s[i])
		case '\\':
			slashes++
			buf.WriteByte(s[i])
		default:
			slashes = 0
			buf.WriteByte(s[i])
		}
	}
	for ; slashes > 0; slashes-- {
		buf.WriteByte('\\')
	}
	return `"` + buf.String() + `"`
}

func (a argv) Unquote(s string) (string, error) {
	var (
		buf     strings.Builder
		inQuote bool
		slashes int
	)
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"':
			if slashes > 0 {
				if slashes%2 == 0 {
					for ; slashes > 0; slashes -= 2 {
						buf.WriteByte('\\')
					}
					inQuote = !inQuote
				} else {
					for slashes--; slashes > 0; slashes -= 2 {
						buf.WriteByte('\\')
					}
					buf.WriteByte(s[i])
				}
			} else {
				inQuote = !inQuote
			}
		case '\\':
			if !inQuote {
				return "", &quote.SyntaxError{
					Msg:    fmt.Sprintf("character %#U outside of quoted string", s[i]),
					Offset: i + 1,
				}
			}
			slashes++
		default:
			if !inQuote {
				return "", &quote.SyntaxError{
					Msg:    fmt.Sprintf("character %#U outside of quoted string", s[i]),
					Offset: i + 1,
				}
			}
			for ; slashes > 0; slashes-- {
				buf.WriteByte('\\')
			}
			buf.WriteByte(s[i])
		}
	}
	if inQuote {
		return "", &quote.SyntaxError{
			Msg:    "unterminated quoted string",
			Offset: len(s),
		}
	}
	return buf.String(), nil
}

// Argv quotes and unquotes strings, surrounded by double quotes ("…"),
// as specified by Microsoft for the CommandLineToArgvW function.
//
// For example, the following string:
//
//  a b:"c d" 'e''f'  "g\""
//
// Would be quoted as:
//
//  "a b:\"c d\" 'e''f'  \"g\\\"\""
//
// See https://docs.microsoft.com/en-us/cpp/c-language/parsing-c-command-line-arguments?view=msvc-160
// for details.
var Argv quote.Quoting = argv{}
