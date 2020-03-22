package calculus

// Token represents a single expression token
type Token interface {
	String() string
}

// Tokenizer is implemented by an expression tokenizer.
type Tokenizer interface {
	// Tokenize is called when an expression needs to be split into understandable tokens.
	Tokenize(expr string) (tokens []Token, err error)
}

type tokenizer struct {
	// TODO: add fields
}

// NewTokenizer creates a new Tokenizer
func NewTokenizer() Tokenizer {
	return &tokenizer{}
}

// Tokenize implements the Tokenizer interface
func (t *tokenizer) Tokenize(expr string) (tokens []Token, err error) {

	return
}
