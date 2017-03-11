package parser_test

import (
	"testing"
	"ntoolkit/assert"
	"ntoolkit/parser"
	"fmt"
	"strings"
)

func TestCollapse(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		tokens := parser.NewTokens(nil)
		tokens.Push(&parser.Token{Type: 1})
		tokens.Push(&parser.Token{Type: 2})
		tokens.Push(&parser.Token{Type: 3})
		tokens.Push(&parser.Token{Type: 4})
		tokens.Push(&parser.Token{Type: 5})

		start := tokens.Front.FindNext(2)
		end := tokens.Front.FindNext(4)
		token := tokens.Collapse(start, end)

		T.Assert(tokens.Count() == 3)
		T.Assert(token.Children.Count() == 3)
		T.Assert(token.Children.Front.Type == 2)
		T.Assert(token.Children.Front.Next.Type == 3)
		T.Assert(token.Children.Front.Next.Next.Type == 4)

		output := tokens.Debug(func(T uint32) string {
			return fmt.Sprintf("TokenType%d", T)
		})

		expected := `
Tokens (3 found):
  Token: TokenType1
  Token: TokenType0
    Tokens (3 found):
      Token: TokenType2
      Token: TokenType3
      Token: TokenType4
  Token: TokenType5`

		T.Assert(strings.TrimSpace(output) == strings.TrimSpace(expected))
	})
}


