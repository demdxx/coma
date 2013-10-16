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
  "regexp"
)

/** ***************************************************************************
 * Rule regular Expression
 ******************************************************************************/

type RegexpRule struct {
  Rule
  Expression string // regexp
}

func MakeRegexpRule(numericCode int, name string, expression string, right bool, fn postProcessFunc) *RegexpRule {
  rule := &RegexpRule{Rule: Rule{NumericCode: numericCode, Name: name, right: right, postPrecess: fn}, Expression: expression}
  return rule
}

func (rule *RegexpRule) Test(node *ast.Node, rules []IRule) int {
  // Is regular expression
  if nil != node.Token {
    matched, _ := regexp.MatchString(rule.Expression, node.Token.Value)
    if matched {
      return 1
    }
  }
  return -1
}

func (rule *RegexpRule) ToString(detail bool) string {
  result := rule.Name
  if detail {
    result += " = " + rule.Expression
  }
  return result
}
