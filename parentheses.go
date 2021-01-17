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

type depthStack struct {
	stack [][2]int // stores pairs of index and depth
}

func (d *depthStack) increment(idx int) {
	d.stack = append(d.stack, [2]int{idx, d.current() + 1})
}

func (d *depthStack) decrement(idx int) {
	d.stack = append(d.stack, [2]int{idx, d.current() - 1})
}

func (d *depthStack) current() int {
	if len(d.stack) == 0 {
		return 0
	}
	return d.stack[len(d.stack)-1][1]
}

func (d *depthStack) clear() {
	d.stack = nil
}
