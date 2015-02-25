package mm_test

import (
	"strings"
	"testing"

	mm "github.com/audreylim/go-markdown"
)

// Ensure the scanner can scan tokens correctly.
func TestScanner_Scan(t *testing.T) {
	var tests = []struct {
		s   string
		tok md.ItemType
		l   string
	}{
		{s: ``, tok: md.EOF, l: "\x00"},
		{s: "#", tok: md.HEX, l: "#"},
		{s: "*", tok: md.SINGLESTAR, l: "*"},
		{s: "**", tok: md.DOUBLESTAR, l: "**"},
		{s: " ", tok: md.WS, l: ""},
		{s: "\t", tok: md.WS, l: ""},
		{s: "\n", tok: md.NEWLINE, l: ""},
		// FIXME: STRINGLIT test case not breaking test loop.
		//{s: "a", tok: md.STRINGLIT, l: "a"},
	}

	for i, tt := range tests {
		s := md.NewScanner(strings.NewReader(tt.s))
		tok, l := s.Scan()

		if tt.tok != tok {
			t.Errorf("%d. %q token mismatch: exp=%q got=%q <%q>", i, tt.s, tt.tok, tok, l)
		}
		if tt.l != l {
			t.Errorf("%d. %q literal mismatch: exp=%q got=%q", i, tt.s, tt.l, l)
		}
	}

}
