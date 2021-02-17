package windows

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sergeymakinen/go-quote"
	"github.com/sergeymakinen/go-quote/internal/test"
)

func TestPSSingleQuote_Quote(t *testing.T) {
	tests := []struct {
		Name, Input, Output string
	}{
		{
			Name:   "empty string",
			Input:  "",
			Output: "''",
		},
		{
			Name:   "special char escaping",
			Input:  "a'b",
			Output: "'a''b'",
		},
	}
	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			quoted := PSSingleQuote.Quote(td.Input)
			test.TestDiff(t, "PSSingleQuote.Quote() ", td.Output, quoted)
			unquoted, err := PSSingleQuote.Unquote(quoted)
			if err != nil {
				t.Fatalf("PSSingleQuote.Unquote() = _, %v; want nil", err)
			}
			test.TestDiff(t, "PSSingleQuote.Unquote()", td.Input, unquoted)
		})
	}
}

func TestPSSingleQuote_Unquote(t *testing.T) {
	unquoted, err := PSSingleQuote.Unquote("'ab''''''cd''ef'")
	if err != nil {
		t.Fatalf("PSSingleQuote.Unquote() = _, %v; want nil", err)
	}
	test.TestDiff(t, "PSSingleQuote.Unquote()", "ab'''cd'ef", unquoted)
}

func TestPSSingleQuote_Unquote_ShouldFail(t *testing.T) {
	tests := []struct {
		Name, Input string
		Err         error
	}{
		{
			Name:  "unterminated string #1",
			Input: "'a",
			Err: &quote.SyntaxError{
				Msg:    "unterminated quoted string",
				Offset: 2,
			},
		},
		{
			Name:  "unterminated string #2",
			Input: "'a''",
			Err: &quote.SyntaxError{
				Msg:    "unterminated quoted string",
				Offset: 4,
			},
		},
		{
			Name:  "char outside of string",
			Input: "a",
			Err: &quote.SyntaxError{
				Msg:    "character U+0061 'a' outside of quoted string",
				Offset: 1,
			},
		},
		{
			Name:  "char after string",
			Input: "'a'a",
			Err: &quote.SyntaxError{
				Msg:    "character U+0061 'a' outside of quoted string",
				Offset: 4,
			},
		},
	}
	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			_, err := PSSingleQuote.Unquote(td.Input)
			if diff := cmp.Diff(td.Err, err); diff != "" {
				t.Errorf("PSSingleQuote.Unquote() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPSSingleQuote_Quote_Unquote_InputTests(t *testing.T) {
	for _, it := range test.InputTests('\'', '$', '`') {
		t.Run(it.Name, func(t *testing.T) {
			quoted := PSSingleQuote.Quote(it.Input)
			unquoted, err := PSSingleQuote.Unquote(quoted)
			if err != nil {
				t.Fatalf("PSSingleQuote.Unquote() = _, %v; want nil", err)
			}
			test.TestDiff(t, "PSSingleQuote.Unquote()", it.Input, unquoted)
		})
	}
}

func TestPSDoubleQuote_Quote_Unquote(t *testing.T) {
	tests := []struct {
		Name, Input, Output, PwshOutput string
	}{
		{
			Name:   "empty string",
			Input:  "",
			Output: `""`,
		},
		{
			Name:       "ANSI escaping",
			Input:      "\a\b\x1B\f\n\r\t\v",
			Output:     "\"`a`b\x1B`f`n`r`t`v\"",
			PwshOutput: "\"`a`b`e`f`n`r`t`v\"",
		},
		{
			Name:   "special char escaping",
			Input:  "\"$`",
			Output: "\"`\"`$``\"",
		},
		{
			Name:       "`u byte escaping",
			Input:      "\x00\x01\x02\x03\x04\x05\x06\x0E\x0F\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1C\x1D\x1E\x1F",
			Output:     "\"`0\x01\x02\x03\x04\x05\x06\x0E\x0F\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1C\x1D\x1E\x1F\"",
			PwshOutput: "\"`0`u{01}`u{02}`u{03}`u{04}`u{05}`u{06}`u{0E}`u{0F}`u{10}`u{11}`u{12}`u{13}`u{14}`u{15}`u{16}`u{17}`u{18}`u{19}`u{1A}`u{1C}`u{1D}`u{1E}`u{1F}\"",
		},
		{
			Name:       "`u 4 char escaping",
			Input:      "\u0378\u0379\u0380\u0381\u0382",
			Output:     "\"\u0378\u0379\u0380\u0381\u0382\"",
			PwshOutput: "\"`u{0378}`u{0379}`u{0380}`u{0381}`u{0382}\"",
		},
		{
			Name:       "`u 6 char escaping",
			Input:      "\U0001000C\U00010027\U0001003B\U0001003E\U0001004E",
			Output:     "\"\U0001000C\U00010027\U0001003B\U0001003E\U0001004E\"",
			PwshOutput: "\"`u{01000C}`u{010027}`u{01003B}`u{01003E}`u{01004E}\"",
		},
	}
	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			quoted := PSDoubleQuote.Quote(td.Input)
			test.TestDiff(t, "PSDoubleQuote.Quote() ", td.Output, quoted)
			unquoted, err := PSDoubleQuote.Unquote(quoted)
			if err != nil {
				t.Fatalf("PSDoubleQuote.Unquote() = _, %v; want nil", err)
			}
			test.TestDiff(t, "PSDoubleQuote.Unquote()", td.Input, unquoted)
		})
		t.Run(td.Name+";pwsh.exe", func(t *testing.T) {
			quoted := PwshDoubleQuote.Quote(td.Input)
			expected := td.PwshOutput
			if expected == "" {
				expected = td.Output
			}
			test.TestDiff(t, "PwshDoubleQuote.Quote() ", expected, quoted)
			unquoted, err := PwshDoubleQuote.Unquote(quoted)
			if err != nil {
				t.Fatalf("PwshDoubleQuote.Unquote() = _, %v; want nil", err)
			}
			test.TestDiff(t, "PwshDoubleQuote.Unquote()", td.Input, unquoted)
		})
	}
}

