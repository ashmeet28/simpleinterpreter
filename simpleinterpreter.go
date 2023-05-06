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
	source  string
	current int
	line    int
}

func ScannerIsAtEnd(s *Scanner) bool {
	return (*s).current >= len((*s).source)
}

func ScannerMatch(s *Scanner, expected string) bool {
	if ScannerIsAtEnd(s) {
		return false
	}
	if string([]byte{(*s).source[(*s).current]}) != expected {
		return false
	}
	(*s).current = (*s).current + 1
	return true
}

func ScannerAdvance(s *Scanner) string {
	var c string = string([]byte{(*s).source[(*s).current]})
	(*s).current = (*s).current + 1
	return c
}

func ScannerPeek(s *Scanner) string {
	if ScannerIsAtEnd(s) {
		return "\x00"
	}
	return string([]byte{(*s).source[(*s).current]})
}

func ScannerPeekNext(s *Scanner) string {
	if ((*s).current + 1) >= len((*s).source) {
		return "\x00"
	}
	return string([]byte{(*s).source[(*s).current+1]})
}

func ScannerIsDigit(v string) bool {
	return (len(v) == 1) && (v[0] >= 0x30) && (v[0] <= 0x39)
}

func ScannerMakeToken(s *Scanner, TokenT string, TokenS string) Token {
	return Token{TokenT, TokenS, (*s).line}
}

func ScannerScanToken(s *Scanner) Token {
	if ScannerIsAtEnd(s) {
		return ScannerMakeToken(s, "eof", "")
	}
	var c string = ScannerAdvance(s)
	switch c {
	case "(":
		return ScannerMakeToken(s, "left_paren", "")
	case ")":
		return ScannerMakeToken(s, "right_paren", "")
	case "{":
		return ScannerMakeToken(s, "left_brace", "")
	case "}":
		return ScannerMakeToken(s, "right_brace", "")
	case ",":
		return ScannerMakeToken(s, "comma", "")
	case ".":
		return ScannerMakeToken(s, "dot", "")
	case "-":
		return ScannerMakeToken(s, "minus", "")
	case "+":
		return ScannerMakeToken(s, "plus", "")
	case ";":
		return ScannerMakeToken(s, "semicolon", "")
	case "*":
		return ScannerMakeToken(s, "star", "")
	case "!":
		if ScannerMatch(s, "=") {
			return ScannerMakeToken(s, "bang_equal", "")
		} else {
			return ScannerMakeToken(s, "bang", "")
		}
	case "=":
		if ScannerMatch(s, "=") {
			return ScannerMakeToken(s, "equal_equal", "")
		} else {
			return ScannerMakeToken(s, "equal", "")
		}
	case "<":
		if ScannerMatch(s, "=") {
			return ScannerMakeToken(s, "less_equal", "")
		} else {
			return ScannerMakeToken(s, "less", "")
		}
	case ">":
		if ScannerMatch(s, "=") {
			return ScannerMakeToken(s, "greater_equal", "")
		} else {
			return ScannerMakeToken(s, "greater", "")
		}
	case "\x22":
	case "\x20":
	case "\x0d":
	case "\x09":
	case "\x0a":
		return ScannerMakeToken(s, "new_line", "")
	default:
		log.Fatalln("Unexpected character", c, "at line", (*s).line)
	}
	return ScannerMakeToken(s, "", "")
}

func main() {
	d, e := os.ReadFile(os.Args[1])
	if e != nil {
		log.Fatal(e)
	}

	var s = Scanner{string(d), 0, 1}
	fmt.Println(ScannerScanToken(&s))
}
