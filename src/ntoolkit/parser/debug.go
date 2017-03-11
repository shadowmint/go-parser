package parser

import (
	"fmt"
	"strings"
)

// debugIndentLine prints a string indented by the given indentation
func debugIndentLine(value string, indent ...int) string {
	indentValue := 0
	if len(indent) > 0 {
		indentValue = indent[0]
	}
	space := ""
	for i := 0; i < indentValue; i++ {
		space += "  "
	}
	return fmt.Sprintf("%s%s", space, value)
}

// debugIndentBlock splits a block on newlines and indents them all, then returns the result
func debugIndentBlock(block string, indent ...int) string {
	buffer := ""
	lines := strings.Split(block, "\n")
	for i := range lines {
		if len(lines[i]) > 0 {
			buffer += debugIndentLine(lines[i], indent...)
		}
		if i != len(lines)-1 {
			buffer += "\n"
		}
	}
	return buffer
}

// debugIndent returns the indent value
func debugIndent(indent ...int) int {
	if len(indent) > 0 {
		return indent[0]
	}
	return 0
}

// debugIndentPlus returns the indent value + 1
func debugIndentPlus(indent ...int) int {
	return debugIndent(indent...) + 1
}
