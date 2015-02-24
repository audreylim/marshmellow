package md

import (
	"io"
	"strings"
)

type Parser struct {
	s *Scanner

	// Temporary storage of values.
	tempSlice  []string
	tempString string
	tempHeader string

	Formatter []string
	Stringlit []string
}

func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

type stateFn func(*Parser) stateFn

func (p *Parser) Parse() {
	p.run(p.stateParse())
}

// Chain states until EOF.
func (p *Parser) run(state stateFn) stateFn {
	newstate := state
	if newstate != nil {
		p.run(newstate)
	}
	return nil
}

// State functions.
func (p *Parser) stateParse() stateFn {
	// t stands for token, and l stands for string literal returned by Scan() functions.
	t, l := p.s.Scan()

	switch t {
	case HEX:
		p.tempHeader = l
		return p.stateHeader()
	case SINGLESTAR:
		return p.stateSingleStar()
	case DOUBLESTAR:
		return p.stateDoubleStar()
	case STRINGLIT:
		p.tempString = l
		return p.statePara()
	case EOF:
		return nil
	case NEWLINE:
		return p.stateParse()
	}

	return nil
}

func (p *Parser) stateHeader() stateFn {
	t1, l1 := p.s.Scan()

	switch t1 {
	case WS:
		_, l2 := p.s.Scan()
		if len(p.tempHeader) < 7 {
			p.Formatter = append(p.Formatter, p.tempHeader)
		} else {
			// Presume h6 if more than 6#.
			p.Formatter = append(p.Formatter, "######")
		}
		p.Stringlit = append(p.Stringlit, l2)

		p.tempHeader = ""
		return p.stateParse()

	default:
		// FIXME: not redirected here when ###noWS after **string*** line. Returns WS as next token instead.
		p.tempSlice = append(p.tempSlice, p.tempHeader+l1)
		return p.statePara()
	}

	return p.stateParse()
}

func (p *Parser) stateDoubleStar() stateFn {
	t1, l1 := p.s.Scan()

	switch t1 {
	case WS:
		_, l2 := p.s.Scan()
		p.Formatter = append(p.Formatter, "para")
		p.Stringlit = append(p.Stringlit, "**"+l2)
		return p.stateParse()

	case STRINGLIT:
		if ta, _ := p.s.Scan(); ta == DOUBLESTAR {
			p.Formatter = append(p.Formatter, "para")
			p.Stringlit = append(p.Stringlit, "<b>"+l1+"</b>")
			return p.stateParse()
		}

	case NEWLINE:
		p.tempSlice = append(p.tempSlice, p.tempString)

	default:
		p.tempSlice = append(p.tempSlice, "**"+l1)
		return p.statePara()
	}

	p.Formatter = append(p.Formatter, "**")
	p.Stringlit = append(p.Stringlit, strings.Join(p.tempSlice, ""))

	p.tempSlice = []string{}

	return p.stateParse()
}

func (p *Parser) stateSingleStar() stateFn {
	t1, l1 := p.s.Scan()

	switch t1 {
	case WS:
		_, l2 := p.s.Scan()
		p.tempSlice = append(p.tempSlice, "<li>"+l2+"</li>\n")
		t3, _ := p.s.Scan()
		if t3 == NEWLINE {
			t4, l4 := p.s.Scan()
			if t4 == SINGLESTAR {
				return p.stateSingleStar()
			}
			if t4 == STRINGLIT {
				p.Formatter = append(p.Formatter, "bullet")
				p.Stringlit = append(p.Stringlit, strings.Join(p.tempSlice, ""))

				p.tempSlice = []string{}
				p.tempString = l4
				return p.statePara()
			}
		}

	case SINGLESTAR:
		return p.stateSingleStar()

	case STRINGLIT:
		if tok, _ := p.s.Scan(); tok == SINGLESTAR {
			p.Formatter = append(p.Formatter, "para")
			p.Stringlit = append(p.Stringlit, "<i>"+l1+"</i>")

			return p.stateParse()
		}

		joinSlice := strings.Join(p.tempSlice, "")
		if joinSlice != "" {
			p.Formatter = append(p.Formatter, "bullet")
			p.Stringlit = append(p.Stringlit, joinSlice)
			p.tempSlice = []string{}
		}

		p.Formatter = append(p.Formatter, "para")
		p.Stringlit = append(p.Stringlit, "*"+l1)

		return p.stateParse()

	case NEWLINE:
		p.tempSlice = append(p.tempSlice, p.tempString)
	}

	p.Formatter = append(p.Formatter, "*")
	p.Stringlit = append(p.Stringlit, strings.Join(p.tempSlice, ""))

	p.tempSlice = []string{}
	return p.stateParse()
}

func (p *Parser) statePara() stateFn {
	t1, l1 := p.s.Scan()
	p.tempSlice = append(p.tempSlice, p.tempString)
	p.tempString = ""

	switch t1 {
	case WS:
		p.tempSlice = append(p.tempSlice, " ")
		return p.statePara()

	case SINGLESTAR:
		typ, inlineString := p.checkIfItalics(t1, l1)
		if typ == "italics" {
			p.tempSlice = append(p.tempSlice, "<i>"+inlineString+"</i>")
			return p.statePara()
		} else if typ == "new line" {
			p.tempSlice = append(p.tempSlice, l1+inlineString)
		} else {
			p.tempSlice = append(p.tempSlice, l1+inlineString)
			return p.statePara()
		}

	case DOUBLESTAR:
		typ, inlineString := p.checkIfBold(t1, l1)
		if typ == "bold" {
			p.tempSlice = append(p.tempSlice, "<b>"+inlineString+"</b>")
			return p.statePara()
		} else {
			p.tempSlice = append(p.tempSlice, l1+inlineString)
		}

	case STRINGLIT:
		p.tempSlice = append(p.tempSlice, l1)
		return p.statePara()
	}

	p.Formatter = append(p.Formatter, "para")
	p.Stringlit = append(p.Stringlit, strings.Join(p.tempSlice, ""))

	p.tempSlice = []string{}
	return p.stateParse()
}

func (p *Parser) checkIfItalics(t ItemType, l string) (string, string) {
	t2, l2 := p.s.Scan()
	t3, _ := p.s.Scan()

	if t3 != SINGLESTAR || t2 == WS {
		if t3 == NEWLINE {
			return "new line", l2
		}
		return "not italics", l2
	}

	if t2 == NEWLINE || t3 == NEWLINE {
		return "new line", ""
	}

	return "italics", l2
}

func (p *Parser) checkIfBold(t1 ItemType, l string) (string, string) {
	_, l2 := p.s.Scan()
	t3, _ := p.s.Scan()
	switch t3 {
	case DOUBLESTAR:
		return "bold", l2
	}

	return "not bold", l2
}
