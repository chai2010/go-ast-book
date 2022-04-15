package main

type Stack struct {
	stk []int
}

func (p *Stack) push(i int) {
	p.stk = append(p.stk, i)
}
func (p *Stack) pop() int {
	if len(p.stk) < 1 {
		panic("stk is empty unable to pop")
	}
	result := p.stk[len(p.stk)-1]
	p.stk = p.stk[:len(p.stk)-1]
	return result
}
