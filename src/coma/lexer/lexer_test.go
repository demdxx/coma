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
	"io"
	"strings"
	"testing"
	"coma/token"
)

var content = `
x = 100 + 900
a = -678
c = x + a * 2
`

var options = []interface{}{
	&Option{OptType: Option_SeparateAlpha, Value: ".,:;`~!?@#$%%^&*-+=_/\\|(){}[]<>\"'"},
	&Option{OptType: Option_AbsSeparateAlpha, Value: " \t\r\n"},

	&Lex{code: token.KeyWord, lit: "put", prepare: PrepareAs_Literal},
	&Lex{code: token.Operator, lit: "+", prepare: PrepareAs_Literal},
	&Lex{code: token.Digit, lit: "0123456789", prepare: PrepareAs_Alpha},
	&Lex{code: token.Word, lit: "abcdefx", prepare: PrepareAs_Alpha},
	&Lex{code: 0, lit: "/* comment */", prepare: PrepareAs_Special},
}

func TestLexer(t *testing.T) {
	var reader io.ReadSeeker = strings.NewReader(content)
	lexer := MakeLexerByOptions(options)
	lexer.SetStream(reader)

	for {
		tt, _, w := lexer.NextToken()
		if token.Separator != tt {
			t.Logf("%d : \"%s\"", tt, w)
			if tt == token.Error || tt == token.EOF {
				break
			}
		}
	}
}
