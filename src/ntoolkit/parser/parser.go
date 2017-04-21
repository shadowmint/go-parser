package parser

import (
	"ntoolkit/futures"
	"ntoolkit/errors"
)

// Parse returns a deferred promise for a parsed token stream.
func Parse(tokenizer Tokenizer, classifier Classifier, stream func(func(data string, finished bool))) *DeferredTokens {
	rtn := &DeferredTokens{}
	root := NewTokens(nil)
	parseTokens(root, tokenizer, stream).Then(func() {
		maxIterations := (1 + root.Count()) * 10
		if err := classifyTokens(root, classifier, maxIterations); err != nil {
			rtn.Reject(err)
		} else {
			rtn.Resolve(root)
		}
	}, func(err error) {
		rtn.Reject(err)
	})
	return rtn
}

// Split all incoming data into strings
func parseTokens(tokens *Tokens, tokenizer Tokenizer, stream func(func(data string, finished bool))) *futures.Deferred {
	rtn := &futures.Deferred{}
	tokenizer.Enter(tokens)
	go stream(func(data string, finished bool) {
		if len(data) > 0 {
			tokenizer.Process(data)
		}
		if finished {
			tokenizer.Close()
			rtn.Resolve()
		}
	})
	return rtn
}

// Classify all tokens running at most maxInterations
func classifyTokens(tokens *Tokens, classifier Classifier, maxInterations int) error {
	resolved := false
	for i := 0; i < maxInterations; i++ {
		changed := false
		tokens.Walk(func(t *Token) bool {
			changed = t.Classify(classifier) || changed
			return changed
		})
		if !changed {
			resolved = true
			break
		}
	}
	if !resolved {
		return errors.Fail(ErrMaxClassifyIterations{}, nil, "Unable to resolve tokens after maximum iteration count")
	}
	return nil
}
