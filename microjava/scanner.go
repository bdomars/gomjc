package microjava

import "io"

// Scanner tokenizes input
type Scanner struct {
	infile *io.Reader
}

// NewScanner initializes a new scanner for a file
func NewScanner(r *io.Reader) *Scanner {
	return &Scanner{
		infile: r,
	}
}
