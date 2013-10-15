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

package parser

import (
  "coma/lexer"
  "coma/rules"
  "coma/token"
  "coma/ast"
  "io"
  "strings"
  "testing"
  "fmt"
)

var content = `
x = 100 + 900
a = -678
c = x + a * 2
put x, a, c
`

var options = []interface{}{
  // Lexer options
  lexer.Option{lexer.Option_SeparateAlpha, ".,:;`~!?@#$%%^&*-+=_/\\|(){}[]<>\"' \t\n"},
  lexer.Option{lexer.Option_AbsSeparateAlpha, " \t\r"},

  lexer.MakeLex(token.KeyWord, "put", lexer.PrepareAs_Literal),
  lexer.MakeLex(token.Digit, "0123456789", lexer.PrepareAs_Alpha),
  lexer.MakeLex(token.Operator, "\n", lexer.PrepareAs_Literal),
  lexer.MakeLex(token.Word, "abcdefghijklmnopqrstuvwxyz", lexer.PrepareAs_Alpha),

  // Parser options
  rules.MakeRule(-1, "put", "^put$", nil),
  rules.MakeRule(0, "numb", "^[0-9]+$", nil),
  rules.MakeRule(1, "name", "^[a-z_$][a-z0-9_$]*$", nil),
  rules.MakeRule(2, "exp", "numb|name", nil),
  rules.MakeRule(3, "exp", "'(' exp ')'", nil),
  rules.MakeRule(6, "exp", "exp '*' exp", nil),
  rules.MakeRule(7, "exp", "exp '/' exp", nil),
  rules.MakeRule(8, "exp", "exp '%%' exp", nil),
  rules.MakeRule(4, "exp", "exp '+' exp", nil),
  rules.MakeRule(5, "exp", "exp '-' exp", nil),
  rules.MakeRule(9, "exp", "'-'exp", rules.PostProcessConcat),
  rules.MakeRule(10, "exp", "exp ',' exp", nil),
  rules.MakeRule(11, "exp", "exp '=' exp", nil),
  rules.MakeRule(11, "exp", "put exp", nil),
  rules.MakeRightRule(11, "exp", "exp '\n'", nil),
}

// func beforeParseRuleCallback(p *Parser, rule rules.IRule) int {
//   return 0
// }

// func afterParseRuleCallback(p *Parser, rule rules.IRule) int {
//   return 0
// }

// func parserErrorCallback(p *Parser, err string) {

// }

func TestParser(t *testing.T) {
  var reader io.ReadSeeker = strings.NewReader(content)

  parser := MakeParserByOptions(options)
  parser.SetStream(reader)

  // Set callbacks
  // parser.BeforeParse = beforeParseRuleCallback
  // parser.AfterParse = afterParseRuleCallback
  // parser.OnError = parserErrorCallback

  // Run parser
  parser.Parse(func(n *ast.Node) bool {
    fmt.Println(n.ToString(true))
    return true
  })
}
