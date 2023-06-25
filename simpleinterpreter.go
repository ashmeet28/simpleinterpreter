package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Token struct {
	t int
	s string
	l int
}

type Scanner struct {
	source  []byte
	current int
	line    int
}

type Compiler struct {
	source  []Token
	current int
}

const (
	PREC_ILLEGAL int = iota

	PREC_NONE
	PREC_OR
	PREC_AND
	PREC_COMP
	PREC_ADD
	PREC_MUL
	PREC_UNARY
	PREC_CALL
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

var TokenPrecedence = map[int]int{
	TOKEN_TYPE_PLUS: PREC_ADD,
	TOKEN_TYPE_STAR: PREC_MUL,
}

func ScannerIsAtEnd(s *Scanner) bool {
	return s.current >= len(s.source)
}

func ScannerMatch(s *Scanner, expected string) bool {
	if ScannerIsAtEnd(s) || (string(s.source[s.current:s.current+1]) != expected) {
		return false
	}
	s.current = s.current + 1
	return true
}

func ScannerAdvance(s *Scanner) string {
	c := string(s.source[s.current : s.current+1])
	s.current = s.current + 1
	return c
}

func ScannerPeek(s *Scanner) string {
	if ScannerIsAtEnd(s) {
		return ""
	}
	return string(s.source[s.current : s.current+1])
}

func ScannerPeekNext(s *Scanner) string {
	if (s.current + 1) >= len(s.source) {
		return ""
	}
	return string(s.source[s.current+1 : s.current+2])
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
		return Token{TOKEN_TYPE_EOF, "", s.line}
	}

	start := s.current
	c := ScannerAdvance(s)

	switch c {
	case "(":
		return Token{TOKEN_TYPE_LEFT_PAREN, "(", s.line}
	case ")":
		return Token{TOKEN_TYPE_RIGHT_PAREN, ")", s.line}
	case "{":
		return Token{TOKEN_TYPE_LEFT_BRACE, "{", s.line}
	case "}":
		return Token{TOKEN_TYPE_RIGHT_BRACE, "}", s.line}
	case ",":
		return Token{TOKEN_TYPE_COMMA, ",", s.line}
	case ".":
		return Token{TOKEN_TYPE_DOT, ".", s.line}
	case "-":
		return Token{TOKEN_TYPE_MINUS, "-", s.line}
	case "+":
		return Token{TOKEN_TYPE_PLUS, "+", s.line}
	case ";":
		return Token{TOKEN_TYPE_SEMICOLON, ";", s.line}
	case "*":
		return Token{TOKEN_TYPE_STAR, "*", s.line}
	case "!":
		if ScannerMatch(s, "=") {
			return Token{TOKEN_TYPE_BANG_EQUAL, "!=", s.line}
		} else {
			return Token{TOKEN_TYPE_BANG, "!", s.line}
		}
	case "=":
		if ScannerMatch(s, "=") {
			return Token{TOKEN_TYPE_EQUAL_EQUAL, "==", s.line}
		} else {
			return Token{TOKEN_TYPE_EQUAL, "=", s.line}
		}
	case "<":
		if ScannerMatch(s, "=") {
			return Token{TOKEN_TYPE_LESS_EQUAL, "<=", s.line}
		} else {
			return Token{TOKEN_TYPE_LESS, "<", s.line}
		}
	case ">":
		if ScannerMatch(s, "=") {
			return Token{TOKEN_TYPE_GREATER_EQUAL, ">=", s.line}
		} else {
			return Token{TOKEN_TYPE_GREATER, ">", s.line}
		}
	case "\x22":
		for (!ScannerIsAtEnd(s)) && (ScannerPeek(s) != "\x22") {
			if ScannerIsPrintable(ScannerPeek(s)) {
				ScannerAdvance(s)
			} else {
				return Token{TOKEN_TYPE_ERROR, "Unexpected character in string literal", s.line}
			}
		}
		if ScannerIsAtEnd(s) {
			return Token{TOKEN_TYPE_ERROR, "Unterminated string", s.line}
		}
		ScannerAdvance(s)
		return Token{TOKEN_TYPE_STRING, string(s.source[start+1 : s.current-1]), s.line}
	case "\x20":
		return Token{TOKEN_TYPE_SPACE, "\x20", s.line}
	case "\x0a":
		s.line = s.line + 1
		return Token{TOKEN_TYPE_NEW_LINE, "\x0a", s.line}
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
			return Token{TOKEN_TYPE_NUMBER, string(s.source[start:s.current]), s.line}
		} else if ScannerIsAlphabet(c) {
			for (!ScannerIsAtEnd(s)) && (ScannerIsAlphabet(ScannerPeek(s)) || ScannerIsDigit(ScannerPeek(s))) {
				ScannerAdvance(s)
			}
			tempString := string(s.source[start:s.current])
			t := ScannerIsTokenKeywordOrIdentifierType(tempString)
			if t == TOKEN_TYPE_IDENTIFIER {
				return Token{TOKEN_TYPE_IDENTIFIER, tempString, s.line}
			} else {
				return Token{t, tempString, s.line}
			}
		}
	}

	return Token{TOKEN_TYPE_ERROR, "Unknown error", s.line}
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
		if (t.t == TOKEN_TYPE_SPACE) || (t.t == TOKEN_TYPE_NEW_LINE) {
			continue
		}
		tokens = append(tokens, t)
		if t.t == TOKEN_TYPE_EOF {
			break
		} else if t.t == TOKEN_TYPE_ERROR {
			log.Fatalln("Error while tokenization - Line", t.l, "-", t.s)
		}
	}

	return tokens
}

