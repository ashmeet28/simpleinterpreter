package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Token struct {
	T int
	S string
	L int
}

type Scanner struct {
	Source  []byte
	Current int
	Line    int
}

type Compiler struct {
	Source  []Token
	Current int
}

const (
	PREC_ILLEGAL int = iota

	PREC_NONE
	PREC_OR
	PREC_AND
	PREC_COMP
	PREC_ADD
	PREC_MUL
)

const (
	TOKEN_TYPE_ILLEGAL int = iota

	TOKEN_TYPE_LEFT_PAREN
	TOKEN_TYPE_RIGHT_PAREN
	TOKEN_TYPE_LEFT_BRACE
	TOKEN_TYPE_RIGHT_BRACE
	TOKEN_TYPE_COMMA
	TOKEN_TYPE_DOT
	TOKEN_TYPE_MINUS
	TOKEN_TYPE_PLUS
	TOKEN_TYPE_SEMICOLON
	TOKEN_TYPE_SLASH
	TOKEN_TYPE_STAR

	TOKEN_TYPE_BANG
	TOKEN_TYPE_BANG_EQUAL
	TOKEN_TYPE_EQUAL
	TOKEN_TYPE_EQUAL_EQUAL
	TOKEN_TYPE_GREATER
	TOKEN_TYPE_GREATER_EQUAL
	TOKEN_TYPE_LESS
	TOKEN_TYPE_LESS_EQUAL

	TOKEN_TYPE_IDENTIFIER
	TOKEN_TYPE_STRING
	TOKEN_TYPE_NUMBER

	TOKEN_TYPE_AND
	TOKEN_TYPE_ELSE
	TOKEN_TYPE_FALSE
	TOKEN_TYPE_FOR
	TOKEN_TYPE_FUNC
	TOKEN_TYPE_IF
	TOKEN_TYPE_NIL
	TOKEN_TYPE_OR
	TOKEN_TYPE_PRINT
	TOKEN_TYPE_RETURN
	TOKEN_TYPE_TRUE
	TOKEN_TYPE_VAR
	TOKEN_TYPE_WHILE

	TOKEN_TYPE_SPACE
	TOKEN_TYPE_NEW_LINE

	TOKEN_TYPE_ERROR
	TOKEN_TYPE_EOF
)

var InfixOpPrec = map[int]int{
	TOKEN_TYPE_PLUS:  PREC_ADD,
	TOKEN_TYPE_MINUS: PREC_ADD,
	TOKEN_TYPE_STAR:  PREC_MUL,
}

func ScannerIsAtEnd(s *Scanner) bool {
	return s.Current >= len(s.Source)
}

func ScannerMatch(s *Scanner, expected string) bool {
	if ScannerIsAtEnd(s) || (string(s.Source[s.Current:s.Current+1]) != expected) {
		return false
	}
	s.Current = s.Current + 1
	return true
}

func ScannerAdvance(s *Scanner) string {
	c := string(s.Source[s.Current : s.Current+1])
	s.Current = s.Current + 1
	return c
}

func ScannerPeek(s *Scanner) string {
	if ScannerIsAtEnd(s) {
		return ""
	}
	return string(s.Source[s.Current : s.Current+1])
}

func ScannerPeekNext(s *Scanner) string {
	if (s.Current + 1) >= len(s.Source) {
		return ""
	}
	return string(s.Source[s.Current+1 : s.Current+2])
}

func ScannerIsDigit(v string) bool {
	return (len(v) == 1) && (v[0] >= 0x30) && (v[0] <= 0x39)
}
func ScannerIsPrintable(v string) bool {
	return (len(v) == 1) && (v[0] >= 0x20) && (v[0] <= 0x7e)
}

func ScannerIsAlphabet(v string) bool {
	return (len(v) == 1) && (((v[0] >= 0x41) && (v[0] <= 0x5a)) || ((v[0] >= 0x61) && (v[0] <= 0x7a)) || (v[0] == 0x5f))
}

func ScannerIsTokenKeywordOrIdentifierType(s string) int {
	keywords := map[string]int{
		"if":     TOKEN_TYPE_IF,
		"else":   TOKEN_TYPE_ELSE,
		"for":    TOKEN_TYPE_FOR,
		"while":  TOKEN_TYPE_WHILE,
		"func":   TOKEN_TYPE_FUNC,
		"nil":    TOKEN_TYPE_NIL,
		"print":  TOKEN_TYPE_PRINT,
		"return": TOKEN_TYPE_RETURN,
		"true":   TOKEN_TYPE_TRUE,
		"false":  TOKEN_TYPE_FALSE,
		"var":    TOKEN_TYPE_VAR,
	}
	elem, ok := keywords[s]
	if ok {
		return elem
	}
	return TOKEN_TYPE_IDENTIFIER
}

