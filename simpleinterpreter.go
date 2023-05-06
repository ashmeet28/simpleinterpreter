package main

import (
	"fmt"
	"log"
	"os"
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
	var start int
	start = (*s).current
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
		if !(b == 0xa || ((b >= 0x20) && (b <= 0x7e))) {
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
	var s = Scanner{source, 0, 1}
	var t Token

	for {
		t = ScannerScanToken(&s)
		tokens = append(tokens, t)
		if t.t == "EOF" {
			break
		} else if t.t == "ERROR" {
			log.Fatalln("Error while tokenization - Line", t.l, "-", t.s)
		}
	}

	return tokens
}

func ComplierCompile(tokens []Token) {
	fmt.Println(tokens)
}

func main() {
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	ComplierCompile(ScannerScan(data))
}
