package parser

const TokenTypeNone = 0

// Tokenizer receives an arbitrary stream of text data in chunks and converts it into tokens.
// Typical usage would be: t := ... ; defer (func() { t.Close() }) ; t.Enter(context) ; for ... { t.Process(data) }
type Tokenizer interface {
	// Enter starts processing a data stream with a context object and resets the internal tokenizer state
	Enter(context *Tokens)

	// Process reads incoming string data and pushes tokens to the context
	Process(data string)

	// Close should close and resolve the tokenizer.
	Close()
}

// Classifier processes a context and reassigns token types based on contextual information.
type Classifier interface {
	// Classify reassigns token types contextually. Classify should return true if any tokens were modified.
	// Classify will be run multiple times until no reassignment is performed.
	Classify(token *Token) bool
}