func ScannerScanToken(s *Scanner) Token {
	if ScannerIsAtEnd(s) {
		return Token{TOKEN_TYPE_EOF, "", s.Line}
	}

	start := s.Current
	c := ScannerAdvance(s)

	switch c {
	case "(":
		return Token{TOKEN_TYPE_LEFT_PAREN, "(", s.Line}
	case ")":
		return Token{TOKEN_TYPE_RIGHT_PAREN, ")", s.Line}
	case "{":
		return Token{TOKEN_TYPE_LEFT_BRACE, "{", s.Line}
	case "}":
		return Token{TOKEN_TYPE_RIGHT_BRACE, "}", s.Line}
	case ",":
		return Token{TOKEN_TYPE_COMMA, ",", s.Line}
	case ".":
		return Token{TOKEN_TYPE_DOT, ".", s.Line}
	case "-":
		return Token{TOKEN_TYPE_MINUS, "-", s.Line}
	case "+":
		return Token{TOKEN_TYPE_PLUS, "+", s.Line}
	case ";":
		return Token{TOKEN_TYPE_SEMICOLON, ";", s.Line}
	case "*":
		return Token{TOKEN_TYPE_STAR, "*", s.Line}
	case "!":
		if ScannerMatch(s, "=") {
			return Token{TOKEN_TYPE_BANG_EQUAL, "!=", s.Line}
		} else {
			return Token{TOKEN_TYPE_BANG, "!", s.Line}
		}
	case "=":
		if ScannerMatch(s, "=") {
			return Token{TOKEN_TYPE_EQUAL_EQUAL, "==", s.Line}
		} else {
			return Token{TOKEN_TYPE_EQUAL, "=", s.Line}
		}
	case "<":
		if ScannerMatch(s, "=") {
			return Token{TOKEN_TYPE_LESS_EQUAL, "<=", s.Line}
		} else {
			return Token{TOKEN_TYPE_LESS, "<", s.Line}
		}
	case ">":
		if ScannerMatch(s, "=") {
			return Token{TOKEN_TYPE_GREATER_EQUAL, ">=", s.Line}
		} else {
			return Token{TOKEN_TYPE_GREATER, ">", s.Line}
		}
	case "\x22":
		for (!ScannerIsAtEnd(s)) && (ScannerPeek(s) != "\x22") {
			if ScannerIsPrintable(ScannerPeek(s)) {
				ScannerAdvance(s)
			} else {
				return Token{TOKEN_TYPE_ERROR, "Unexpected character in string literal", s.Line}
			}
		}
		if ScannerIsAtEnd(s) {
			return Token{TOKEN_TYPE_ERROR, "Unterminated string", s.Line}
		}
		ScannerAdvance(s)
		return Token{TOKEN_TYPE_STRING, string(s.Source[start+1 : s.Current-1]), s.Line}
	case "\x20":
		return Token{TOKEN_TYPE_SPACE, "\x20", s.Line}
	case "\x0a":
		s.Line = s.Line + 1
		return Token{TOKEN_TYPE_NEW_LINE, "\x0a", s.Line}
	default:
		if ScannerIsDigit(c) {
			for ScannerIsDigit(ScannerPeek(s)) {
				ScannerAdvance(s)
			}
			if (ScannerPeek(s) == ".") && ScannerIsDigit(ScannerPeekNext(s)) {
				ScannerAdvance(s)
				for ScannerIsDigit(ScannerPeek(s)) {
					ScannerAdvance(s)
				}
			}
			return Token{TOKEN_TYPE_NUMBER, string(s.Source[start:s.Current]), s.Line}
		} else if ScannerIsAlphabet(c) {
			for (!ScannerIsAtEnd(s)) && (ScannerIsAlphabet(ScannerPeek(s)) || ScannerIsDigit(ScannerPeek(s))) {
				ScannerAdvance(s)
			}
			tempString := string(s.Source[start:s.Current])
			t := ScannerIsTokenKeywordOrIdentifierType(tempString)
			if t == TOKEN_TYPE_IDENTIFIER {
				return Token{TOKEN_TYPE_IDENTIFIER, tempString, s.Line}
			} else {
				return Token{t, tempString, s.Line}
			}
		}
	}

	return Token{TOKEN_TYPE_ERROR, "Unknown error", s.Line}
}

