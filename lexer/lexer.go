// Copyright (c) 2014 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate stringer -type Type

package lexer

import (
	"fmt"
	"io"
	"text/scanner"
)

// Type categorizes a token.
type Type int

const (
	EOF   Type = iota // end of file
	Error             // error, value is the text of the lexeme
	Ident             // alphanumeric identifier
	// Literals
	False  // false
	Number // integer number
	True   // true
	// Keywords
	Else  // else
	If    // if
	Var   // var
	While // while
	// Types
	Bool // bool
	Int  // int
	// Assignment operator
	Assign // =
	// Arithmetic operators
	Multiply // *
	Divide   // /
	Plus     // +
	Minus    // -
	// Relationship operators
	Less           // <
	LessOrEqual    // <=
	Greater        // >=
	GreaterOrEqual // >
	Equal          // ==
	NotEqual       // !=
	// Logical operators
	Not // !
	And // &&
	Or  // ||
	// Delimiters
	LeftParen  // (
	RightParen // )
	LeftBrace  // {
	RightBrace // }
)

// All the keywords.
var keywords = map[string]Type{
	"bool":  Bool,
	"else":  Else,
	"false": False,
	"if":    If,
	"int":   Int,
	"true":  True,
	"var":   Var,
	"while": While,
}

// Token represents a token.
type Token struct {
	Pos  scanner.Position // position in the input string
	Text string           // text of this token
	Type Type             // tpye of this token
}

func (t Token) String() string {
	switch t.Type {
	case EOF:
		return "EOF"
	case Error:
		return t.Text
	}
	if len(t.Text) > 10 {
		return fmt.Sprintf("%s:%d:%d %.10q...", t.Pos.Filename, t.Pos.Line, t.Pos.Column, t.Text)
	}
	return fmt.Sprintf("%s:%d:%d %q", t.Pos.Filename, t.Pos.Line, t.Pos.Column, t.Text)
}

// Lexer represents the lexical analyser.
type Lexer struct {
	Tokens  chan Token
	scanner scanner.Scanner
	braces  int // nesting of braces
	parens  int // nesting of parens
}

// Lex creates a new lexer for the input source.
func Lex(filename string, src io.Reader) *Lexer {
	l := &Lexer{Tokens: make(chan Token)}
	l.scanner.Init(src)
	l.scanner.Filename = filename
	go l.run()
	return l
}

// emit emits a token to the channel of tokens.
func (l *Lexer) emit(pos scanner.Position, text string, typ Type) {
	l.Tokens <- Token{pos, text, typ}
}

// emitHere emits the current token to the channel of tokens.
func (l *Lexer) emitHere(typ Type) {
	l.Tokens <- Token{l.scanner.Position, l.scanner.TokenText(), typ}
}

// emitIfNext emits a token of type t1 if r matches the next rune. Otherwise a token of type t2 is emitted.
func (l *Lexer) emitIfNext(r rune, t1, t2 Type) {
	if l.scanner.Peek() == r {
		pos := l.scanner.Position
		text := l.scanner.TokenText()
		l.scanner.Scan()
		text += l.scanner.TokenText()
		l.emit(pos, text, t1)
	} else {
		l.emit(l.scanner.Position, l.scanner.TokenText(), t2)
	}
}

// error emits an error token.
func (l *Lexer) errorf(format string, args ...interface{}) {
	l.emit(l.scanner.Position, fmt.Sprintf(format, args...), Error)
}

// expect emits an token if r matches the next rune. Otherwise an error is emitted.
func (l *Lexer) expect(r rune, typ Type, err string) {
	if l.scanner.Peek() == r {
		pos := l.scanner.Position
		text := l.scanner.TokenText()
		l.scanner.Scan()
		text += l.scanner.TokenText()
		l.emit(pos, text, typ)
	} else {
		l.errorf(err)
	}
}

// run starts the lexer.
func (l *Lexer) run() {
	for state := lexSource; state != nil; {
		state = state(l)
	}
	close(l.Tokens)
}

// stateFn is a state of the lexer. It is a function that returns the next state.
type stateFn func(*Lexer) stateFn

// lexSource scans the source.
func lexSource(l *Lexer) stateFn {
	switch l.scanner.Scan() {
	case scanner.EOF:
		l.emitHere(EOF)
		return nil
	case scanner.Ident:
		return lexIdent
	case scanner.Int:
		l.emitHere(Number)
	case '*':
		l.emitHere(Multiply)
	case '/':
		l.emitHere(Divide)
	case '+':
		l.emitHere(Plus)
	case '-':
		l.emitHere(Minus)
	case '=':
		l.emitIfNext('=', Equal, Assign)
	case '<':
		l.emitIfNext('=', LessOrEqual, Less)
	case '>':
		l.emitIfNext('=', GreaterOrEqual, Greater)
	case '!':
		l.emitIfNext('=', NotEqual, Not)
	case '&':
		l.expect('&', And, "expected && operator")
	case '|':
		l.expect('|', Or, "expected || operator")
	case '{':
		return lexLeftBrace
	case '}':
		return lexRightBrace
	case '(':
		return lexLeftParen
	case ')':
		return lexRightParen
	default:
		l.errorf("unrecognized token %v", l.scanner.TokenText())
	}
	return lexSource
}

// lexIdent scans an alphanumeric identifier.
func lexIdent(l *Lexer) stateFn {
	typ, ok := keywords[l.scanner.TokenText()]
	if !ok {
		typ = Ident
	}
	l.emitHere(typ)
	return lexSource
}

// lexLeftBrace scans a left brace and keeps track of the nesting.
func lexLeftBrace(l *Lexer) stateFn {
	l.braces++
	l.emitHere(LeftBrace)
	return lexSource
}

// lexRightBrace scans a right brace and keeps track of the nesting.
func lexRightBrace(l *Lexer) stateFn {
	l.braces--
	if l.braces < 0 {
		l.errorf("unexpected }")
	} else {
		l.emitHere(RightBrace)
	}
	return lexSource
}

// lexLeftParen scans a left paren and keeps track of the nesting.
func lexLeftParen(l *Lexer) stateFn {
	l.parens++
	l.emitHere(LeftParen)
	return lexSource
}

// lexRightParen scans a right paren and keeps track of the nesting.
func lexRightParen(l *Lexer) stateFn {
	l.parens--
	if l.parens < 0 {
		l.errorf("unexpected )")
	} else {
		l.emitHere(RightParen)
	}
	return lexSource
}
