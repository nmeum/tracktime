package parser

import (
	"fmt"
)

type ParserError struct {
	Name string
	Line uint
	Msg  string
}

func (p ParserError) Error() string {
	return fmt.Sprintf("%s:%d %s", p.Name, p.Line, p.Msg)
}
