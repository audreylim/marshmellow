package md

import (
	"bufio"
	"bytes"
	"io"
)

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// Scan() returns token and corresponding string literal.
func (s *Scanner) Scan() (ItemType, string) {
	// Reads the next rune.
	ch := s.read()

	switch {
	case isWhitespace(ch):
		return s.scanWhitespace()
	case isNewLine(ch):
		s.unread()
		return s.scanNewLine()
	case isStar(ch):
		s.unread()
		return s.scanStar()
	case isHex(ch):
		s.unread()
		a, b := s.scanHex()
		return a, b
	case ch == eof:
		return EOF, string(eof)
	}

	s.unread()
	return s.scanStringLiteral()
}

func (s *Scanner) scanWhitespace() (ItemType, string) {
	// Create a buffer and read the current ch into it.
	var buf bytes.Buffer

	// Read every subsequent ws ch into the buffer.
	// Non ws ch will cause the loop to exit.
	for {
		if ch := s.read(); !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, ""
}

func (s *Scanner) scanNewLine() (ItemType, string) {
	var buf bytes.Buffer

	for {
		if ch := s.read(); !isNewLine(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return NEWLINE, ""
}

func (s *Scanner) scanHex() (ItemType, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); !isHex(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return HEX, buf.String()
}

func (s *Scanner) scanStar() (ItemType, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	ch := s.read()
	if !isStar(ch) {
		s.unread()
		return SINGLESTAR, buf.String()
	}

	buf.WriteRune(ch)
	return DOUBLESTAR, buf.String()
}

func (s *Scanner) scanStringLiteral() (ItemType, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); !isStringLiteral(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return STRINGLIT, buf.String()
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread() { _ = s.r.UnreadRune() }

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' }

func isNewLine(ch rune) bool { return ch == '\n' }

func isHex(ch rune) bool { return ch == '#' }

func isStar(ch rune) bool { return ch == '*' }

func isStringLiteral(ch rune) bool {
	return (ch != '\t' && ch != '\n' && ch != '#' && ch != '*')
}

var eof = rune(0)
