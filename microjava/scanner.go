package microjava

import (
	"bytes"
	"io"
	"strconv"

	ascii "github.com/galsondor/go-ascii"
)

// Scanner tokenizes input
type Scanner struct {
	reader   io.ByteReader
	line     int
	column   int
	position int
	currChar byte
}

// NewScanner initializes a new scanner for a reader
func NewScanner(r io.ByteReader) *Scanner {
	return &Scanner{
		reader: r,
		line:   1,
	}
}

// NextToken returns the next scanned token of the input stream
func (s *Scanner) NextToken() *Token {
	s.skipNonPrint()
	token := &Token{
		line:   s.line,
		column: s.column,
	}
	if s.currChar == '\u0080' {
		token.kind = tcEOF
	} else if ascii.IsDigit(s.currChar) {
		s.readNumber(token)
	} else if ascii.IsLetter(s.currChar) {
		s.readName(token)
	} else if s.currChar == '\'' {
		//s.readCharacter(token)
	} else if s.currChar == '/' {
		s.NextChar()
		if s.currChar == '/' {
			for {
				s.NextChar() // skipping comments
				if s.currChar == '\n' || s.currChar == '\u0080' {
					break
				}
			}
		} else {
			token.kind = tcSlash
			token.data = "/"
		}
	}
	return token
}

// NextChar advances the scanner to the next character in the stream
func (s *Scanner) NextChar() {

	c, err := s.reader.ReadByte()
	if err != nil {
		panic(err)
	}

	s.currChar = c
	s.column++
	s.position++

	if s.currChar == '\n' {
		s.line++
		s.column = 0
	}
}

func (s *Scanner) skipNonPrint() {
	for s.currChar <= ' ' {
		s.NextChar()
	}
}

func (s *Scanner) readNumber(token *Token) {
	var lexeme bytes.Buffer

	for ascii.IsDigit(s.currChar) {
		lexeme.WriteByte(s.currChar)
		s.NextChar()
	}

	token.kind = tcNumber
	token.data = lexeme.String()

	value, err := strconv.ParseInt(token.data, 10, 32)
	if err != nil {
		panic(err)
	}
	token.value = int(value)

}

func (s *Scanner) readName(token *Token) {
	var lexeme bytes.Buffer

	for ascii.IsLetter(s.currChar) {
		lexeme.WriteByte(s.currChar)
		s.NextChar()
	}

	token.data = lexeme.String()
	token.kind = GetKeywordKind(token.data)
}
