package microjava

import (
	"bytes"
	"fmt"
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
		s.nextChar()
		if s.currChar == '/' {
			for {
				s.nextChar() // skipping comments
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
func (s *Scanner) nextChar() {

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
		s.nextChar()
	}
}

func (s *Scanner) readNumber(token *Token) {
	var lexeme bytes.Buffer

	for ascii.IsDigit(s.currChar) {
		lexeme.WriteByte(s.currChar)
		s.nextChar()
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
		s.nextChar()
	}

	token.data = lexeme.String()
	token.kind = GetKeywordKind(token.data)
}

func (s *Scanner) readCharacter(token *Token) {
	var lexeme bytes.Buffer
	lexeme.WriteByte(s.currChar)
	s.nextChar()
	lexeme.WriteByte(s.currChar)
	if s.currChar == '\'' {
		s.errorEmptyChar(token)
	} else if s.currChar == '\\' {
		s.readEscapedChar(token, &lexeme)
	} else {
		s.readCloseChar(token, &lexeme)
	}
}

func (s *Scanner) errorEmptyChar(token *Token) {
	token.errorMsg = "empty character token"
	token.kind = tcNone
	token.data = "''"
	token.column = s.column
	s.nextChar()
}

func (s *Scanner) readEscapedChar(token *Token, lexeme *bytes.Buffer) {
	s.nextChar()
	lexeme.WriteByte(s.currChar)
	if s.currChar == 'n' || s.currChar == 't' || s.currChar == 'r' {
		s.readCloseChar(token, lexeme)
	} else {
		token.errorMsg = "invalid character escape sequence"
		token.kind = tcNone
		token.column = s.column
		s.skipUntilCloseChar(lexeme)
		token.data = lexeme.String()
		s.nextChar()
	}
}

func (s *Scanner) readCloseChar(token *Token, lexeme *bytes.Buffer) {
	s.nextChar()
	lexeme.WriteByte(s.currChar)

	if s.currChar == '\'' {
		token.data = lexeme.String()
		token.kind = tcCharCon
		if token.data[1] == '\\' {
			charVal := token.data[1]
			if charVal == 'n' {
				token.charValue = '\n'
			} else if charVal == 'r' {
				token.charValue = '\r'
			} else if charVal == 't' {
				token.charValue = '\t'
			}
		} else {
			token.charValue = token.data[1]
		}
		s.nextChar()
	} else {
		token.errorMsg = "unclosed char constant"
		token.kind = tcNone
		token.column = s.column
		s.skipUntilCloseChar(lexeme)
		token.data = lexeme.String()
		s.nextChar()
	}
}

func (s *Scanner) skipUntilCloseChar(lexeme *bytes.Buffer) {
	for {
		s.nextChar()
		lexeme.WriteByte(s.currChar)

		if s.currChar == '\'' || s.currChar == '\u0080' {
			break
		}
	}
}

func (s *Scanner) readOperator(token *Token) bool {
	if s.currChar == '=' {
		s.nextChar()
		if s.currChar == '=' {
			token.kind = tcEql
			s.nextChar()
		} else {
			token.kind = tcAssign
		}
	} else if s.currChar == '!' {
		s.nextChar()
		if s.currChar == '=' {
			token.kind = tcNeq
		} else {
			token.kind = tcNone
			token.errorMsg = "invalid operator"
			token.data = fmt.Sprintf("!%v", s.currChar)
			s.nextChar()
			return false
		}
	} else {
		return false
	}
	token.data = "poop"
	return true
}
