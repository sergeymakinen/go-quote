package unix

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sergeymakinen/go-quote"
	"github.com/sergeymakinen/go-quote/internal/test"
)

func TestANSIC_Quote_QuoteBinary_Unquote_UnquoteBinary(t *testing.T) {
	tests := []struct {
		Name, Input, Output, OutputBinary string
	}{
		{
			Name:   "empty string",
			Input:  "",
			Output: "$''",
		},
		{
			Name:   "ANSI escaping",
			Input:  "\a\b\x1B\f\n\r\t\v",
			Output: `$'\a\b\e\f\n\r\t\v'`,
		},
		{
			Name:   "special char escaping",
			Input:  `"'?\`,
			Output: `$'\"\'\?\\'`,
		},
		{
			Name:   `\x escaping`,
			Input:  "\x00\x01\x02\x03\x04\x05\x06\x0E\x0F\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1C\x1D\x1E\x1F",
			Output: `$'\x00\x01\x02\x03\x04\x05\x06\x0E\x0F\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1C\x1D\x1E\x1F'`,
		},
		{
			Name:         `\u escaping`,
			Input:        "\u0378\u0379\u0380\u0381\u0382",
			Output:       `$'\u0378\u0379\u0380\u0381\u0382'`,
			OutputBinary: `$'\xCD\xB8\xCD\xB9\xCE\x80\xCE\x81\xCE\x82'`,
		},
		{
			Name:         `\U escaping`,
			Input:        "\U0001000C\U00010027\U0001003B\U0001003E\U0001004E",
			Output:       `$'\U0001000C\U00010027\U0001003B\U0001003E\U0001004E'`,
			OutputBinary: `$'\xF0\x90\x80\x8C\xF0\x90\x80\xA7\xF0\x90\x80\xBB\xF0\x90\x80\xBE\xF0\x90\x81\x8E'`,
		},
	}
	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			quoted := ANSIC.Quote(td.Input)
			test.TestDiff(t, "ANSIC.Quote() ", td.Output, quoted)
			unquoted, err := ANSIC.Unquote(quoted)
			if err != nil {
				t.Fatalf("ANSIC.Unquote() = _, %v; want nil", err)
			}
			test.TestDiff(t, "ANSIC.Unquote()", td.Input, unquoted)
		})
		t.Run(td.Name+";binary", func(t *testing.T) {
			quoted := ANSIC.QuoteBinary([]byte(td.Input))
			expected := td.OutputBinary
			if expected == "" {
				expected = td.Output
			}
			test.TestDiff(t, "ANSIC.QuoteBinary() ", expected, quoted)
			unquoted, err := ANSIC.UnquoteBinary(quoted)
			if err != nil {
				t.Fatalf("ANSIC.UnquoteBinary() = _, %v; want nil", err)
			}
			test.TestDiff(t, "ANSIC.UnquoteBinary()", td.Input, string(unquoted))
		})
	}
}

func TestANSIC_Unquote_UnquoteBinary(t *testing.T) {
	tests := []struct {
		Name, Input, Output string
	}{
		{
			Name:   "multiple strings",
			Input:  "$'a'$''$'b'$''",
			Output: "ab",
		},
		{
			Name:   "ANSI escaping",
			Input:  `$'\a\b\e\E\f\n\r\t\v'`,
			Output: "\a\b\x1B\x1B\f\n\r\t\v",
		},
		{
			Name:   `\c escaping`,
			Input:  `$'\c?\c@\cA\cB\cC\cD\cE\cF\cG\cH\cI\cJ\cK\cL\cM\cN\cO\cP\cQ\cR\cS\cT\cU\cV\cW\cX\cY\cZ\c[\c\\c]\c^\c_\ca\cb\cc\cd\ce\cf\cg\ch\ci\cj\ck\cl\cm\cn\co\cp\cq\cr\cs\ct\cu\cv\cw\cx\cy\cz'`,
			Output: "\x7F\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1B\x1C\x1D\x1E\x1F\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A",
		},
		{
			Name:   "short escaping",
			Input:  `$'\0\x0\u378\U1000C'`,
			Output: "\x00\x00\u0378\U0001000C",
		},
		{
			Name:   "unnecessary escaping",
			Input:  `$'\p\z'`,
			Output: `\p\z`,
		},
	}
	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			unquoted, err := ANSIC.Unquote(td.Input)
			if err != nil {
				t.Fatalf("ANSIC.Unquote() = _, %v; want nil", err)
			}
			test.TestDiff(t, "ANSIC.Unquote()", td.Output, unquoted)
		})
		t.Run(td.Name+";binary", func(t *testing.T) {
			unquoted, err := ANSIC.UnquoteBinary(td.Input)
			if err != nil {
				t.Fatalf("ANSIC.UnquoteBinary() = _, %v; want nil", err)
			}
			test.TestDiff(t, "ANSIC.UnquoteBinary()", td.Output, string(unquoted))
		})
	}
}

func TestANSIC_Unquote_UnquoteBinary_ShouldFail(t *testing.T) {
	tests := []struct {
		Name, Input string
		Err         error
	}{
		{
			Name:  "unterminated quoting start",
			Input: "$",
			Err: &quote.SyntaxError{
				Msg:    "character U+0024 '$' outside of quoted string",
				Offset: 1,
			},
		},
		{
			Name:  "unterminated string",
			Input: `$'a`,
			Err: &quote.SyntaxError{
				Msg:    "unterminated quoted string",
				Offset: 3,
			},
		},
		{
			Name:  "unterminated escape sequence",
			Input: `$'\`,
			Err: &quote.SyntaxError{
				Msg:    "unterminated escape sequence",
				Offset: 3,
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
			Input: "$'a'a",
			Err: &quote.SyntaxError{
				Msg:    "character U+0061 'a' outside of quoted string",
				Offset: 5,
			},
		},
		{
			Name:  `unterminated \c`,
			Input: `$'\c`,
			Err: &quote.SyntaxError{
				Msg:    "unterminated escape sequence `\\c`",
				Offset: 4,
			},
		},
		{
			Name:  `invalid \c`,
			Input: `$'\c+a`,
			Err: &quote.SyntaxError{
				Msg:    "invalid character U+002B '+' in escape sequence `\\c`",
				Offset: 5,
			},
		},
		{
			Name:  `unterminated \x`,
			Input: `$'\x`,
			Err: &quote.SyntaxError{
				Msg:    "unterminated escape sequence `\\x`",
				Offset: 4,
			},
		},
		{
			Name:  `invalid \U`,
			Input: `$'\Uffffffff `,
			Err: &quote.SyntaxError{
				Msg:    "invalid escape sequence `\\Uffffffff`",
				Offset: 12,
			},
		},
		{
			Name:  "invalid octal escape",
			Input: `$'\777`,
			Err: &quote.SyntaxError{
				Msg:    "invalid escape sequence `\\777`",
				Offset: 6,
			},
		},
	}
	for _, td := range tests {
		t.Run(td.Name, func(t *testing.T) {
			_, err := ANSIC.Unquote(td.Input)
			if diff := cmp.Diff(td.Err, err); diff != "" {
				t.Errorf("ANSIC.Unquote() mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run(td.Name+";binary", func(t *testing.T) {
			_, err := ANSIC.UnquoteBinary(td.Input)
			if diff := cmp.Diff(td.Err, err); diff != "" {
				t.Errorf("ANSIC.UnquoteBinary() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestANSIC_Quote_Unquote_InputTests(t *testing.T) {
	for _, it := range test.InputTests('\'', '\t', '\n', ' ', '$', '"') {
		t.Run(it.Name, func(t *testing.T) {
			if strings.HasPrefix(it.Name, "bytes:") {
				t.Skipf("it.Name=%s", it.Name)
			}
			quoted := ANSIC.Quote(it.Input)
			unquoted, err := ANSIC.Unquote(quoted)
			if err != nil {
				t.Fatalf("ANSIC.Unquote() = _, %v; want nil", err)
			}
			test.TestDiff(t, "ANSIC.Unquote()", it.Input, unquoted)
		})
	}
}

func TestANSIC_QuoteBinary_UnquoteBinary_InputTests(t *testing.T) {
	for _, it := range test.InputTests('\'', '\t', '\n', ' ', '$', '"') {
		t.Run(it.Name, func(t *testing.T) {
			quoted := ANSIC.QuoteBinary([]byte(it.Input))
			unquoted, err := ANSIC.UnquoteBinary(quoted)
			if err != nil {
				t.Fatalf("ANSIC.UnquoteBinary() = _, %v; want nil", err)
			}
			test.TestDiff(t, "ANSIC.UnquoteBinary()", it.Input, string(unquoted))
		})
	}
}
