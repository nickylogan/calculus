package calculus

// LeftParen is a left parenthesis
type LeftParen struct{}

func (l LeftParen) String() string {
	return "("
}

// RightParen is a right parentheses
type RightParen struct{}

func (r RightParen) String() string {
	return ")"
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
