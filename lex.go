package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const eof = -1

var (
	digits = "0123456789"
)

type stateFn func(*lexer) stateFn

type lexer struct {
	in  *bufio.Reader
	out chan token // channel of scanned tokens
	buf []rune
}

func newLexer(in io.Reader, out chan token) *lexer {
	l := &lexer{
		in:  bufio.NewReader(in),
		out: out,
		buf: make([]rune, 0, 32),
	}
	return l
}

func (l *lexer) next() rune {
	r, _, err := l.in.ReadRune()
	switch err {
	case nil:
	case io.EOF:
		return eof
	default:
		l.errorf("lex error in next: %v", err)
		return eof
	}
	l.buf = append(l.buf, r)
	return r
}

func (l *lexer) discard() {
	if len(l.buf) >= 1 {
		l.buf = l.buf[0 : len(l.buf)-1]
	}
}

func (l *lexer) backup() {
	l.discard()
	err := l.in.UnreadRune()
	if err != nil {
		l.errorf("lex error in backup: %v", err)
	}
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) clearBuf() {
	l.buf = l.buf[0:0]
}

func (l *lexer) emit(t tokenType) {
	l.out <- token{string(l.buf), t}
	l.clearBuf()
}

func (l *lexer) emitError(format string, args ...interface{}) {
	l.out <- token{fmt.Sprintf(format, args...), tokenError}
}

func (l *lexer) fatalf(format string, args ...interface{}) stateFn {
	l.emitError(format, args...)
	return nil
}

func (l *lexer) skipUntil(good string) {
	for {
		r := l.next()
		if r == eof || strings.ContainsRune(good, r) {
			l.backup()
			return
		}
	}
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.emitError(format, args...)
	l.skipUntil("\n\r")
	l.next()
	l.clearBuf()
	return lexWhitespace
}

func isDigit(r rune) bool {
	return strings.ContainsRune(digits, r)
}

func isOperator(r rune) bool {
	_, ok := operators[string(r)]
	return ok
}

func isTerminator(r rune) bool {
	switch r {
	case '\n', '\r':
		return true
	}
	return false
}

func isWhitespace(r rune) bool {
	switch r {
	case ' ', '\t':
		return true
	}
	return false
}

func lexWhitespace(l *lexer) stateFn {
	r := l.next()
	switch {
	case isDigit(r):
		return lexNum
	case isOperator(r):
		l.backup()
		return lexOperator
	case isWhitespace(r):
		l.discard()
		return lexWhitespace
	case isTerminator(r):
		l.emit(tokenEnd)
		return lexWhitespace
	case r == '(':
		l.backup()
		return lexLeftParen
	case r == ')':
		l.backup()
		return lexRightParen
	case r == eof:
		return nil
	}
	return l.errorf("illegal rune in lexRoot: %c", r)
}

func lexNum(l *lexer) stateFn {
	r := l.next()
	switch {
	case isDigit(r):
		return lexNum
	case isOperator(r):
		l.backup()
		l.emit(tokenNumber)
		return lexOperator
	case isWhitespace(r):
		l.discard()
		l.emit(tokenNumber)
		return lexWhitespace
	case isTerminator(r):
		l.backup()
		l.emit(tokenNumber)
		l.next()
		l.emit(tokenEnd)
		return lexWhitespace
	case r == '(':
		l.backup()
		l.emit(tokenNumber)
		return lexLeftParen
	case r == ')':
		l.backup()
		l.emit(tokenNumber)
		return lexRightParen
	case r == eof:
		return nil
	}
	return l.errorf("illegal rune in lexNum: %c", r)
}

func lexOperator(l *lexer) stateFn {
	r := l.next()
	switch {
	case isOperator(r):
		l.emit(tokenOperator)
		return lexWhitespace
	case r == eof:
		return nil
	default:
		return l.errorf("illegal rune in lexOperator: %c", r)
	}
}

func lexLeftParen(l *lexer) stateFn {
	switch r := l.next(); r {
	case '(':
		l.emit(tokenLeftParen)
		return lexWhitespace
	default:
		return l.errorf("illegal rune in lexLeftParen: %c", r)
	}
}

func lexRightParen(l *lexer) stateFn {
	switch r := l.next(); r {
	case ')':
		l.emit(tokenRightParen)
		return lexWhitespace
	default:
		return l.errorf("illegal rune in lexRightParen: %c", r)
	}
}

func lex(in io.Reader, c chan token) {
	defer close(c)
	lexer := newLexer(in, c)
	for fn := lexWhitespace; fn != nil; {
		fn = fn(lexer)
	}
}
