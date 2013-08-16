package main

import (
	"math"
)

type associativity int

const (
	left associativity = iota
	right
)

type operator struct {
	associativity
	precedence int
	fn         func(int, int) int
}

var operators = map[string]operator{
	"+": {left, 2, add},
	"-": {left, 2, sub},
	"*": {left, 3, mul},
	"/": {left, 3, div},
	"^": {right, 4, pow},
}

func getOperator(t token) operator {
	op, ok := operators[t.val]
	if !ok {
		panic("bad operator " + t.val)
	}
	return op
}

func add(x, y int) int {
	return x + y
}

func sub(x, y int) int {
	return x - y
}

func mul(x, y int) int {
	return x * y
}

func div(x, y int) int {
	return x / y
}

func pow(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}
