package main

import (
	"fmt"
	"os"
)

func run(tokens []token) int {
	stack := make([]int, 0, len(tokens))
	pop := func() int {
		if len(stack) < 1 {
			panic("too many operators in run")
		}
		n := stack[len(stack)-1]
		stack = stack[0 : len(stack)-1]
		return n
	}
	push := func(n int) {
		stack = append(stack, n)
	}
	for _, t := range tokens {
		switch t.typ {
		case tokenNumber:
			push(t.Number())
		case tokenOperator:
			op := t.Operator()
			n0, n1 := pop(), pop()
			push(op.fn(n1, n0))
		default:
			panic("bad token in run")
		}
	}
	if len(stack) != 1 {
		panic("bad stack count at end of expression")
	}
	return stack[0]
}

func shunt(tokens []token) []token {
	output := make([]token, 0, len(tokens))
	keep := func(t token) {
		output = append(output, t)
	}
	stack := make([]token, 0, len(tokens))
	push := func(t token) {
		stack = append(stack, t)
	}
	shrink := func() {
		stack = stack[0 : len(stack)-1]
	}
	peek := func() token {
		if len(stack) == 0 {
			panic("peek on empty stack")
		}
		return stack[len(stack)-1]
	}
	pop := func() token {
		defer shrink()
		return peek()
	}
	transfer := func() {
		keep(pop())
	}
	handleOperator := func(t token) {
		o1 := getOperator(t)
		for len(stack) > 0 && peek().typ == tokenOperator {
			o2 := getOperator(peek())
			if o1.precedence < o2.precedence ||
				(o1.associativity == left && o1.precedence == o2.precedence) {
				transfer()
			} else {
				break
			}
		}
		push(t)
	}
	for _, t := range tokens {
		switch t.typ {
		case tokenNumber:
			keep(t)
		case tokenOperator:
			handleOperator(t)
		case tokenLeftParen:
			push(t)
		case tokenRightParen:
			for {
				t1 := pop()
				if t1.typ == tokenLeftParen {
					break
				}
				keep(t1)
			}
		default:
			panic("bad token")
		}
	}
	for len(stack) > 0 {
		transfer()
	}
	return output
}

func main() {
	c := make(chan token)
	go lex(os.Stdin, c)
	buf := make([]token, 0, 32)
	for t := range c {
		if t.typ == tokenEnd {
			tokens := make([]token, len(buf))
			copy(tokens, buf)
			fmt.Println(run(shunt(tokens)))
			buf = buf[0:0]
			continue
		}
		buf = append(buf, t)
	}
}
