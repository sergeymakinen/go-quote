// Package quote defines interfaces shared by other packages
// that quote command-line arguments and variables.
//
// See the documentation for the unix and windows packages for more information.
package quote

// Quoting quotes and and unquotes textual command-line arguments and variables.
type Quoting interface {
	// MustQuote reports whether s must be quoted in order
	// to appear correctly as a single command-line argument or variable.
	MustQuote(s string) bool

	// Quote returns s quoted such that it appears correctly
	// as a single command-line argument or variable.
	Quote(s string) string

	// Unquote interprets s as a quoted string, returning
	// the string value that s quotes.
	Unquote(s string) (string, error)
}

// BinaryQuoting quotes and and unquotes binary command-line arguments and variables.
type BinaryQuoting interface {
	Quoting

	// QuoteBinary returns b quoted such that it appears correctly
	// as a single command-line argument or variable.
	QuoteBinary(b []byte) string

	// Unquote interprets s as a quoted string, returning
	// the bytes that s quotes.
	UnquoteBinary(s string) ([]byte, error)
}

// SyntaxError represents an error during unquoting of the string.
type SyntaxError struct {
	Msg    string // description of error
	Offset int    // error occurred after reading Offset bytes
}

func (e *SyntaxError) Error() string { return "syntax error: " + e.Msg }
