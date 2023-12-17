package unix

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sergeymakinen/go-quote"
	"github.com/sergeymakinen/go-quote/internal/testutil"
)

func Example() {
	filename := `Long File With 'Single' & "Double" Quotes.txt`
	fmt.Println(SingleQuote.MustQuote(filename))
	// Echoing inside 'sh -c' requires a quoting to make it safe with an arbitrary string
	quoted := SingleQuote.Quote(filename)
	fmt.Println([]string{
		"sh",
		"-c",
		fmt.Sprintf("echo %s | callme", quoted),
	})
	unquoted, _ := SingleQuote.Unquote(quoted)
	fmt.Println(unquoted)
	// Output:
	// true
	// [sh -c echo 'Long File With '"'"'Single'"'"' & "Double" Quotes.txt' | callme]
	// Long File With 'Single' & "Double" Quotes.txt
}

func TestSingleQuote_Quote(t *testing.T) {
	tests := []struct {
		Name, Input, Output string
	}{
		{
			Name:   "empty string",
			Input:  "",
			Output: "''",
		},
		{
			Name:   "special char escaping #1",
			Input:  `"'`,
			Output: `'"'"'"''`,
		},
		{
			Name:   "special char escaping #2",
			Input:  `a"'b`,
			Output: `'a"'"'"'b'`,
		},
	}
	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			quoted := SingleQuote.Quote(td.Input)
			testutil.TestDiff(t, "SingleQuote.Quote() ", td.Output, quoted)
			unquoted, err := SingleQuote.Unquote(quoted)
			if err != nil {
				t.Fatalf("SingleQuote.Unquote() = _, %v; want nil", err)
			}
			testutil.TestDiff(t, "SingleQuote.Unquote()", td.Input, unquoted)
		})
	}
}

func TestSingleQuote_Unquote(t *testing.T) {
	tests := []struct {
		Name, Input, Output string
	}{
		{
			Name:   "multiple strings",
			Input:  `'ab""'"'"'cd''ef'`,
			Output: `ab""'cdef`,
		},
		{
			Name:   "multiple single quotes",
			Input:  `'ab'"'''"`,
			Output: "ab'''",
		},
		{
			Name:   "unnecessary escaping",
			Input:  `'\p\z'`,
			Output: `\p\z`,
		},
	}
	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			unquoted, err := SingleQuote.Unquote(td.Input)
			if err != nil {
				t.Fatalf("SingleQuote.Unquote() = _, %v; want nil", err)
			}
			testutil.TestDiff(t, "SingleQuote.Unquote()", td.Output, unquoted)
		})
	}
}

func TestSingleQuote_Unquote_ShouldFail(t *testing.T) {
	tests := []struct {
		Name, Input string
		Err         error
	}{
		{
			Name:  "unterminated string",
			Input: `'a`,
			Err: &quote.SyntaxError{
				Msg:    "unterminated quoted string",
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
			Input: "'a'a",
			Err: &quote.SyntaxError{
				Msg:    "character U+0061 'a' outside of quoted string",
				Offset: 4,
			},
		},
		{
			Name:  "not single quote in double quotes",
			Input: `'a'"b"`,
			Err: &quote.SyntaxError{
				Msg:    "unsupported character U+0062 'b' in double quoted string",
				Offset: 5,
			},
		},
	}
	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			_, err := SingleQuote.Unquote(td.Input)
			if diff := cmp.Diff(td.Err, err); diff != "" {
				t.Errorf("SingleQuote.Unquote() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSingleQuote_Quote_Unquote_InputTests(t *testing.T) {
	for _, it := range testutil.InputTests('\'', '\t', '\n', ' ', '"') {
		t.Run(it.Name, func(t *testing.T) {
			quoted := SingleQuote.Quote(it.Input)
			unquoted, err := SingleQuote.Unquote(quoted)
			if err != nil {
				t.Fatalf("SingleQuote.Unquote() = _, %v; want nil", err)
			}
			testutil.TestDiff(t, "SingleQuote.Unquote()", it.Input, unquoted)
		})
	}
}

func TestDoubleQuote_Quote_Unquote(t *testing.T) {
	tests := []struct {
		Name, Input, Output string
	}{
		{
			Name:   "empty string",
			Input:  "",
			Output: `""`,
		},
		{
			Name:   "special char escaping",
			Input:  "!\"$\\`",
			Output: "\"\\!\\\"\\$\\\\\\`\"",
		},
	}
	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			quoted := DoubleQuote.Quote(td.Input)
			testutil.TestDiff(t, "DoubleQuote.Quote() ", td.Output, quoted)
			unquoted, err := DoubleQuote.Unquote(quoted)
			if err != nil {
				t.Fatalf("DoubleQuote.Unquote() = _, %v; want nil", err)
			}
			testutil.TestDiff(t, "DoubleQuote.Unquote()", td.Input, unquoted)
		})
	}
}

func TestDoubleQuote_Unquote(t *testing.T) {
	tests := []struct {
		Name, Input, Output string
	}{
		{
			Name:   "multiple strings",
			Input:  `"a""""""b"""`,
			Output: "ab",
		},
		{
			Name:   "unnecessary escaping",
			Input:  `"\p\z"`,
			Output: `\p\z`,
		},
	}
	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			unquoted, err := DoubleQuote.Unquote(td.Input)
			if err != nil {
				t.Fatalf("DoubleQuote.Unquote() = _, %v; want nil", err)
			}
			testutil.TestDiff(t, "DoubleQuote.Unquote()", td.Output, unquoted)
		})
	}
}

func TestDoubleQuote_Unquote_ShouldFail(t *testing.T) {
	tests := []struct {
		Name, Input string
		Err         error
	}{
		{
			Name:  "unterminated string #1",
			Input: `"a`,
			Err: &quote.SyntaxError{
				Msg:    "unterminated quoted string",
				Offset: 2,
			},
		},
		{
			Name:  "unterminated string #2",
			Input: `"""`,
			Err: &quote.SyntaxError{
				Msg:    "unterminated quoted string",
				Offset: 3,
			},
		},
		{
			Name:  "unterminated escape sequence",
			Input: `"\`,
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
			Name:  "unescaped !",
			Input: `"!"`,
			Err: &quote.SyntaxError{
				Msg:    "unescaped special character U+0021 '!'",
				Offset: 2,
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
			Name:  "unescaped `",
			Input: "\"`\"",
			Err: &quote.SyntaxError{
				Msg:    "unescaped special character U+0060 '`'",
				Offset: 2,
			},
		},
	}
	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			_, err := DoubleQuote.Unquote(td.Input)
			if diff := cmp.Diff(td.Err, err); diff != "" {
				t.Errorf("DoubleQuote.Unquote() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDoubleQuote_Quote_Unquote_InputTests(t *testing.T) {
	for _, it := range testutil.InputTests('"', '\t', '\n', ' ', '$', '\'') {
		t.Run(it.Name, func(t *testing.T) {
			quoted := DoubleQuote.Quote(it.Input)
			unquoted, err := DoubleQuote.Unquote(quoted)
			if err != nil {
				t.Fatalf("DoubleQuote.Unquote() = _, %v; want nil", err)
			}
			testutil.TestDiff(t, "DoubleQuote.Unquote()", it.Input, unquoted)
		})
	}
}
