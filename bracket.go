package yamp

// Bracket represents bracket tokens such as parentheses.
type Bracket interface {
	Token

	IsLeft() bool
	IsRight() bool
}

type bracket int

var _ Bracket = bracket(0)

const (
	// LeftParen represents a left parenthesis.
	LeftParen bracket = 1
	// RightParen represents a right parenthesis.
	RightParen bracket = 2
)

func (b bracket) String() string {
	switch b {
	case LeftParen:
		return "("
	case RightParen:
		return ")"
	default:
		return ""
	}
}

// IsLeft implements Bracket.
func (b bracket) IsLeft() bool {
	switch b {
	case LeftParen:
		return true
	default:
		return false
	}
}

// IsRight implements Bracket.
func (b bracket) IsRight() bool {
	switch b {
	case RightParen:
		return true
	default:
		return false
	}
}

type bracketDepth struct {
	index int
	depth int
	b     Bracket
}

type bracketStack struct {
	stack []bracketDepth
}

func (s *bracketStack) increment(idx int, b Bracket) {
	s.stack = append(s.stack, bracketDepth{
		index: idx,
		depth: s.depth() + 1,
		b:     b,
	})
}

func (s *bracketStack) decrement(idx int, b Bracket) {
	s.stack = append(s.stack, bracketDepth{
		index: idx,
		depth: s.depth() - 1,
		b:     b,
	})
}

func (s bracketStack) depth() int {
	if len(s.stack) == 0 {
		return 0
	}
	return s.stack[len(s.stack)-1].depth
}

func (s *bracketStack) clear() {
	s.stack = nil
}
