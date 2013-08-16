package main

import (
	"fmt"
    "errors"
	"os"
)

func run(tokens []token) (int, error) {
	stack := make([]int, 0, len(tokens))
	pop := func() int {
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
            if len(stack) < 2 {
                return 0, errors.New("expression has extra operators")
            }
			n0, n1 := pop(), pop()
			push(op.fn(n1, n0))
		default:
			return 0, errors.New("expression has unexpected token")
		}
	}
	if len(stack) != 1 {
        return 0, errors.New("expression has extra operands")
	}
	return stack[0], nil
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
        switch t.typ {
        case tokenEnd:
			tokens := make([]token, len(buf))
			copy(tokens, buf)
            n, err := run(shunt(tokens))
            if err != nil {
                fmt.Fprintln(os.Stderr, err.Error())
            } else {
                fmt.Println(n)
            }
			buf = buf[0:0]
			continue
        case tokenError:
            fmt.Println(t)
            buf = buf[0:0]
        default:
            buf = append(buf, t)
		}
	}
}
