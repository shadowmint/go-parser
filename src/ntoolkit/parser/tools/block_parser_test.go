package tools_test

import (
	"testing"
	"ntoolkit/assert"
	"ntoolkit/parser/tools"
	"strings"
)

func TestBlockParserParse(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		p := tools.NewBlockParser()
		p.Parse("Hello World\"Chunk here\" ok ok \"Free chunk\" .... \"chunk on \n a newline\" OK?")
		tokens, err := p.Finished()
		T.Assert(err == nil)
		T.Assert(tokens != nil)
		T.Assert(tokens.Count() == 9)

		quoted := tokens.Front.FindNext(tools.TokenTypeQuotedBlock)
		T.Assert(quoted != nil)
		T.Assert(quoted.Children.Count() == 2)
		T.Assert(*quoted.Children.Front.Raw == "Chunk")
		T.Assert(*quoted.Children.Back.Raw == "here")
	})
}

func TestBlockParserDebug(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		p := tools.NewBlockParser()
		p.Parse("\"One Two\" ok ok \"\" .... \"Three four\"")
		tokens, err := p.Finished()

		T.Assert(err == nil)

		output := tokens.Debug(p.TokenTypes)
		T.Assert(strings.TrimSpace(output) == strings.TrimSpace(`
Tokens (6 found):
  Token: Quoted Block
    Tokens (2 found):
      Token: Block ('One')
      Token: Block ('Two')
  Token: Block ('ok')
  Token: Block ('ok')
  Token: Quoted Block
    Tokens (1 found):
      Token: Block ('')
  Token: Block ('....')
  Token: Quoted Block
    Tokens (2 found):
      Token: Block ('Three')
      Token: Block ('four')`))
	})
}
