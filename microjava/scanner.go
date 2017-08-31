package microjava

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

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
		Line:   s.line,
		Column: s.column,
	}
	if s.currChar == '\u0080' {
		token.Kind = tcEOF
	} else if ascii.IsDigit(s.currChar) {
		s.readNumber(token)
	} else if ascii.IsLetter(s.currChar) {
		s.readName(token)
	} else if s.currChar == '\'' {
		s.readCharacter(token)
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
			token.Kind = tcSlash
			token.Data = "/"
		}
	} else if !s.readOperator(token) {
		if token.Data == "" {
			token.Kind = tcNone
			token.ErrorMsg = "invalid symbol"
		}
	}
	return token
}

// NextChar advances the scanner to the next character in the stream
func (s *Scanner) nextChar() {

	c, err := s.reader.ReadByte()
	if err == io.EOF {
		c = '\u0080'
	} else if err != nil {
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

	token.Kind = tcNumber
	token.Data = lexeme.String()

	value, err := strconv.ParseInt(token.Data, 10, 32)
	if err != nil {
		panic(err)
	}
	token.Value = int(value)

}

func (s *Scanner) readName(token *Token) {
	var lexeme bytes.Buffer

	for ascii.IsLetter(s.currChar) {
		lexeme.WriteByte(s.currChar)
		s.nextChar()
	}

	token.Data = lexeme.String()
	token.Kind = GetKeywordKind(token.Data)
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
	token.ErrorMsg = "empty character token"
	token.Kind = tcNone
	token.Data = "''"
	token.Column = s.column
	s.nextChar()
}

func (s *Scanner) readEscapedChar(token *Token, lexeme *bytes.Buffer) {
	s.nextChar()
	lexeme.WriteByte(s.currChar)
	if s.currChar == 'n' || s.currChar == 't' || s.currChar == 'r' {
		s.readCloseChar(token, lexeme)
	} else {
		token.ErrorMsg = "invalid character escape sequence"
		token.Kind = tcNone
		token.Column = s.column
		s.skipUntilCloseChar(lexeme)
		token.Data = lexeme.String()
		s.nextChar()
	}
}

func (s *Scanner) readCloseChar(token *Token, lexeme *bytes.Buffer) {
	s.nextChar()
	lexeme.WriteByte(s.currChar)

	if s.currChar == '\'' {
		token.Data = lexeme.String()
		token.Kind = tcCharCon
		if token.Data[1] == '\\' {
			charVal := token.Data[1]
			if charVal == 'n' {
				token.CharValue = '\n'
			} else if charVal == 'r' {
				token.CharValue = '\r'
			} else if charVal == 't' {
				token.CharValue = '\t'
			}
		} else {
			token.CharValue = token.Data[1]
		}
		s.nextChar()
	} else {
		token.ErrorMsg = "unclosed char constant"
		token.Kind = tcNone
		token.Column = s.column
		s.skipUntilCloseChar(lexeme)
		token.Data = lexeme.String()
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
			token.Kind = tcEql
			s.nextChar()
		} else {
			token.Kind = tcAssign
		}
	} else if s.currChar == '!' {
		s.nextChar()
		if s.currChar == '=' {
			token.Kind = tcNeq
		} else {
			token.Kind = tcNone
			token.ErrorMsg = "invalid operator"
			token.Data = fmt.Sprintf("!%v", s.currChar)
			s.nextChar()
			return false
		}
	} else if s.currChar == '<' {
		s.nextChar()
		if s.currChar == '=' {
			token.Kind = tcLeq
			s.nextChar()
		} else {
			token.Kind = tcLss
		}
	} else if s.currChar == '>' {
		s.nextChar()
		if s.currChar == '=' {
			token.Kind = tcGeq
			s.nextChar()
		} else {
			token.Kind = tcGtr
		}
	} else if strings.IndexByte("+-*%;,.()[]{}", s.currChar) != -1 {
		token.Kind = GetOperatorKind(string(s.currChar))
		s.nextChar()
	} else {
		return false
	}
	token.Data = GetTokenName(token.Kind)
	return true
}
