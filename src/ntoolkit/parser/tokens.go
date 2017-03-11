package parser

import (
	"fmt"
)

// Tokens is a management container for a set of Token objects
type Tokens struct {
	Parent *Token
	Front  *Token
	Back   *Token
}

// NewTokens returns a new empty tokens object
func NewTokens(parent *Token) *Tokens {
	return &Tokens{Parent: parent, Front: nil, Back: nil}
}

// Count the tokens
func (tokens *Tokens) Count() int {
	count := 0
	for marker := tokens.Front; marker != nil; marker = marker.Next {
		count += 1
	}
	return count
}

// Push a new token onto the end of the token list
func (tokens *Tokens) Push(token *Token) {
	token.Parent = tokens
	if tokens.Back == nil {
		tokens.Back = token.Last()
		tokens.Front = token.First()
	} else {
		tokens.Back.InsertAfter(token)
		tokens.Back = token.Last()
	}
}

// Pop a token off the end of the token list
func (tokens *Tokens) Pop() *Token {
	back := tokens.Back
	if tokens.Back == nil {
		return nil
	} else if tokens.Back == tokens.Front {
		tokens.Back = nil
		tokens.Front = nil
	} else {
		tmp := tokens.Back.Prev
		tmp.TruncateEnd()
		tokens.Back = tmp
	}
	return back
}

// Shift a token off the front of the token list
func (tokens *Tokens) Shift() *Token {
	front := tokens.Front
	if tokens.Front == nil {
		return nil
	} else if tokens.Back == tokens.Front {
		tokens.Back = nil
		tokens.Front = nil
	} else {
		tmp := tokens.Front.Next
		tmp.TruncateStart()
		tokens.Front = tmp
	}
	return front
}

// Collapse collapses the tokens between start and end into a single new token and returns it.
func (tokens *Tokens) Collapse(start *Token, end *Token) *Token {
	prev := start.TruncateStart()
	next := end.TruncateEnd()
	token := &Token{}
	token.Children = NewTokens(token)
	token.InsertAfter(next)
	token.InsertBefore(prev)
	token.Children.Push(start)

	// Rebind...
	if tokens.Front == start {
		tokens.Front = token
	}
	if tokens.Back == end {
		tokens.Back = end
	}

	return token
}

// Walk all nodes in this list until out of tokens or the function returns false.
func (tokens *Tokens) Walk(op func(t *Token) bool) {
	for marker := tokens.Front; marker != nil; marker = marker.Next {
		if halt := op(marker); halt {
			break
		}
	}
}

// Debug this group by printing itself and all children
func (tokens *Tokens) Debug(types func(T uint32) string, indent ...int) string {
	indentValue := debugIndent(indent...)
	if tokens.Parent == nil {
		indentValue = 1
	}
	buffer := ""
	count := 0
	for marker := tokens.Front; marker != nil; marker = marker.Next {
		count += 1
		buffer += marker.Debug(types, indentValue)
	}
	if count != 0 {
		tmp := buffer
		buffer = fmt.Sprintf("Tokens (%d found):\n%s", count, debugIndentBlock(tmp, indentValue))
	} else {
		buffer = "Tokens (none found)"
	}
	return buffer
}
