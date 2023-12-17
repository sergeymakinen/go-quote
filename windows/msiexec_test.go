package windows

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sergeymakinen/go-quote"
	"github.com/sergeymakinen/go-quote/internal/testutil"
)

func TestMsiexec_Quote(t *testing.T) {
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
			Input:  `a"b`,
			Output: `"a""b"`,
		},
	}
	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			quoted := Msiexec.Quote(td.Input)
			testutil.TestDiff(t, "Msiexec.Quote() ", td.Output, quoted)
			unquoted, err := Msiexec.Unquote(quoted)
			if err != nil {
				t.Fatalf("Msiexec.Unquote() = _, %v; want nil", err)
			}
			testutil.TestDiff(t, "Msiexec.Unquote()", td.Input, unquoted)
		})
	}
}

func TestMsiexec_Unquote(t *testing.T) {
	unquoted, err := Msiexec.Unquote(`"ab""""""cd""ef"`)
	if err != nil {
		t.Fatalf("Msiexec.Unquote() = _, %v; want nil", err)
	}
	testutil.TestDiff(t, "Msiexec.Unquote()", `ab"""cd"ef`, unquoted)
}

func TestMsiexec_Unquote_ShouldFail(t *testing.T) {
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
			Input: `"a""`,
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
			Input: `"a"a`,
			Err: &quote.SyntaxError{
				Msg:    "character U+0061 'a' outside of quoted string",
				Offset: 4,
			},
		},
	}
	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			_, err := Msiexec.Unquote(td.Input)
			if diff := cmp.Diff(td.Err, err); diff != "" {
				t.Errorf("Msiexec.Unquote() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestMsiexec_Quote_Unquote_InputTests(t *testing.T) {
	for _, it := range testutil.InputTests('"') {
		t.Run(it.Name, func(t *testing.T) {
			quoted := Msiexec.Quote(it.Input)
			unquoted, err := Msiexec.Unquote(quoted)
			if err != nil {
				t.Fatalf("Msiexec.Unquote() = _, %v; want nil", err)
			}
			testutil.TestDiff(t, "Msiexec.Unquote()", it.Input, unquoted)
		})
	}
}
