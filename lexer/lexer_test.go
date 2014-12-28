// Copyright (c) 2014 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lexer

import (
	"strings"
	"testing"
	"text/scanner"
)

type lexerTest struct {
	name   string
	input  string
	tokens []Token
}

// All the lexer tests.
var lexerTests = []lexerTest{
	{"empty", "", []Token{newToken(EOF, "", 0, 0)}},
	{"identifiers", "a foo _a _1 a1 ifelse", []Token{
		newToken(Ident, "a", 1, 1),
		newToken(Ident, "foo", 1, 3),
		newToken(Ident, "_a", 1, 7),
		newToken(Ident, "_1", 1, 10),
		newToken(Ident, "a1", 1, 13),
		newToken(Ident, "ifelse", 1, 16),
		newToken(EOF, "", 1, 22)}},
	{"literals", "42 35 true false", []Token{
		newToken(Number, "42", 1, 1),
		newToken(Number, "35", 1, 4),
		newToken(True, "true", 1, 7),
		newToken(False, "false", 1, 12),
		newToken(EOF, "", 1, 17)}},
	{"keywords", "bool else false if int true var for", []Token{
		newToken(Bool, "bool", 1, 1),
		newToken(Else, "else", 1, 6),
		newToken(False, "false", 1, 11),
		newToken(If, "if", 1, 17),
		newToken(Int, "int", 1, 20),
		newToken(True, "true", 1, 24),
		newToken(Var, "var", 1, 29),
		newToken(For, "for", 1, 33),
		newToken(EOF, "", 1, 36)}},
	{"operators", "= * / + - < <= == != >= > && || !", []Token{
		newToken(Assign, "=", 1, 1),
		newToken(Multiply, "*", 1, 3),
		newToken(Divide, "/", 1, 5),
		newToken(Plus, "+", 1, 7),
		newToken(Minus, "-", 1, 9),
		newToken(Less, "<", 1, 11),
		newToken(LessOrEqual, "<=", 1, 13),
		newToken(Equal, "==", 1, 16),
		newToken(NotEqual, "!=", 1, 19),
		newToken(GreaterOrEqual, ">=", 1, 22),
		newToken(Greater, ">", 1, 25),
		newToken(And, "&&", 1, 27),
		newToken(Or, "||", 1, 30),
		newToken(Not, "!", 1, 33),
		newToken(EOF, "", 1, 34)}},
	{"delimiters", "{}()", []Token{
		newToken(LeftBrace, "{", 1, 1),
		newToken(RightBrace, "}", 1, 2),
		newToken(LeftParen, "(", 1, 3),
		newToken(RightParen, ")", 1, 4),
		newToken(EOF, "", 1, 5)}},
	{"line comment", "a // comment\nb", []Token{
		newToken(Ident, "a", 1, 1),
		newToken(Ident, "b", 2, 1),
		newToken(EOF, "", 2, 2)}},
	{"block comment", "a /* x \n x */ b c/* x */d", []Token{
		newToken(Ident, "a", 1, 1),
		newToken(Ident, "b", 2, 7),
		newToken(Ident, "c", 2, 9),
		newToken(Ident, "d", 2, 17),
		newToken(EOF, "", 2, 18)}},
	{"errors", "& | } ) ?", []Token{
		newToken(Error, "expected && operator", 1, 1),
		newToken(Error, "expected || operator", 1, 3),
		newToken(Error, "unexpected }", 1, 5),
		newToken(Error, "unexpected )", 1, 7),
		newToken(Error, "unrecognized token ?", 1, 9),
		newToken(EOF, "", 1, 10)}},
}

// TestLex runs all lexer tests.
func TestLex(t *testing.T) {
	for _, test := range lexerTests {
		tokens := collect(&test)
		if !equal(tokens, test.tokens) {
			t.Errorf("%s: got '%v' expected '%v'", test.name, tokens, test.tokens)
		}
	}
}

// collect scans all tokens of a lexer test and puts them into a slice.
func collect(test *lexerTest) (tokens []Token) {
	l := Lex(test.name, strings.NewReader(test.input))
	for t := range l.Tokens {
		tokens = append(tokens, t)
		if t.Type == EOF {
			break
		}
	}
	return
}

// equal checks whether two slices of tokens are the same.
func equal(t1, t2 []Token) bool {
	if len(t1) != len(t2) {
		return false
	}
	for i := range t1 {
		if t1[i].Type != t2[i].Type || t1[i].Text != t2[i].Text {
			return false
		}
		if t1[i].Pos.Line != t2[i].Pos.Line || t1[i].Pos.Column != t2[i].Pos.Column {
			return false
		}
	}
	return true
}

// newToken creates a new tokens.
func newToken(typ Type, text string, line, column int) Token {
	pos := scanner.Position{}
	pos.Line = line
	pos.Column = column
	return Token{pos, text, typ}
}