func TestPSDoubleQuote_Unquote(t *testing.T) {
	tests := []struct {
		Name, Input, Output string
	}{
		{
			Name:   "multiple strings",
			Input:  `"a""""b"""`,
			Output: "ab",
		},
		{
			Name:   "ANSI escaping",
			Input:  "\"`a`b`e`f`n`r`t`v\"",
			Output: "\a\b\x1B\f\n\r\t\v",
		},
		{
			Name:   "short escaping",
			Input:  "\"`0`u{0}`u{378}`u{1000C}\"",
			Output: "\x00\x00\u0378\U0001000C",
		},
		{
			Name:   "unnecessary escaping",
			Input:  `"\p\z"`,
			Output: `\p\z`,
		},
	}
	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			unquoted, err := PSDoubleQuote.Unquote(td.Input)
			if err != nil {
				t.Fatalf("PSDoubleQuote.Unquote() = _, %v; want nil", err)
			}
			test.TestDiff(t, "PSDoubleQuote.Unquote()", td.Output, unquoted)
		})
		t.Run(td.Name+";pwsh.exe", func(t *testing.T) {
			unquoted, err := PwshDoubleQuote.Unquote(td.Input)
			if err != nil {
				t.Fatalf("PwshDoubleQuote.Unquote() = _, %v; want nil", err)
			}
			test.TestDiff(t, "PwshDoubleQuote.Unquote()", td.Output, unquoted)
		})
	}
}

