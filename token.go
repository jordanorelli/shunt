package main

import (
	"fmt"
	"strconv"
)

type tokenType int

const (
	tokenError tokenType = iota
	tokenEOF
	tokenNumber
	tokenLeftParen
	tokenRightParen
	tokenOperator
	tokenEnd
)

type token struct {
	val string
	typ tokenType
}

func (t token) String() string {
	switch t.typ {
	case tokenEOF:
		return "EOF"
	case tokenError:
		return t.val
	}
	if len(t.val) > 10 {
		return fmt.Sprintf("%.10q...", t.val)
	}
	return fmt.Sprintf("%q", t.val)
}

func (t token) Number() int {
	if t.typ != tokenNumber {
		panic("can't get number of non-number token")
	}
	i, err := strconv.ParseInt(t.val, 10, 64)
	if err != nil {
		panic(err.Error())
	}
	return int(i)
}

func (t token) Operator() operator {
	if t.typ != tokenOperator {
		panic("can't get operator of non-operator token")
	}
	return getOperator(t)
}
