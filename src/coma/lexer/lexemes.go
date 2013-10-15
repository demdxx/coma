// Project coma
//
// Copyright (C) <2013> Dmitry Ponomarev <demdxx@gmail.com>
//
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the Software
// is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
// PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
// SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package lexer

import (
	"regexp"
	"strings"
)

const (
	PrepareAs_Special = iota
	PrepareAs_Literal
	PrepareAs_Alpha
	PrepareAs_RegExp
)

/** ***************************************************************************
 * Structures
 ******************************************************************************/

type LexBase interface {
	test(l *Lexer, s string) bool
	GetCode() int
}

type LexCustom struct {
	LexBase
	code    	int
	testFunc 	func(*Lexer, string) bool
}

type Lex struct {
	LexBase
	code    int
	lit     string
	prepare int
}

func MakeLex(t int, l string, p int) *Lex {
	return &Lex{code: t, lit: l, prepare: p}
}

func MakeLexCustom(t int, test func(*Lexer, string) bool) *LexCustom {
	return &LexCustom{code: t, testFunc: test}
}

/** ***************************************************************************
 * Lexem Methods
 ******************************************************************************/

func (self Lex) test(l *Lexer, s string) bool {
	switch self.prepare {
	case PrepareAs_Special:
		return self.testSpecial(l, s)
	case PrepareAs_Literal:
		return self.testLiteral(l, s)
	case PrepareAs_Alpha:
		return self.testAlpha(l, s)
	case PrepareAs_RegExp:
		return self.testRegexp(l, s)
	}
	return false
}

func (self *Lex) GetCode() int {
	return self.code
}

func (self *Lex) testSpecial(l *Lexer, word string) bool {
	// fmt.Println("testSpecial", word, self.lit)
	return false
}

func (self *Lex) testLiteral(l *Lexer, word string) bool {
	// fmt.Println("testLiteral", word, self.lit)
	return self.lit == word
}

func (self *Lex) testAlpha(l *Lexer, word string) bool {
	// fmt.Println("testAlpha", word, self.lit)
	for _, c := range word {
		if -1 == strings.Index(self.lit, string(c)) {
			return false
		}
	}
	return true
}

func (self *Lex) testRegexp(l *Lexer, word string) bool {
	// fmt.Println("testRegexp", word, self.lit)
	matched, _ := regexp.MatchString(self.lit, word)
	return matched
}

/** ***************************************************************************
 * Custom lexem Methods
 ******************************************************************************/

func (self LexCustom) GetCode() int {
	return self.code
}

func (self LexCustom) test(l *Lexer, s string) bool {
	return self.testFunc(l, s)
}
