package windows

import (
	"strings"

	"github.com/sergeymakinen/go-quote"
)

var (
	cmdUnsafeChars   = "!\"&'+,;<=>[]^`{}~"
	cmdQuoteReplacer = strings.NewReplacer(
		"\t", "^\t",
		" ", "^ ",
		"!", "^!",
		`"`, `^"`,
		"&", "^&",
		"'", "^'",
		"+", "^+",
		",", "^,",
		";", "^;",
		"<", "^<",
		"=", "^=",
		">", "^>",
		"[", "^[",
		"]", "^]",
		"^", "^^",
		"`", "^`",
		"{", "^{",
		"}", "^}",
		"~", "^~",
	)
	cmdUnquoteReplacer = strings.NewReplacer(
		"^\t", "\t",
		"^ ", " ",
		"^!", "!",
		`^"`, `"`,
		"^&", "&",
		"^'", "'",
		"^+", "+",
		"^,", ",",
		"^;", ";",
		"^<", "<",
		"^=", "=",
		"^>", ">",
		"^[", "[",
		"^]", "]",
		"^^", "^",
		"^`", "`",
		"^{", "{",
		"^}", "}",
		"^~", "~",
	)
)

type cmd struct{}

func (cmd) MustQuote(s string) bool {
	return strings.ContainsAny(s, cmdUnsafeChars)
}

func (cmd) Quote(s string) string {
	return cmdQuoteReplacer.Replace(s)
}

func (cmd) Unquote(s string) (string, error) {
	return cmdUnquoteReplacer.Replace(s), nil
}

// Cmd quotes and unquotes strings containing characters special to the Windows command interpreter (cmd.exe).
//
// For example, the following string:
//
//  a b:"c d" 'e''f'  "g\""
//
// Would be quoted as:
//
//  a b:^"c d^" ^'e^'^'f^'  ^"g\^"^"
//
// See https://docs.microsoft.com/en-us/archive/blogs/twistylittlepassagesallalike/everyone-quotes-command-line-arguments-the-wrong-way
// for details.
var Cmd quote.Quoting = cmd{}
