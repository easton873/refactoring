package hotel

import (
	"fmt"
	"strings"
)

var Printer OutputStrategy

type OutputStrategy interface {
	Println(...any)
	Printf(string, ...any)
}

type ConsolePrinter struct{}

func (c ConsolePrinter) Println(a ...any) {
	fmt.Println(a...)
}

func (c ConsolePrinter) Printf(s string, a ...any) {
	fmt.Printf(s, a...)
}

type StringPrinter struct {
	*strings.Builder
}

func (s StringPrinter) Println(a ...any) {
	s.WriteString(fmt.Sprintln(a...))
}

func (s StringPrinter) Printf(s2 string, a ...any) {
	s.WriteString(fmt.Sprintf(s2, a...))
}
