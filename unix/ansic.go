package unix

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/sergeymakinen/go-quote"
)

type ansiC struct {
	unixQuote
}

func (ansiC) Quote(s string) string {
	var buf strings.Builder
	for _, r := range s {
		switch r {
		case '\a':
			buf.WriteString(`\a`)
		case '\b':
			buf.WriteString(`\b`)
		case '\x1B':
			buf.WriteString(`\e`)
		case '\f':
			buf.WriteString(`\f`)
		case '\n':
			buf.WriteString(`\n`)
		case '\r':
			buf.WriteString(`\r`)
		case '\t':
			buf.WriteString(`\t`)
		case '\v':
			buf.WriteString(`\v`)
		case '"', '\'', '?', '\\':
			buf.WriteRune('\\')
			buf.WriteRune(r)
		default:
			switch {
			case r < 0x20:
				buf.WriteString(fmt.Sprintf("\\x%02X", r))
			case strconv.IsPrint(r):
				buf.WriteRune(r)
			default:
				if r < 0x10000 {
					buf.WriteString(fmt.Sprintf("\\u%04X", r))
				} else {
					buf.WriteString(fmt.Sprintf("\\U%08X", r))
				}
			}
		}
	}
	return "$'" + buf.String() + "'"
}

func (ansiC) Unquote(s string) (string, error) {
	var (
		r       rune
		inQuote bool
		buf     strings.Builder
		b       = make([]rune, 9)
	)
	for i, width := 0, 0; i < len(s); i += width {
		r, width = utf8.DecodeRuneInString(s[i:])
		if !inQuote {
			if strings.HasPrefix(s[i:], "$'") {
				inQuote = true
				width++
				continue
			}
			return "", &quote.SyntaxError{
				Msg:    fmt.Sprintf("character %#U outside of quoted string", r),
				Offset: i + 1,
			}
		} else if r == '\'' {
			inQuote = false
			continue
		}
		if r != '\\' {
			buf.WriteRune(r)
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
		case 'a':
			buf.WriteRune('\a')
		case 'b':
			buf.WriteRune('\b')
		case 'e':
			buf.WriteRune('\x1B')
		case 'E':
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
		case '"', '?', '\'', '\\':
			buf.WriteRune(r)
		case 'c':
			if i += width; i >= len(s) {
				return "", &quote.SyntaxError{
					Msg:    "unterminated escape sequence `\\c`",
					Offset: len(s),
				}
			}
			r, width = utf8.DecodeRuneInString(s[i:])
			switch {
			case r == '?':
				buf.WriteByte(0x7F)
			case r >= 'a' && r <= 'z':
				r -= 'a' - 'A'
				fallthrough
			case r >= '@' && r <= '_':
				buf.WriteRune(r - '@')
			default:
				return "", &quote.SyntaxError{
					Msg:    fmt.Sprintf("invalid character %#U in escape sequence `\\c`", r),
					Offset: i + 1,
				}
			}
		case 'x', 'u', 'U':
			n := 0
			switch s[i] {
			case 'x':
				n = 2
			case 'u':
				n = 4
			case 'U':
				n = 8
			}
			b = append(b[:0], r)
			for ; n > 0 && i+width < len(s); n-- {
				r1, width1 := utf8.DecodeRuneInString(s[i+width:])
				if (r1 >= '0' && r1 <= '9') || (r1 >= 'a' && r1 <= 'f') || (r1 >= 'A' && r1 <= 'F') {
					b = append(b, r1)
					i += width
					width = width1
				} else {
					break
				}
			}
			if len(b) == 1 {
				return "", &quote.SyntaxError{
					Msg:    fmt.Sprintf("unterminated escape sequence `\\%s`", string(b)),
					Offset: i + 1,
				}
			}
			v, err := strconv.ParseUint(string(b[1:]), 16, 64)
			if err != nil || v > utf8.MaxRune {
				return "", &quote.SyntaxError{
					Msg:    fmt.Sprintf("invalid escape sequence `\\%s`", string(b)),
					Offset: i + 1,
				}
			}
			buf.WriteRune(rune(v))
		case '0', '1', '2', '3', '4', '5', '6', '7':
			b = append(b[:0], r)
			for n := 2; n > 0 && i+width < len(s); n-- {
				r1, width1 := utf8.DecodeRuneInString(s[i+width:])
				if r1 >= '0' && r1 <= '7' {
					b = append(b, r1)
					i += width
					width = width1
				} else {
					break
				}
			}
			v, err := strconv.ParseUint(string(b), 8, 8)
			if err != nil {
				return "", &quote.SyntaxError{
					Msg:    fmt.Sprintf("invalid escape sequence `\\%s`", string(b)),
					Offset: i + 1,
				}
			}
			buf.WriteByte(byte(v))
		default:
			buf.WriteRune('\\')
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

func (ansiC) QuoteBinary(b []byte) string {
	var buf strings.Builder
	for i := 0; i < len(b); i++ {
		switch b[i] {
		case '\a':
			buf.WriteString(`\a`)
		case '\b':
			buf.WriteString(`\b`)
		case '\x1B':
			buf.WriteString(`\e`)
		case '\f':
			buf.WriteString(`\f`)
		case '\n':
			buf.WriteString(`\n`)
		case '\r':
			buf.WriteString(`\r`)
		case '\t':
			buf.WriteString(`\t`)
		case '\v':
			buf.WriteString(`\v`)
		case '"', '\'', '?', '\\':
			buf.WriteString(`\`)
			buf.WriteByte(b[i])
		default:
			switch {
			case b[i] < 0x20 || b[i] > 0x7F || !strconv.IsPrint(rune(b[i])):
				buf.WriteString(fmt.Sprintf("\\x%02X", b[i]))
			default:
				buf.WriteByte(b[i])
			}
		}
	}
	return "$'" + buf.String() + "'"
}

func (ansiC) UnquoteBinary(s string) ([]byte, error) {
	var (
		buf     bytes.Buffer
		inQuote bool
		b       = make([]byte, 9)
	)
	for i := 0; i < len(s); i++ {
		if !inQuote {
			if strings.HasPrefix(s[i:], "$'") {
				inQuote = true
				i++
				continue
			}
			return nil, &quote.SyntaxError{
				Msg:    fmt.Sprintf("character %#U outside of quoted string", s[i]),
				Offset: i + 1,
			}
		} else if s[i] == '\'' {
			inQuote = false
			continue
		}
		if s[i] != '\\' {
			buf.WriteByte(s[i])
			continue
		}
		if i++; i >= len(s) {
			return nil, &quote.SyntaxError{
				Msg:    "unterminated escape sequence",
				Offset: len(s),
			}
		}
		switch s[i] {
		case 'a':
			buf.WriteByte('\a')
		case 'b':
			buf.WriteByte('\b')
		case 'e', 'E':
			buf.WriteByte('\x1B')
		case 'f':
			buf.WriteByte('\f')
		case 'n':
			buf.WriteByte('\n')
		case 'r':
			buf.WriteByte('\r')
		case 't':
			buf.WriteByte('\t')
		case 'v':
			buf.WriteByte('\v')
		case '"', '?', '\'', '\\':
			buf.WriteByte(s[i])
		case 'c':
			if i++; i >= len(s) {
				return nil, &quote.SyntaxError{
					Msg:    "unterminated escape sequence `\\c`",
					Offset: len(s),
				}
			}
			switch c := s[i]; {
			case s[i] == '?':
				buf.WriteByte(0x7F)
			case c >= 'a' && c <= 'z':
				c -= 'a' - 'A'
				fallthrough
			case c >= '@' && c <= '_':
				buf.WriteByte(c - '@')
			default:
				return nil, &quote.SyntaxError{
					Msg:    fmt.Sprintf("invalid character %#U in escape sequence `\\c`", c),
					Offset: i + 1,
				}
			}
		case 'x', 'u', 'U':
			n := 0
			isByte := false
			switch s[i] {
			case 'x':
				n = 2
				isByte = true
			case 'u':
				n = 4
			case 'U':
				n = 8
			}
			b = b[:1]
			b[0] = s[i]
			for i++; n > 0 && i < len(s); n-- {
				if (s[i] >= '0' && s[i] <= '9') || (s[i] >= 'a' && s[i] <= 'f') || (s[i] >= 'A' && s[i] <= 'F') {
					b = append(b, s[i])
					i++
				} else {
					break
				}
			}
			i--
			if len(b) == 1 {
				return nil, &quote.SyntaxError{
					Msg:    fmt.Sprintf("unterminated escape sequence `\\%s`", string(b)),
					Offset: i + 1,
				}
			}
			v, err := strconv.ParseUint(string(b[1:]), 16, 64)
			if err != nil || v > utf8.MaxRune {
				return nil, &quote.SyntaxError{
					Msg:    fmt.Sprintf("invalid escape sequence `\\%s`", string(b)),
					Offset: i + 1,
				}
			}
			if isByte {
				buf.WriteByte(byte(v))
			} else {
				buf.WriteRune(rune(v))
			}
		case '0', '1', '2', '3', '4', '5', '6', '7':
			b = b[:0]
			for n := 3; n > 0 && i < len(s); n-- {
				if s[i] >= '0' && s[i] <= '7' {
					b = append(b, s[i])
					i++
				} else {
					break
				}
			}
			i--
			v, err := strconv.ParseUint(string(b), 8, 8)
			if err != nil {
				return nil, &quote.SyntaxError{
					Msg:    fmt.Sprintf("invalid escape sequence `\\%s`", string(b)),
					Offset: i + 1,
				}
			}
			buf.WriteByte(byte(v))
		default:
			buf.WriteByte('\\')
			buf.WriteByte(s[i])
		}
	}
	if inQuote {
		return nil, &quote.SyntaxError{
			Msg:    "unterminated quoted string",
			Offset: len(s),
		}
	}
	return buf.Bytes(), nil
}

// ANSIC quotes and unquotes strings, surrounded by single quotes with a dollar sign prefix ($'…'),
// as specified by Bash and Zsh (ANSI-C quoting).
//
// For example, the following string:
//
//  a b:"c d" 'e''f'  "g\""
//
// Would be quoted as:
//
//  $'a b:\"c d\" \'e\'\'f\'  \"g\\\"\"'
//
// See https://www.gnu.org/software/bash/manual/html_node/ANSI_002dC-Quoting.html
// for details.
var ANSIC quote.BinaryQuoting = ansiC{}
