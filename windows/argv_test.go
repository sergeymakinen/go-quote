package windows

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sergeymakinen/go-quote"
	"github.com/sergeymakinen/go-quote/internal/testutil"
)

func TestArgv_Quote(t *testing.T) {
	tests := []struct {
		Name, Input, Output string
	}{
		{
			Name:   "empty string",
			Input:  "",
			Output: `""`,
		},
		{
			Name:   "special char escaping #1",
			Input:  `\"\`,
			Output: `"\\\"\\"`,
		},
		{
			Name:   "special char escaping #2",
			Input:  `a\\b\\\c\\""`,
			Output: `"a\\b\\\c\\\\\"\""`,
		},
	}
	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			quoted := Argv.Quote(td.Input)
			testutil.TestDiff(t, "Argv.Quote() ", td.Output, quoted)
			unquoted, err := Argv.Unquote(quoted)
			if err != nil {
				t.Fatalf("Argv.Unquote() = _, %v; want nil", err)
			}
			testutil.TestDiff(t, "Argv.Unquote()", td.Input, unquoted)
		})
	}
}

func TestArgv_Unquote(t *testing.T) {
	tests := []struct {
		Name, Input, Output string
	}{
		{
			Name:   "multiple strings",
			Input:  `"ab""\"""""cd""ef"`,
			Output: `ab"cdef`,
		},
		{
			Name:   "unnecessary escaping",
			Input:  `"\p\z"`,
			Output: `\p\z`,
		},
	}
	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			unquoted, err := Argv.Unquote(td.Input)
			if err != nil {
				t.Fatalf("Argv.Unquote() = _, %v; want nil", err)
			}
			testutil.TestDiff(t, "Argv.Unquote()", td.Output, unquoted)
		})
	}
}

func TestArgv_Unquote_ShouldFail(t *testing.T) {
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
			Input: `\`,
			Err: &quote.SyntaxError{
				Msg:    `character U+005C '\' outside of quoted string`,
				Offset: 1,
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
	}
	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			_, err := Argv.Unquote(td.Input)
			if diff := cmp.Diff(td.Err, err); diff != "" {
				t.Errorf("Argv.Unquote() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestArgv_Quote_Unquote_InputTests(t *testing.T) {
	for _, it := range testutil.InputTests('"') {
		t.Run(it.Name, func(t *testing.T) {
			quoted := Argv.Quote(it.Input)
			unquoted, err := Argv.Unquote(quoted)
			if err != nil {
				t.Fatalf("Argv.Unquote() = _, %v; want nil", err)
			}
			testutil.TestDiff(t, "Argv.Unquote()", it.Input, unquoted)
		})
	}
}
