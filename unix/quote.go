package unix

import "regexp"

var reUnsafeChars = regexp.MustCompile("[\\x00-\\x24&'()*;<=>?\\[\\]^`\\x7B-\\x7F\\x{00A0}]")

type unixQuote struct{}

func (unixQuote) MustQuote(s string) bool {
	return reUnsafeChars.MatchString(s)
}
