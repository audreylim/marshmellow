package mm_test

import (
	"strings"
	"testing"

	"github.com/audreylim/go-markdown"
)

// Ensure the scanner can scan tokens correctly.
func TestScanner_Scan(t *testing.T) {
	var tests = []struct {
		s   string
		tok mm.ItemType
		l   string
	}{
		{s: ``, tok: mm.EOF, l: "\x00"},
		{s: "#", tok: mm.HEX, l: "#"},
		{s: "*", tok: mm.SINGLESTAR, l: "*"},
		{s: "**", tok: mm.DOUBLESTAR, l: "**"},
		{s: " ", tok: mm.WS, l: ""},
		{s: "\t", tok: mm.WS, l: ""},
		{s: "\n", tok: mm.NEWLINE, l: ""},
		// FIXME: STRINGLIT test case not breaking test loop.
		//{s: "a", tok: mm.STRINGLIT, l: "a"},
	}

	for i, tt := range tests {
		s := mm.NewScanner(strings.NewReader(tt.s))
		tok, l := s.Scan()

		if tt.tok != tok {
			t.Errorf("%d. %q token mismatch: exp=%q got=%q <%q>", i, tt.s, tt.tok, tok, l)
		}
		if tt.l != l {
			t.Errorf("%d. %q literal mismatch: exp=%q got=%q", i, tt.s, tt.l, l)
		}
	}

}
