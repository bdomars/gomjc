package microjava

import "io"

// Scanner tokenizes input
type Scanner struct {
	infile   io.Reader
	line     int
	column   int
	position int
	currChar byte
}

// NewScanner initializes a new scanner for a reader
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		infile: r,
		line:   1,
	}
}

// NextToken returns the next scanned token of the input stream
func (s *Scanner) NextToken() *Token {
	t := new(Token)
	t.kind = tcPlus
	return t
}
