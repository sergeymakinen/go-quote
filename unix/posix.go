// Package unix contains quoting interfaces for Unix-based shells.
package unix

import (
	"fmt"
	"strings"

	"github.com/sergeymakinen/go-quote"
)

type singleQuote struct {
	unixQuote
}

func (singleQuote) Quote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'"'"'`) + "'"
}

func (singleQuote) Unquote(s string) (string, error) {
	var (
		buf                          strings.Builder
		inSingleQuote, inDoubleQuote bool
	)
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\'':
			if inDoubleQuote {
				buf.WriteByte(s[i])
			} else {
				inSingleQuote = !inSingleQuote
			}
		case '"':
			if inSingleQuote {
				buf.WriteByte(s[i])
			} else {
				inDoubleQuote = !inDoubleQuote
			}
		default:
			if inDoubleQuote {
				return "", &quote.SyntaxError{
					Msg:    fmt.Sprintf("unsupported character %#U in double quoted string", s[i]),
					Offset: i + 1,
				}
			}
			if !inSingleQuote {
				return "", &quote.SyntaxError{
					Msg:    fmt.Sprintf("character %#U outside of quoted string", s[i]),
					Offset: i + 1,
				}
			}
			buf.WriteByte(s[i])
		}
	}
	if inSingleQuote || inDoubleQuote {
		return "", &quote.SyntaxError{
			Msg:    "unterminated quoted string",
			Offset: len(s),
		}
	}
	return buf.String(), nil
}

// SingleQuote quotes and unquotes strings, surrounded by single quotes (')
// as specified by POSIX.
//
// For example, the following string:
//
//  a b:"c d" 'e''f'  "g\""
//
// Would be quoted as:
//
//  'a b:"c d" '"'"'e'"'"''"'"'f'"'"'  "g\""'
//
// See https://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html#tag_18_02_02
// for details.
var SingleQuote quote.Quoting = singleQuote{}

type doubleQuote struct {
	unixQuote
}

var doubleQuoteReplacer = strings.NewReplacer(
	"!", `\!`,
	`"`, `\"`,
	"$", `\$`,
	`\`, `\\`,
	"`", "\\`",
)

func (doubleQuote) Quote(s string) string {
	return `"` + doubleQuoteReplacer.Replace(s) + `"`
}

func (doubleQuote) Unquote(s string) (string, error) {
	var (
		buf     strings.Builder
		inQuote bool
	)
	for i := 0; i < len(s); i++ {
		if s[i] == '"' {
			inQuote = !inQuote
			continue
		}
		if !inQuote {
			return "", &quote.SyntaxError{
				Msg:    fmt.Sprintf("character %#U outside of quoted string", s[i]),
				Offset: i + 1,
			}
		}
		escape := false
		if s[i] == '\\' {
			escape = true
			if i++; i >= len(s) {
				return "", &quote.SyntaxError{
					Msg:    "unterminated escape sequence",
					Offset: len(s),
				}
			}
		}
		switch s[i] {
		case '!', '"', '$', '\\', '`':
			if !escape {
				return "", &quote.SyntaxError{
					Msg:    fmt.Sprintf("unescaped special character %#U", s[i]),
					Offset: i + 1,
				}
			}
			escape = false
			fallthrough
		default:
			if escape {
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

// DoubleQuote quotes and unquotes strings, surrounded by double quotes (")
// as specified by POSIX.
//
// For example, the following string:
//
//  a b:"c d" 'e''f'  "g\""
//
// Would be quoted as:
//
//  "a b:\"c d\" 'e''f'  \"g\\\"\""
//
// See https://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html#tag_18_02_03
// for details.
var DoubleQuote quote.Quoting = doubleQuote{}
