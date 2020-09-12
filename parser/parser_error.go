package parser

import (
	"fmt"
)

type ParserError struct {
	Line uint
	Msg  string
}

func (p ParserError) Error() string {
	return fmt.Sprintf("%s:%d %s", "stdin", p.Line, p.Msg)
}
