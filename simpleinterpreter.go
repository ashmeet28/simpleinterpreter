package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Token struct {
	t string
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

type ParseRule struct {
	prefix     func(c *Compiler)
	infix      func(c *Compiler)
	precedence int
}

var ParseRules map[string]ParseRule

var (
	PREC_NONE       int = 1
	PREC_ASSIGNMENT int = 2
	PREC_OR         int = 3
	PREC_AND        int = 4
	PREC_EQUALITY   int = 5
	PREC_COMPARISON int = 6
	PREC_TERM       int = 7
	PREC_FACTOR     int = 8
	PREC_UNARY      int = 9
	PREC_CALL       int = 10
	PREC_PRIMARY    int = 11
)

func ScannerIsAtEnd(s *Scanner) bool {
	return (*s).current >= len((*s).source)
}

func ScannerMatch(s *Scanner, expected string) bool {
	if ScannerIsAtEnd(s) || (string((*s).source[(*s).current:(*s).current+1]) != expected) {
		return false
	}
	(*s).current = (*s).current + 1
	return true
}

func ScannerAdvance(s *Scanner) string {
	var c string = string((*s).source[(*s).current : (*s).current+1])
	(*s).current = (*s).current + 1
	return c
}

func ScannerPeek(s *Scanner) string {
	if ScannerIsAtEnd(s) {
		return ""
	}
	return string((*s).source[(*s).current : (*s).current+1])
}

func ScannerPeekNext(s *Scanner) string {
	if ((*s).current + 1) >= len((*s).source) {
		return ""
	}
	return string((*s).source[(*s).current+1 : (*s).current+2])
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

func ScannerIsKeyword(v string) bool {
	keywords := []string{"if", "else", "for", "while", "func", "nil", "print", "return", "true", "false", "var"}
	for _, w := range keywords {
		if v == w {
			return true
		}
	}
	return false
}

func ScannerScanToken(s *Scanner) Token {
	if ScannerIsAtEnd(s) {
		return Token{"EOF", "", (*s).line}
	}

	var string_1 string
	var start int = (*s).current
	var c string = ScannerAdvance(s)

	switch c {
	case "(":
		return Token{"LEFT_PAREN", "(", (*s).line}
	case ")":
		return Token{"RIGHT_PAREN", ")", (*s).line}
	case "{":
		return Token{"LEFT_BRACE", "{", (*s).line}
	case "}":
		return Token{"RIGHT_BRACE", "}", (*s).line}
	case ",":
		return Token{"COMMA", ",", (*s).line}
	case ".":
		return Token{"DOT", ".", (*s).line}
	case "-":
		return Token{"MINUS", "-", (*s).line}
	case "+":
		return Token{"PLUS", "+", (*s).line}
	case ";":
		return Token{"SEMICOLON", ";", (*s).line}
	case "*":
		return Token{"STAR", "*", (*s).line}
	case "!":
		if ScannerMatch(s, "=") {
			return Token{"BANG_EQUAL", "!=", (*s).line}
		} else {
			return Token{"BANG", "!", (*s).line}
		}
	case "=":
		if ScannerMatch(s, "=") {
			return Token{"EQUAL_EQUAL", "==", (*s).line}
		} else {
			return Token{"EQUAL", "=", (*s).line}
		}
	case "<":
		if ScannerMatch(s, "=") {
			return Token{"LESS_EQUAL", "<=", (*s).line}
		} else {
			return Token{"LESS", "<", (*s).line}
		}
	case ">":
		if ScannerMatch(s, "=") {
			return Token{"GREATER_EQUAL", ">=", (*s).line}
		} else {
			return Token{"GREATER", ">", (*s).line}
		}
	case "\x22":
		for (!ScannerIsAtEnd(s)) && (ScannerPeek(s) != "\x22") {
			if ScannerIsPrintable(ScannerPeek(s)) {
				ScannerAdvance(s)
			} else {
				return Token{"ERROR", "Unexpected character in string literal", (*s).line}
			}
		}

		if ScannerIsAtEnd(s) {
			return Token{"ERROR", "Unterminated string", (*s).line}
		}

		ScannerAdvance(s)
		return Token{"STRING", string((*s).source[start+1 : (*s).current-1]), (*s).line}
	case "\x20":
		return Token{"SPACE", "\x20", (*s).line}
	case "\x0a":
		(*s).line = (*s).line + 1
		return Token{"NEW_LINE", "\x0a", (*s).line}
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
			return Token{"NUMBER", string((*s).source[start:(*s).current]), (*s).line}
		} else if ScannerIsAlphabet(c) {
			for (!ScannerIsAtEnd(s)) && (ScannerIsAlphabet(ScannerPeek(s)) || ScannerIsDigit(ScannerPeek(s))) {
				ScannerAdvance(s)
			}
			string_1 = string((*s).source[start:(*s).current])
			if ScannerIsKeyword(string_1) {
				return Token{string_1, string_1, (*s).line}
			} else {
				return Token{"IDENTIFIER", string_1, (*s).line}
			}
		}
	}

	return Token{"ERROR", "Internal error", (*s).line}
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
	var s Scanner = Scanner{source, 0, 1}
	var t Token

	for {
		t = ScannerScanToken(&s)
		if (t.t == "SPACE") || (t.t == "NEW_LINE") {
			continue
		}
		tokens = append(tokens, t)
		if t.t == "EOF" {
			break
		} else if t.t == "ERROR" {
			log.Fatalln("Error while tokenization - Line", t.l, "-", t.s)
		}
	}

	return tokens
}

func CompilerAdvance(c *Compiler) {
	(*c).current = (*c).current + 1
}

func CompilerPrevious(c *Compiler) Token {
	return (*c).source[(*c).current-1]
}

func CompilerCurrent(c *Compiler) Token {
	return (*c).source[(*c).current]
}

func CompilerNumber(c *Compiler) {
	fmt.Println("Comp Number", CompilerPrevious(c), CompilerCurrent(c))

	s, err := strconv.ParseFloat(CompilerPrevious(c).s, 64)
	if err != nil {
		log.Fatalln("Error while compiling - Line", CompilerPrevious(c).l, "-", "Unable to parse", CompilerPrevious(c).s, "to float.")
	}
	fmt.Println("OP_PUSH_FLOAT")
	fmt.Println(s)
}

func CompilerUnary(c *Compiler) {
	fmt.Println("unary c", CompilerPrevious(c), CompilerCurrent(c))

	var t string = CompilerPrevious(c).t

	CompilerParseExpression(c, PREC_UNARY)

	fmt.Println("unary opert", CompilerPrevious(c), CompilerCurrent(c))

	switch t {
	case "MINUS":
		fmt.Println("OP_NEGATE")
	}
}

func CompilerBinary(c *Compiler) {
	fmt.Println("Binary c", CompilerPrevious(c), CompilerCurrent(c))
	var operatorType string = CompilerPrevious(c).t
	var rule ParseRule = CompilerGetRule(operatorType)
	CompilerParseExpression(c, rule.precedence+1)

	fmt.Println("Binary opert", CompilerPrevious(c), CompilerCurrent(c))

	switch operatorType {
	case "PLUS":
		fmt.Println("OP_ADD")
	case "MINUS":
		fmt.Println("OP_SUBTRACT")
	case "STAR":
		fmt.Println("OP_MULTIPLY")
	case "SLASH":
		fmt.Println("OP_DIVIDE")
	}
}

func CompilerGetRule(t string) ParseRule {
	return ParseRules[t]
}

func CompilerConsume(c *Compiler, t string, e string) {
	if CompilerCurrent(c).t == t {
		CompilerAdvance(c)
	} else {
		log.Fatalln("Error while compiling - Line", CompilerCurrent(c).l, "-", e)
	}
}

func CompilerGrouping(c *Compiler) {
	fmt.Println("Group ", CompilerPrevious(c), CompilerCurrent(c))
	CompilerParseExpression(c, PREC_ASSIGNMENT)
	CompilerConsume(c, "RIGHT_PAREN", "Expect ) after expression.")
}

func CompilerParseExpression(c *Compiler, precedence int) {
	CompilerAdvance(c)
	fmt.Println("ParseE prefix", CompilerPrevious(c), CompilerCurrent(c))

	var prefixRule func(c *Compiler) = CompilerGetRule(CompilerPrevious(c).t).prefix

	if prefixRule == nil {
		log.Fatalln("Unexpected expression.")
	}

	prefixRule(c)

	fmt.Println("ParseE infix", CompilerPrevious(c), CompilerCurrent(c))
	fmt.Println("for loop before", precedence, CompilerGetRule(CompilerCurrent(c).t).precedence)

	for precedence <= CompilerGetRule(CompilerCurrent(c).t).precedence {

		fmt.Println("for loop", precedence, CompilerGetRule(CompilerCurrent(c).t).precedence)
		fmt.Println("ParseE infix in for", CompilerPrevious(c), CompilerCurrent(c))
		CompilerAdvance(c)
		var infixRule func(c *Compiler) = CompilerGetRule(CompilerPrevious(c).t).infix
		infixRule(c)
	}
}

func ComplierCompile(tokens []Token) {
	fmt.Println(tokens)
	var c Compiler = Compiler{tokens, 0}
	CompilerParseExpression(&c, PREC_ASSIGNMENT)
	CompilerConsume(&c, "EOF", "Expect end of file.")
}

func CompilerInit() {
	ParseRules = map[string]ParseRule{
		"LEFT_PAREN":  {CompilerGrouping, nil, PREC_NONE},
		"RIGHT_PAREN": {nil, nil, PREC_NONE},
		"NUMBER":      {CompilerNumber, nil, PREC_NONE},
		"MINUS":       {CompilerUnary, CompilerBinary, PREC_TERM},
		"PLUS":        {nil, CompilerBinary, PREC_TERM},
		"SLASH":       {nil, CompilerBinary, PREC_FACTOR},
		"STAR":        {nil, CompilerBinary, PREC_FACTOR},
		"EOF":         {nil, nil, PREC_NONE},
	}
}

func main() {
	d, e := os.ReadFile(os.Args[1])
	if e != nil {
		log.Fatal(e)
	}

	CompilerInit()
	ComplierCompile(ScannerScan(d))
}
