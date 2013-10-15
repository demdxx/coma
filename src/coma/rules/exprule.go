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

package rules

import (
  "coma/ast"
  "coma/token"
  "fmt"
  "strings"
)

/** ***************************************************************************
 * Rule Expression tokens
 ******************************************************************************/

type tokenType int

const (
  token_string tokenType = iota
  token_sep_optional
  token_sep
  token_anything
  token_rule
  token_super
)

type iptoken interface {
  addToken(t *ptoken)
}

type ptoken struct {
  iptoken
  vtype    tokenType
  value    string
  tokens   []*ptoken
  optional bool
  repeat   bool
  isnot    bool
}

/** ***************************************************************************
 * Rule Expression
 ******************************************************************************/

type ExpRule struct {
  Rule
  iptoken
  Expression string // example: exp '+' exp
  tokens     []*ptoken
}

func (rule *ExpRule) init() {
  // Parse expression
  rule.prepare()
}

func MakeExpRule(code int, name string, expression string, right bool, fn postProcessFunc) *ExpRule {
  rule := &ExpRule{
    Rule: Rule{
      NumericCode: code,
      Name: name,
      right: right,
      postPrecess: fn,
    },
    Expression: expression,
  }
  rule.init()
  return rule
}

func (rule *ExpRule) Test(node *ast.Node, rules []IRule) int {
  // fmt.Println("+++++++++++++++++> [", rule.Expression, "]", node.ToString(true))
  // r := node
  // count := 0
  // ln := ""
  // // Parse by tokens
  // for _, t := range rule.tokens {
  //   if nil == r {
  //     return -1
  //   }
  //   next := true
  //   switch t.vtype {
  //   case token_string:
  //     if nil == r.Token || t.value != r.Token.Value {
  //       return -1
  //     }
  //   case token_sep_optional:
  //     if nil != r.Token && token.Separator != r.Token.Code {
  //       next = false
  //     }
  //   case token_sep:
  //     if nil == r.Token || token.Separator != r.Token.Code {
  //       return -1
  //     }
  //   case token_rule:
  //     if !NodeTestRuleName(r, t.value) {
  //       return -1
  //     }
  //   }
  //   ln += r.ToString(true)
  //   if next {
  //     r = r.Next
  //     count++
  //   } else {
  //     next = true
  //   }
  // }
  // fmt.Println("-------------> PR", rule.Expression, node.Token, count, "`", ln, "`")
  return tests(rule, rule.tokens, node, rules)
}

func (rule *ExpRule) ToString(detail bool) string {
  result := rule.Name + " : " + rule.Expression
  if detail {
    result += " DETAIL"
    if nil != rule.tokens {
      for _, it := range rule.tokens {
        result += "\n" + it.toTreeString(0, 4)
      }
    } else {
      result += " => <error>"
    }
  }
  return result
}

func (t *ptoken) toString() string {
  result := ""
  if t.isnot {
    result += "!"
  }
  if t.optional {
    result += "["
  }
  if "\n" == t.value {
    result += "\\n"
  } else {
    result += t.value
  }

  if nil != t.tokens {
    for _, it := range t.tokens {
      result += it.toString()
    }
  }

  if t.optional {
    result += "]"
  }
  if t.repeat {
    if t.optional {
      result += "*"
    } else {
      result += "+"
    }
  }
  return result
}

func (t *ptoken) toTreeString(spaces int, step int) string {
  result := ""
  for i := 0; i < spaces; i++ {
    result += " "
  }
  if t.isnot {
    result += "!"
  }
  if t.optional {
    result += "["
  }
  if "\n" == t.value {
    result += "\\n"
  } else {
    result += t.value
  }

  if nil != t.tokens {
    for _, it := range t.tokens {
      result += "\n" + it.toTreeString(spaces+step, step)
    }
    result += "\n"
  }

  if t.optional {
    result += "]"
  }
  if t.repeat {
    if t.optional {
      result += "*"
    } else {
      result += "+"
    }
  }
  return result
}

/** ***************************************************************************
 * ExpRule prepare
 ******************************************************************************/

var separators string = " \t\r\n"
var alphabet string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ|"

func (r *ExpRule) addToken(t *ptoken) {
  if nil == r.tokens {
    r.tokens = make([]*ptoken, 1)
    r.tokens[0] = t
  } else {
    r.tokens = append(r.tokens, t)
  }
}

func (r *ptoken) addToken(t *ptoken) {
  if nil == r.tokens {
    r.tokens = make([]*ptoken, 1)
    r.tokens[0] = t
  } else {
    r.tokens = append(r.tokens, t)
  }
}

/**
 * Test token item
 * @param Tested node
 * @param rules []
 * @param count matched items
 */