func CompilerAdvance(c *Compiler) {
	c.current = c.current + 1
}

func CompilerCurrent(c *Compiler) Token {
	return c.source[c.current]
}

func CompilerParseNumber(c *Compiler) {
	s, err := strconv.ParseFloat(CompilerCurrent(c).s, 64)
	if err != nil {
		log.Fatalln("Error while compiling - Line", CompilerCurrent(c).l, "-", "Unable to parse", CompilerCurrent(c).s, "to float")
	}
	fmt.Println("OP_PUSH_FLOAT")
	fmt.Println(s)
}

func CompilerConsume(c *Compiler, t int, e string) {
	if CompilerCurrent(c).t == t {
		CompilerAdvance(c)
	} else {
		log.Fatalln("Error while compiling - Line", CompilerCurrent(c).l, "-", e)
	}
}

func CompilerParseExpression(c *Compiler) {
	var OpStack []Token
	var t Token

	for {
		if CompilerCurrent(c).t == TOKEN_TYPE_NUMBER {
			fmt.Println("PUSH NUMBER", CompilerCurrent(c).s)
			CompilerAdvance(c)
		} else if CompilerCurrent(c).t == TOKEN_TYPE_PLUS || CompilerCurrent(c).t == TOKEN_TYPE_STAR {
			for (len(OpStack) != 0) && TokenPrecedence[CompilerCurrent(c).t] <= TokenPrecedence[OpStack[len(OpStack)-1].t] {
				t, OpStack = OpStack[len(OpStack)-1], OpStack[:len(OpStack)-1]
				fmt.Println(t.s)
			}
			OpStack = append(OpStack, CompilerCurrent(c))
			CompilerAdvance(c)
		} else if CompilerCurrent(c).t == TOKEN_TYPE_SEMICOLON {
			CompilerAdvance(c)
			break
		}
	}
	for len(OpStack) != 0 {
		t, OpStack = OpStack[len(OpStack)-1], OpStack[:len(OpStack)-1]
		fmt.Println(t.s)
	}
}

func ComplierCompile(tokens []Token) {
	fmt.Println(tokens)
	c := Compiler{tokens, 0}
	CompilerParseExpression(&c)
	CompilerConsume(&c, TOKEN_TYPE_EOF, "Error - Expect end of file")
}

func main() {
	d, e := os.ReadFile(os.Args[1])
	if e != nil {
		log.Fatal(e)
	}

	ComplierCompile(ScannerScan(d))
}