func TestPSDoubleQuote_Unquote_ShouldFail(t *testing.T) {
	tests := []struct {
		Name, Input string
		Err         error
	}{
		{
			Name:  "unterminated string",
			Input: `"a`,
			Err: &quote.SyntaxError{
				Msg:    "unterminated quoted string",
				Offset: 2,
			},
		},
		{
			Name:  "unterminated escape sequence",
			Input: "\"`",
			Err: &quote.SyntaxError{
				Msg:    "unterminated escape sequence",
				Offset: 2,
			},
		},
		{
			Name:  "char outside of string",
			Input: "a",
			Err: &quote.SyntaxError{
				Msg:    "character U+0061 'a' outside of quoted string",
				Offset: 1,
			},
		},
		{
			Name:  "char after string",
			Input: `"a"a`,
			Err: &quote.SyntaxError{
				Msg:    "character U+0061 'a' outside of quoted string",
				Offset: 4,
			},
		},
		{
			Name:  `unescaped $`,
			Input: `"$"`,
			Err: &quote.SyntaxError{
				Msg:    "unescaped special character U+0024 '$'",
				Offset: 2,
			},
		},
		{
			Name:  "unterminated `u #1",
			Input: "\"`u",
			Err: &quote.SyntaxError{
				Msg:    "unterminated escape sequence `u",
				Offset: 3,
			},
		},
		{
			Name:  "unterminated `u #2",
			Input: "\"`u{",
			Err: &quote.SyntaxError{
				Msg:    "invalid escape sequence '`u'",
				Offset: 4,
			},
		},
		{
			Name:  "unterminated `u #3",
			Input: "\"`u{f",
			Err: &quote.SyntaxError{
				Msg:    "unterminated escape sequence `u",
				Offset: 5,
			},
		},
		{
			Name:  "invalid `u #1",
			Input: "\"`u{ffffff}\"",
			Err: &quote.SyntaxError{
				Msg:    "invalid escape sequence '`u{ffffff}'",
				Offset: 11,
			},
		},
		{
			Name:  "invalid `u #2",
			Input: "\"`u[ffffff}\"",
			Err: &quote.SyntaxError{
				Msg:    "invalid character U+005B '[' in escape sequence '`u'",
				Offset: 4,
			},
		},
		{
			Name:  "invalid `u #3",
			Input: "\"`u{ffffff]\"",
			Err: &quote.SyntaxError{
				Msg:    "invalid character U+005D ']' in escape sequence '`u'",
				Offset: 11,
			},
		},
	}
	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			_, err := PSDoubleQuote.Unquote(td.Input)
			if diff := cmp.Diff(td.Err, err); diff != "" {
				t.Errorf("PSDoubleQuote.Unquote() mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run(td.Name+";pwsh.exe", func(t *testing.T) {
			_, err := PwshDoubleQuote.Unquote(td.Input)
			if diff := cmp.Diff(td.Err, err); diff != "" {
				t.Errorf("PwshDoubleQuote.Unquote() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPSDoubleQuote_Quote_Unquote_InputTests(t *testing.T) {
	for _, it := range test.InputTests('\'', '$', '`') {
		t.Run(it.Name, func(t *testing.T) {
			if strings.HasPrefix(it.Name, "bytes:") {
				t.Skipf("it.Name=%s", it.Name)
			}
			quoted := PSDoubleQuote.Quote(it.Input)
			unquoted, err := PSDoubleQuote.Unquote(quoted)
			if err != nil {
				t.Fatalf("PSDoubleQuote.Unquote() = _, %v; want nil", err)
			}
			test.TestDiff(t, "PSDoubleQuote.Unquote()", it.Input, unquoted)
		})
	}
}

func TestPwshDoubleQuote_Quote_Unquote_InputTests(t *testing.T) {
	for _, it := range test.InputTests('\'', '$', '`') {
		t.Run(it.Name, func(t *testing.T) {
			if strings.HasPrefix(it.Name, "bytes:") {
				t.Skipf("it.Name=%s", it.Name)
			}
			quoted := PwshDoubleQuote.Quote(it.Input)
			unquoted, err := PwshDoubleQuote.Unquote(quoted)
			if err != nil {
				t.Fatalf("PwshDoubleQuote.Unquote() = _, %v; want nil", err)
			}
			test.TestDiff(t, "PwshDoubleQuote.Unquote()", it.Input, unquoted)
		})
	}
}
