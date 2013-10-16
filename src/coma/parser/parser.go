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
  "coma/ast"
  "coma/lexer"
  "coma/rules"
  "coma/token"
  "fmt"
  "io"
)

type Parser struct {
  lexer   *lexer.Lexer
  rules   []rules.IRule
  context Context

  // Events
  BeforeParse func(p *Parser, rule rules.IRule) int
  AfterParse  func(p *Parser, rule rules.IRule) int
  OnError     func(p *Parser, err string)
}

/** ***************************************************************************
 * Parser declaration
 ******************************************************************************/

func MakeParser(l *lexer.Lexer, r []rules.IRule) *Parser {
  parser := new(Parser)
  parser.init(l, r)
  return parser
}

func MakeParserByOptions(opt []interface{}) *Parser {
  parser := new(Parser)
  parser.initByOptions(opt)
  return parser
}

func (p *Parser) init(l *lexer.Lexer, r []rules.IRule) {
  p.lexer = l
  if nil != r {
    p.rules = make([]rules.IRule, len(r))
    copy(p.rules, r)
  }
}

func (p *Parser) initByOptions(opt []interface{}) {
  p.lexer = lexer.MakeLexerByOptions(opt)
  if nil != opt {
    for _, it := range opt {
      switch it.(type) {
      case rules.IRule:
        p.AddRule(it.(rules.IRule))
      default:
        // Somothing else...
      }
    }
  }
}

/** ***************************************************************************
 * Lexer
 ******************************************************************************/

func (p *Parser) AddLex(lex lexer.LexBase) *Parser {
  p.lexer.AddLex(lex)
  return p
}

func (p *Parser) AddStdLex(lex int, lit string, prepare int) *Parser {
  p.lexer.AddStdLex(lex, lit, prepare)
  return p
}

func (p *Parser) AddCustomLex(lex int, test func(*lexer.Lexer, string) bool) *Parser {
  p.lexer.AddCustomLex(lex, test)
  return p
}

func (p *Parser) GetLexer() *lexer.Lexer {
  return p.lexer
}

/** ***************************************************************************
 * Actions
 ******************************************************************************/

func (p *Parser) SetStream(s io.ReadSeeker) *Parser {
  p.lexer.SetStream(s)
  return p
}

func (p *Parser) AddRule(r rules.IRule) *Parser {
  if nil == p.rules {
    p.rules = make([]rules.IRule, 1)
    p.rules[0] = r
  } else {
    p.rules = append(p.rules, r)
  }
  return p
}

/** ***************************************************************************
 * Parser methods
 ******************************************************************************/

func (p *Parser) normalize(astree *ast.Node) {
  // Normalize AST tree
  for j := 0; j < len(p.rules); j++ {
    r := p.rules[j]
    c := astree.Child.RollUp(r.IsRight(), func(cur *ast.Node) (int /* count */, interface{} /* rule */) {
      count := r.Test(cur, p.rules)
      // fmt.Println("count", count, cur.Token)
      if count >= 0 {
        return count, r
      }
      return 0, nil
    })
    if c > 0 && j > 0 {
      // Start agane
      j = -1
    }
  }
}

func (p *Parser) checkAst(astree *ast.Node) bool {
  // Check first level ast
  // for _, n := range ast.Childrens {
  // 	if nil == m.Rule {
  // 		return false
  // 	}
  // }
  return true
}

func (p *Parser) Parse(f func(*ast.Node) bool) bool {
  astree := &ast.Node{}

  // Construct AST
  for {
    line, simbol := 0, 0
    tokType, _, word := p.lexer.NextToken()
    // if token.Separator == tokType {
    // 	continue
    // }

    // fmt.Printf(">>> %d : \"%s\"\n", tokType, word)
    if tokType == token.Error || tokType == token.EOF {
      break
    }

    // Create new tocken
    tok := &token.Token{Code: tokType, Value: word, Line: line, Simbol: simbol}

    if nil == astree.AddChild(&ast.Node{Token: tok}) {
      // TODO Error...
      fmt.Println("ERROOOOOOOOOOOOOOOOOOOOOOOR!")
      return false
    }
  }

  // Normalize AST tree
  p.normalize(astree)

  // if nil != astree.Child {
  //   astree.Child.Clean()
  // }
  // fmt.Println("#@#@#", astree.ToTreeString(0, 4))

  // Check first level AST
  if p.checkAst(astree) {
    // Process ast
    return astree.Child.Each(f)
  }
  return false
}
