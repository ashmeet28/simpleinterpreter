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

func ScannerScanToken(s *Scanner) Token {
	if ScannerIsAtEnd(s) {
		return Token{"eof", "", (*s).line}
	}
	var start int
	start = (*s).current
	var c string = ScannerAdvance(s)

	switch c {
	case "(":
		return Token{"left_paren", "(", (*s).line}
	case ")":
		return Token{"right_paren", ")", (*s).line}
	case "{":
		return Token{"left_brace", "{", (*s).line}
	case "}":
		return Token{"right_brace", "}", (*s).line}
	case ",":
		return Token{"comma", ",", (*s).line}
	case ".":
		return Token{"dot", ".", (*s).line}
	case "-":
		return Token{"minus", "-", (*s).line}
	case "+":
		return Token{"plus", "+", (*s).line}
	case ";":
		return Token{"semicolon", ";", (*s).line}
	case "*":
		return Token{"star", "*", (*s).line}
	case "!":
		if ScannerMatch(s, "=") {
			return Token{"bang_equal", "!=", (*s).line}
		} else {
			return Token{"bang", "!", (*s).line}
		}
	case "=":
		if ScannerMatch(s, "=") {
			return Token{"equal_equal", "==", (*s).line}
		} else {
			return Token{"equal", "=", (*s).line}
		}
	case "<":
		if ScannerMatch(s, "=") {
			return Token{"less_equal", "<=", (*s).line}
		} else {
			return Token{"less", "<", (*s).line}
		}
	case ">":
		if ScannerMatch(s, "=") {
			return Token{"greater_equal", ">=", (*s).line}
		} else {
			return Token{"greater", ">", (*s).line}
		}
	case "\x22":
		for (!(ScannerIsAtEnd(s))) && (ScannerPeek(s) != "\x22") {
			if ScannerIsPrintable(ScannerPeek(s)) {
				ScannerAdvance(s)
			} else {
				log.Fatalln("Error while tokenization [ Line", (*s).line, "]", "- Unexpected character in string literal")
			}
		}

		if ScannerIsAtEnd(s) {
			log.Fatalln("Error while tokenization [ Line", (*s).line, "]", "- Unterminated string")
		}

		ScannerAdvance(s)
		return Token{"string", string((*s).source[start+1 : (*s).current-1]), (*s).line}
	case "\x20":
		return Token{"whitespace", "\x20", (*s).line}
	case "\x0d":
		return Token{"whitespace", "\x0d", (*s).line}
	case "\x09":
		return Token{"whitespace", "\x09", (*s).line}
	case "\x0a":
		(*s).line = (*s).line + 1
		return Token{"new_line", "\x0a", (*s).line}
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
			return Token{"number", string((*s).source[start:(*s).current]), (*s).line}
		}

		log.Fatalln("Error while tokenization [ Line", (*s).line, "]", "- Unexpected character", c)
	}

	return Token{"error", "", (*s).line}
}

func main() {
	d, e := os.ReadFile(os.Args[1])
	if e != nil {
		log.Fatal(e)
	}

	var s = Scanner{d, 0, 1}
	for {
		t := ScannerScanToken(&s)
		fmt.Println(t)
		if t.t == "eof" {
			break
		}
	}
}