func (t *ptoken) test(rule IRule, n *ast.Node, rules []IRule) int {
  switch t.vtype {
  case token_string:
    if nil == n.Token || t.value != n.Token.Value {
      return -1
    }
  case token_sep_optional:
    if nil == n.Token || token.Separator != n.Token.Code {
      return 0
    }
  case token_sep:
    if nil == n.Token || token.Separator != n.Token.Code {
      return -1
    }
  case token_rule:
    if !NodeTestRuleName(n, t.value) {
      // Find rule in scope
      if nil != rules && !rule.TestName(t.value) {
        for _, r := range rules {
          if r.TestName(t.value) {
            // Process roll up AST
            count := r.Test(n, rules)
            if count > 0 {
              n.Concat(count-1, r)
              return 1
            }
          }
        }
      }
      return -1
    }
  case token_anything:
    return 1
  case token_super:
    return tests(rule, t.tokens, n, rules)
  }
  return 1
}

func tests(rule IRule, tokens []*ptoken, n *ast.Node, rules []IRule) int {
  count := 0
  iter := 0
  r := n
  if nil != tokens {
    for i := 0; i < len(tokens); i++ {
      if nil == r {
        return -1
      }
      t := tokens[i]

      if !t.repeat {
        iter = 0
      }

      if t.optional && i+1 < len(tokens) {
        // look at the front
        rs := tokens[i+1].test(rule, r, rules)
        if rs > 0 {
          i++
          count += rs
          r = r.GetNext(rs - 1)
          if t.repeat { // We can repeat?
            iter++
          }
          continue
        }
      }

      // Test at current position
      rs := t.test(rule, r, rules)
      if rs == -1 {
        if t.isnot {
          rs = 1
        } else if !t.optional && iter < 1 {
          // fmt.Println("###>", t.toString(), "'", r.ToString(true), "'")
          return -1
        }
      }
      if rs > 0 {
        count += rs
        r = r.GetNext(rs - 1)
        if t.repeat { // We can repeat?
          iter++
          i--
        }
      }
    }
  }
  return count
}

func parse_addtoken(t *ptoken, exp string, pos int, container iptoken) int /* offset */ {
  offset := 0
  if pos+1 < len(exp) {
    switch rune(exp[pos+1]) {
    case '*':
      t.repeat = true
      t.optional = true
      offset = 1
    case '+':
      t.repeat = true
      offset = 1
    }
  }
  container.addToken(t)
  return offset
}

func parse_exp_const_string(exp string, begin int) string {
  esc := false
  value := ""
  for i := begin; i < len(exp); i++ {
    c := rune(exp[i])
    switch c {
    case '\'':
      if esc {
        value += "'"
        esc = false
      } else {
        return value
      }
    case '\\':
      if esc {
        value += "\\"
        esc = false
      } else {
        esc = true
      }
    default:
      if esc {
        value += "\\"
      }
      value += string(c)
    }
  }
  return value
}

func parse_exp_until(exp string, begin int, end rune, container iptoken) int /* end index */ {
  value := ""
  not := false
  var t *ptoken = nil
  i := begin
  for ; i < len(exp); i++ {
    c := rune(exp[i])

    // Add if alphabet simbol
    if -1 != strings.Index(alphabet, string(c)) {
      value += string(c)
    } else {
      if len(value) > 0 {
        i--
        t = &ptoken{vtype: token_rule, value: value, isnot: not}
        value = ""
      } else {
        switch c {
        case '\'': // parse const string
          s := parse_exp_const_string(exp, i+1)
          t = &ptoken{vtype: token_string, value: s, isnot: not}
          i += len(s) + 1
          not = false
        case '[':
          t = &ptoken{vtype: token_super, optional: true, isnot: not}
          i = parse_exp_until(exp, i+1, ']', t)
          if -1 == i {
            return -1
          }
          not = false
        case '_':
          t = &ptoken{vtype: token_sep, value: "_", isnot: not}
          not = false
        case '%': // Anything term or noterm
          t = &ptoken{vtype: token_anything, value: "%", isnot: not}
          not = false
        case '+': // Doommy ...
          not = false
        case '!':
          not = true
        case end: // Is end
          if len(value) > 0 {
            container.addToken(&ptoken{vtype: token_rule, value: value, isnot: not})
            not = false
          }
          return i
        default:
          if -1 != strings.Index(separators, string(c)) {
            // Skip separators
            for i++; i < len(exp); i++ {
              if -1 == strings.Index(separators, string(exp[i])) {
                i--
                break
              }
            }
            t = &ptoken{vtype: token_sep_optional, value: " ", isnot: not}
            not = false
          } else {
            // Invalid parse
            fmt.Println("--------------------", string(c))
            return -1
          }
        }
      }

      // Add token to list
      if nil != t {
        i += parse_addtoken(t, exp, i, container)
      }
      t = nil
    }
  }
  if len(value) > 0 {
    container.addToken(&ptoken{vtype: token_rule, value: value, isnot: not})
  }
  return i
}

func (r *ExpRule) prepare() []*ptoken {
  if -1 == parse_exp_until(r.Expression, 0, rune(0), r) {
    // TODO error...
    return nil
  }
  fmt.Println("prepared", r.ToString(true), ";")
  return r.tokens
}
