package microjava

// TokenCode is a type for all token classes used in the scanner
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
	kind   TokenCode
	line   int
	column int
	value  int
	data   string
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
