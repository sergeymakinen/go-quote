package windows

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/sergeymakinen/go-quote"
)

const psUnsafeChars = "\t \"$'`"

type psQuote struct{}

func (psQuote) MustQuote(s string) bool {
	return strings.ContainsAny(s, psUnsafeChars)
}

type psSingleQuote struct {
	psQuote
}

func (psSingleQuote) Quote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

func (psSingleQuote) Unquote(s string) (string, error) {
	var (
		buf     strings.Builder
		inQuote bool
	)
	for i := 0; i < len(s); i++ {
		if s[i] == '\'' {
			if !inQuote {
				inQuote = true
			} else {
				if i+1 < len(s) && s[i+1] == '\'' {
					buf.WriteByte('\'')
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

// PSSingleQuote quotes and unquotes strings, surrounded by single quotes ('…')
// containing characters special to Windows PowerShell (powershell.exe).
//
// For example, the following string:
//
//  a b:"c d" 'e''f'  "g\""
//
// Would be quoted as:
//
//  'a b:"c d" ''e''''f''  "g\""'
//
// See https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_quoting_rules?view=powershell-7.1
// for details.
var PSSingleQuote quote.Quoting = psSingleQuote{}

type basePSDoubleQuote struct {
	psQuote
}

func (basePSDoubleQuote) Unquote(s string) (string, error) {
	var (
		r       rune
		inQuote bool
		buf     strings.Builder
		b       = make([]rune, 6)
	)
	for i, width := 0, 0; i < len(s); i += width {
		r, width = utf8.DecodeRuneInString(s[i:])
		if r == '"' {
			inQuote = !inQuote
			continue
		}
		if !inQuote {
			return "", &quote.SyntaxError{
				Msg:    fmt.Sprintf("character %#U outside of quoted string", r),
				Offset: i + 1,
			}
		}
		if r != '`' {
			switch r {
			case '$':
				return "", &quote.SyntaxError{
					Msg:    fmt.Sprintf("unescaped special character %#U", r),
					Offset: i + 1,
				}
			default:
				buf.WriteRune(r)
			}
			continue
		}
		if i += width; i >= len(s) {
			return "", &quote.SyntaxError{
				Msg:    "unterminated escape sequence",
				Offset: len(s),
			}
		}
		r, width = utf8.DecodeRuneInString(s[i:])
		switch r {
		case '0':
			buf.WriteRune('\000')
		case 'a':
			buf.WriteRune('\a')
		case 'b':
			buf.WriteRune('\b')
		case 'e':
			buf.WriteRune('\x1B')
		case 'f':
			buf.WriteRune('\f')
		case 'n':
			buf.WriteRune('\n')
		case 'r':
			buf.WriteRune('\r')
		case 't':
			buf.WriteRune('\t')
		case 'v':
			buf.WriteRune('\v')
		case '"', '$', '`':
			buf.WriteRune(r)
		case 'u':
			if i += width; i >= len(s) {
				return "", &quote.SyntaxError{
					Msg:    "unterminated escape sequence `u",
					Offset: len(s),
				}
			}
			r, width = utf8.DecodeRuneInString(s[i:])
			if r != '{' {
				return "", &quote.SyntaxError{
					Msg:    fmt.Sprintf("invalid character %#U in escape sequence '`u'", r),
					Offset: i + 1,
				}
			}
			b = b[:0]
			for n := 6; n > 0 && i+width < len(s); n-- {
				r1, width1 := utf8.DecodeRuneInString(s[i+width:])
				if (r1 >= '0' && r1 <= '9') || (r1 >= 'a' && r1 <= 'f') || (r1 >= 'A' && r1 <= 'F') {
					b = append(b, r1)
					i += width
					width = width1
				} else {
					break
				}
			}
			if len(b) == 0 {
				return "", &quote.SyntaxError{
					Msg:    "invalid escape sequence '`u'",
					Offset: i + 1,
				}
			}
			if i += width; i >= len(s) {
				return "", &quote.SyntaxError{
					Msg:    "unterminated escape sequence `u",
					Offset: len(s),
				}
			}
			r, width = utf8.DecodeRuneInString(s[i:])
			if r != '}' {
				return "", &quote.SyntaxError{
					Msg:    fmt.Sprintf("invalid character %#U in escape sequence '`u'", r),
					Offset: i + 1,
				}
			}
			v, err := strconv.ParseUint(string(b), 16, 48)
			if err != nil || v > utf8.MaxRune {
				return "", &quote.SyntaxError{
					Msg:    fmt.Sprintf("invalid escape sequence '`u{%s}'", string(b)),
					Offset: i + 1,
				}
			}
			buf.WriteRune(rune(v))
		default:
			buf.WriteRune(r)
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

func (basePSDoubleQuote) quote(s string, pwsh bool) string {
	var buf strings.Builder
	for _, r := range s {
		switch r {
		case '\000':
			buf.WriteString("`0")
		case '\a':
			buf.WriteString("`a")
		case '\b':
			buf.WriteString("`b")
		case '\x1B':
			if pwsh {
				buf.WriteString("`e")
			} else {
				buf.WriteRune(r)
			}
		case '\f':
			buf.WriteString("`f")
		case '\n':
			buf.WriteString("`n")
		case '\r':
			buf.WriteString("`r")
		case '\t':
			buf.WriteString("`t")
		case '\v':
			buf.WriteString("`v")
		case '"', '$', '`':
			buf.WriteRune('`')
			buf.WriteRune(r)
		default:
			switch {
			case pwsh && (r < 0x20 || !strconv.IsPrint(r)):
				switch {
				case r < 0x7F:
					buf.WriteString(fmt.Sprintf("`u{%02X}", r))
				case r < 0x10000:
					buf.WriteString(fmt.Sprintf("`u{%04X}", r))
				default:
					buf.WriteString(fmt.Sprintf("`u{%06X}", r))
				}
			default:
				buf.WriteRune(r)
			}
		}
	}
	return `"` + buf.String() + `"`
}

type psDoubleQuote struct {
	basePSDoubleQuote
}

func (psDoubleQuote) MustQuote(s string) bool {
	return strings.ContainsAny(s, psUnsafeChars)
}

func (q psDoubleQuote) Quote(s string) string {
	return q.quote(s, false)
}

// PSDoubleQuote quotes and unquotes strings, surrounded by double quotes ("…")
// containing characters special to Windows PowerShell (powershell.exe).
//
// For example, the following string:
//
//  a b:"c d" 'e''f'  "g\""
//
// Would be quoted as:
//
//  "a b:`"c d`" 'e''f'  `"g\`"`""
//
// See https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_quoting_rules?view=powershell-7.1
// and https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_special_characters?view=powershell-7.1
// for details.
var PSDoubleQuote quote.Quoting = psDoubleQuote{}

type pwshDoubleQuote struct {
	basePSDoubleQuote
}

func (pwshDoubleQuote) MustQuote(s string) bool {
	return strings.ContainsAny(s, psUnsafeChars)
}

func (q pwshDoubleQuote) Quote(s string) string {
	return q.quote(s, true)
}

// PwshDoubleQuote quotes and unquotes strings, surrounded by double quotes ("…")
// containing characters special to PowerShell/PowerShell Core (pwsh.exe).
//
// For example, the following string:
//
//  a b:"c d" 'e''f'  "g\""
//
// Would be quoted as:
//
//  "a b:`"c d`" 'e''f'  `"g\`"`""
//
// See https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_quoting_rules?view=powershell-7.1
// and https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_special_characters?view=powershell-7.1
// for details.
var PwshDoubleQuote quote.Quoting = pwshDoubleQuote{}
