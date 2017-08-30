package microjava

import (
	"fmt"
)

// TokenCode is a type for all token kinds used in the scanner
type TokenCode int

const (
	tcNone TokenCode = iota
	tcIdent
	tcNumber
	tcCharCon
	tcPlus
	tcMinus
	tcTimes
	tcSlash
	tcRem
	tcEql
	tcNeq
	tcLss
	tcLeq
	tcGtr
	tcGeq
	tcAssign
	tcSemicolon
	tcComma
	tcPeriod
	tcLpar
	tcRpar
	tcLbrack
	tcRbrack
	tcLbrace
	tcRbrace
	tcClass
	tcElse
	tcFinal
	tcIf
	tcNew
	tcPrint
	tcProgram
	tcRead
	tcReturn
	tcVoid
	tcWhile
	tcEOF
)

// Token holds the representation of one token
type Token struct {
	Kind      TokenCode
	Line      int
	Column    int
	Value     int
	Data      string
	ErrorMsg  string
	CharValue byte
}

// TokenNames will map a token code to a human readable name
var TokenNames = map[TokenCode]string{
	tcNone:      "None",
	tcIdent:     "Identifier",
	tcNumber:    "Number",
	tcCharCon:   "CharacterConstant",
	tcPlus:      "+",
	tcMinus:     "-",
	tcTimes:     "*",
	tcSlash:     "/",
	tcRem:       "%",
	tcEql:       "==",
	tcNeq:       "!=",
	tcLss:       "<",
	tcLeq:       "<=",
	tcGtr:       ">",
	tcGeq:       ">=",
	tcAssign:    "=",
	tcSemicolon: ";",
	tcComma:     ",",
	tcPeriod:    ".",
	tcLpar:      "(",
	tcRpar:      ")",
	tcLbrack:    "[",
	tcRbrack:    "]",
	tcLbrace:    "{",
	tcRbrace:    "}",
	tcClass:     "class",
	tcElse:      "else",
	tcFinal:     "final",
	tcIf:        "if",
	tcNew:       "new",
	tcPrint:     "print",
	tcProgram:   "program",
	tcRead:      "read",
	tcReturn:    "return",
	tcVoid:      "void",
	tcWhile:     "while",
	tcEOF:       "EOF",
}

func (t Token) String() string {
	if t.Kind == tcIdent || t.Kind == tcNumber {
		return fmt.Sprintf("<Token: kind='%v', data='%v'>", TokenNames[t.Kind], t.Data)
	}
	return fmt.Sprintf("<Token: kind='%v'>", TokenNames[t.Kind])

}

// KeywordMap defines a lookup table for the keywords used
var KeywordMap = map[string]TokenCode{
	"class":   tcClass,
	"else":    tcElse,
	"final":   tcFinal,
	"if":      tcIf,
	"new":     tcNew,
	"print":   tcPrint,
	"program": tcProgram,
	"read":    tcRead,
	"return":  tcReturn,
	"void":    tcVoid,
	"while":   tcWhile,
}

// GetKeywordKind returns a TokenCode for a keyword string, or defaults to the Identifier kind
func GetKeywordKind(lexeme string) TokenCode {
	if kind, ok := KeywordMap[lexeme]; ok {
		return kind
	}
	return tcIdent
}

// OperatorMap for looking up token code for an operator
var OperatorMap = map[string]TokenCode{
	"+": tcPlus,
	"-": tcMinus,
	"*": tcTimes,
	"%": tcRem,
	";": tcSemicolon,
	".": tcPeriod,
	"(": tcLpar,
	")": tcRpar,
	"[": tcLbrack,
	"]": tcRbrack,
	"{": tcLbrace,
	"}": tcRbrace,
}

// GetOperatorKind returns a TokenCode for an operator string
func GetOperatorKind(lexeme string) TokenCode {
	if kind, ok := OperatorMap[lexeme]; ok {
		return kind
	}
	return tcNone
}

// IsEOF returns true if this is the last token in stream
func (t *Token) IsEOF() bool {
	return t.Kind == tcEOF
}

// IsError returns true if this is an invalid token with an error message
func (t *Token) IsError() bool {
	return t.ErrorMsg != ""
}