func ScannerIsValidSource(source []byte) bool {
	for _, b := range source {
		if !((b == 0xa) || ((b >= 0x20) && (b <= 0x7e))) {
			return false
		}
	}
	return true
}

func ScannerScan(source []byte) []Token {
	if !ScannerIsValidSource(source) {
		log.Fatalln("Error while tokenization - Invalid source")
	}

	var tokens []Token
	s := Scanner{source, 0, 1}

	for {
		t := ScannerScanToken(&s)
		if (t.T == TOKEN_TYPE_SPACE) || (t.T == TOKEN_TYPE_NEW_LINE) {
			continue
		}
		tokens = append(tokens, t)
		if t.T == TOKEN_TYPE_EOF {
			break
		} else if t.T == TOKEN_TYPE_ERROR {
			log.Fatalln("Error while tokenization - Line", t.L, "-", t.S)
		}
	}

	return tokens
}

func CompilerAdvance(c *Compiler) {
	c.Current = c.Current + 1
}

func CompilerCurrent(c *Compiler) Token {
	return c.Source[c.Current]
}

func CompilerParseNumber(c *Compiler) float64 {
	s, err := strconv.ParseFloat(CompilerCurrent(c).S, 64)
	if err != nil {
		log.Fatalln("Error while compiling - Line", CompilerCurrent(c).L, "-", "Unable to parse", CompilerCurrent(c).S, "to float")
	}
	return s
}

func CompilerConsume(c *Compiler, t int, e string) {
	if CompilerCurrent(c).T == t {
		CompilerAdvance(c)
	} else {
		log.Fatalln("Error while compiling - Line", CompilerCurrent(c).L, "-", e)
	}
}

func CompilerIsInfixOp(t Token) bool {
	_, ok := InfixOpPrec[t.T]
	return ok
}

func CompilerParseExpression(c *Compiler) {
	var opStack []Token
	var tempToken Token
	var currentToken Token

	var parsingState int = 1

	for {
		currentToken = CompilerCurrent(c)

		if parsingState == 1 {
			if currentToken.T == TOKEN_TYPE_NUMBER {
				fmt.Println("PUSH NUMBER", currentToken.S)
				CompilerAdvance(c)
				parsingState = 3
			} else if currentToken.T == TOKEN_TYPE_MINUS {
				opStack = append(opStack, currentToken)
				CompilerAdvance(c)
				parsingState = 2
			}
		} else if parsingState == 2 {
			if currentToken.T == TOKEN_TYPE_NUMBER {
				fmt.Println("PUSH NUMBER", currentToken.S)
				for (len(opStack) != 0) && (opStack[len(opStack)-1].T == TOKEN_TYPE_MINUS) {
					tempToken, opStack = opStack[len(opStack)-1], opStack[:len(opStack)-1]
					fmt.Println("OP_NEGATE")
				}
				CompilerAdvance(c)
				parsingState = 3
			} else if currentToken.T == TOKEN_TYPE_MINUS {
				opStack = append(opStack, currentToken)
				CompilerAdvance(c)
			}
		} else if parsingState == 3 {
			if currentToken.T == TOKEN_TYPE_SEMICOLON {
				CompilerAdvance(c)
				break
			} else if CompilerIsInfixOp(currentToken) {
				for (len(opStack) != 0) && (InfixOpPrec[currentToken.T] <= InfixOpPrec[opStack[len(opStack)-1].T]) {
					tempToken, opStack = opStack[len(opStack)-1], opStack[:len(opStack)-1]
					fmt.Println(tempToken.S)
				}
				opStack = append(opStack, currentToken)
				CompilerAdvance(c)
				parsingState = 1
			}
		}
	}

	for len(opStack) != 0 {
		tempToken, opStack = opStack[len(opStack)-1], opStack[:len(opStack)-1]
		fmt.Println(tempToken.S)
	}
}

func ComplierCompile(tokens []Token) {
	fmt.Println(tokens)
	c := Compiler{tokens, 0}
	CompilerParseExpression(&c)
	CompilerConsume(&c, TOKEN_TYPE_EOF, "Error while compiling- Expect end of file")
}

func main() {
	d, e := os.ReadFile(os.Args[1])
	if e != nil {
		log.Fatal(e)
	}

	ComplierCompile(ScannerScan(d))
}
