package tools

import (
	"ntoolkit/parser"
	"strings"
	"sync"
)

const TokenTypeBlock = 1
const TokenTypeQuote = 2
const TokenTypeQuotedBlock = 3

type SpaceQuoteTokenizer struct {
	next   *parser.Token
	tokens *parser.Tokens
}

// Enter starts processing a data stream with a context object and resets the internal tokenizer state
func (t *SpaceQuoteTokenizer) Enter(context *parser.Tokens) {
	t.tokens = context
	t.next = nil
}

// Process reads incoming string data and pushes tokens to the context
func (t *SpaceQuoteTokenizer) Process(data string) {
	parts := strings.Split(data, "")
	for i := range parts {
		if token := t.process(parts[i]); token != nil {
			t.tokens.Push(token)
		}
	}
}

// Close should close and resolve the tokenizer.
func (t *SpaceQuoteTokenizer) Close() {
	if t.next != nil {
		t.tokens.Push(t.next)
		t.next = nil
	}
}

// Process reads incoming string data and pushes tokens to the context
func (t *SpaceQuoteTokenizer) process(sym string) *parser.Token {

	// Spaces always end the current token, if there is one, and leave no pending token
	if sym == " " {
		if t.next != nil {
			rtn := t.next
			t.next = nil
			return rtn
		}
	}

	// A quote symbol ends a word token, leaving a quote on the stack.
	// A quote symbol on its own yields nothing, leaving a quote on the stack.
	if sym == "\"" {
		rtn := t.next
		t.next = &parser.Token{Type: TokenTypeQuote, Raw: nil}
		return rtn
	}

	// If there's nothing on the stack, create a new symbol
	if t.next == nil {
		t.next = &parser.Token{Type: parser.TokenTypeNone, Raw: &sym}
		return nil
	}

	// If the top of the stack already exists and is None type, append to it
	if t.next != nil && t.next.Type == parser.TokenTypeNone {
		*t.next.Raw += sym
		return nil
	}

	// If the top of the stack is a quote, return that and create a new symbol
	if t.next != nil && t.next.Type == TokenTypeQuote {
		rtn := t.next
		t.next = &parser.Token{Type: parser.TokenTypeNone, Raw: &sym}
		return rtn
	}

	// No idea how we might get here, do nothing.
	return nil
}

type BlockClassifier struct {
}

func (t *BlockClassifier) Classify(token *parser.Token) bool {
	modified := false

	// If there is a quote symbol, push everything until the next quote into a child quoted block
	if token.Type == TokenTypeQuote {
		end := token.FindNext(TokenTypeQuote)
		if end != nil {
			if token.Parent.Parent == nil || token.Parent.Parent.Type != TokenTypeQuotedBlock {
				newToken := token.Parent.Collapse(token, end)
				newToken.Type = TokenTypeQuotedBlock

				// Remove the quote symbols from the quoted block
				newToken.Children.Shift()
				newToken.Children.Pop()

				// Empty quoted blocks need an empty string value in them
				if newToken.Children.Count() == 0 {
					empty := ""
					newToken.Children.Push(&parser.Token{Type: TokenTypeBlock, Raw: &empty})
				}

				modified = true
			}
		}
	} else if token.Type == parser.TokenTypeNone {
		// Everything else just gets to be a block
		modified = true
		token.Type = TokenTypeBlock
	}
	return modified
}

type BlockParser struct {
	tokenizer  *SpaceQuoteTokenizer
	classifier *BlockClassifier
	promise    *parser.DeferredTokens
	handle     func(data string, done bool)
	TokenTypes func(T uint32) string
}

func NewBlockParser() *BlockParser {
	return &BlockParser{
		TokenTypes: BlockParserTokenTypes,
		tokenizer:  &SpaceQuoteTokenizer{},
		classifier: &BlockClassifier{},
		promise:    nil,
		handle:     nil}
}

func (bp *BlockParser) Parse(chunk string) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	if bp.promise == nil {
		bp.promise = parser.Parse(bp.tokenizer, bp.classifier, func(handle func(data string, done bool)) {
			bp.handle = handle
			wg.Done()
		})
	}
	wg.Wait()
	bp.handle(chunk, false)
}

func (bp *BlockParser) Finished() (*parser.Tokens, error) {
	var err error = nil
	var rtn *parser.Tokens = nil
	wg := &sync.WaitGroup{}
	wg.Add(1)
	bp.promise.Then(func(t *parser.Tokens) {
		rtn = t
		wg.Done()
	}, func(terr error) {
		err = terr
		wg.Done()
	})
	bp.handle("", true)
	return rtn, err
}

func BlockParserTokenTypes(T uint32) string {
	switch T {
	case TokenTypeBlock:
		return "Block"
	case TokenTypeQuote:
		return "Quote Symbol"
	case TokenTypeQuotedBlock:
		return "Quoted Block"
	default:
		return "Unknown"
	}
}
