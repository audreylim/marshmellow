package markdown

type ItemType int

const (
	ERROR ItemType = iota // 0
	EOF                   // 1

	SINGLESTAR // 2
	DOUBLESTAR // 3

	HEX       // 4
	STRINGLIT // 5

	WS      // 6
	NEWLINE // 7
)
