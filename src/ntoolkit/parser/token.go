package parser

import (
	"fmt"
	"strings"
)

// Token is an intrusive linked list of lexical tokens
type Token struct {
	Type     uint32
	Raw      *string
	Parent   *Tokens
	Children *Tokens
	Next     *Token
	Prev     *Token
}

// First returns the first token in the current chain, or the token itself
func (t *Token) First() *Token {
	tmp := t
	for tmp.Prev != nil {
		tmp = tmp.Prev
	}
	return tmp
}

// Last returns the last token in the current chain, or the token itself
func (t *Token) Last() *Token {
	tmp := t
	for tmp.Next != nil {
		tmp = tmp.Next
	}
	return tmp
}

// TruncateStart cuts off any preceding tokens and returns the old prev token
func (t *Token) TruncateStart() *Token {
	rtn := t.Prev
	t.Prev = nil
	if rtn != nil {
		rtn.Next = nil
	}
	return rtn
}

// TruncateEnd cuts off any trailing tokens and returns the old next token
func (t *Token) TruncateEnd() *Token {
	rtn := t.Next
	t.Next = nil
	if rtn != nil {
		rtn.Prev = nil
	}
	return rtn
}

// InsertAfter puts a new token after this token
func (t *Token) InsertAfter(token *Token) {
	if token == nil {
		return
	}

	tmp := t.Next
	t.Next = token

	tmpPrev := token.Prev
	token.Prev = t

	if tmp != nil {
		tmp.Prev = token
	}

	if tmpPrev != nil {
		tmpPrev.Next = nil
	}
}

// InsertBefore puts a new token before this token
func (t *Token) InsertBefore(token *Token) {
	if token == nil {
		return
	}

	tmp := t.Prev
	t.Prev = token

	tmpNext := token.Next
	token.Next = t

	if tmp != nil {
		tmp.Next = token
	}

	if tmpNext != nil {
		tmpNext.Prev = nil
	}
}

// Classify attempts to classify this token and all its child tokens; it returns true if any classify did so
func (t *Token) Classify(classifier Classifier) bool {
	changed := classifier.Classify(t)
	if t.Children != nil {
		t.Children.Walk(func(child *Token) bool {
			changed = child.Classify(classifier) || changed
			return changed
		})
	}
	return changed
}

// Find the next token in this token chain of the given type
func (t *Token) FindNext(T uint32) *Token {
	for marker := t.Next; marker != nil; marker = marker.Next {
		if marker.Type == T {
			return marker
		}
	}
	return nil
}

// Debug this token by printing itself and all children
func (t *Token) Debug(types func(T uint32) string, indent ...int) string {
	buffer := fmt.Sprintf("Token: %s", types(t.Type))
	if t.Raw != nil {
		buffer += fmt.Sprintf(" ('%s')", strings.Replace(*t.Raw, "\n", "\\n", -1))
	}
	buffer += "\n"
	if t.Children != nil {
		buffer += debugIndentBlock(t.Children.Debug(types, debugIndent(indent...)), debugIndent(indent...))
	}
	return buffer
}
