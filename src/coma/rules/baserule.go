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
  "regexp"
  "strings"
  // "fmt"
  // "coma/token"
  "coma/ast"
)

/** ***************************************************************************
 * Rule
 ******************************************************************************/

type IRule interface {
  Test(node *ast.Node, rules []IRule) int   // Return count matched nodes
  TestName(name string) bool // name in name|name2 => true OR name|it in name|name2 => true :)
  PostProcess(node *ast.Node) *ast.Node
  ToString(detail bool) string
  IsRight() bool
}

type postProcessFunc func(rule IRule, node *ast.Node) *ast.Node

type Rule struct {
  IRule
  NumericCode int
  Name        string
  right       bool
  rtype       int
  postPrecess postProcessFunc
}

func NodeTestRuleName(node *ast.Node, name string) bool {
  if nil != node {
    if nil != node.Link && node.Link.(IRule).TestName(name) {
      return true
    }
    // check subrules, if only one child
    // if nil != node.Child && nil == node.Child.Next && nil != node.Child.Link {
    //   return NodeTestRuleName(node.Child, name)
    // }
  }
  return false
}

func (rule *Rule) TestName(name string) bool {
  if name == rule.Name {
    return true
  }
  names := strings.Split(rule.Name, "|")
  for _, nm := range names {
    matched, _ := regexp.MatchString("^("+name+")$", nm)
    return matched
  }
  return false
}

func (rule *Rule) PostProcess(node *ast.Node) *ast.Node {
  if nil != rule.postPrecess {
    n := rule.postPrecess(rule, node)
    if nil != n {
      return n
    }
  }
  return node
}

func (rule *Rule) ToString(detail bool) string {
  return rule.Name
}

func (rule *Rule) Test(node *ast.Node, rules []IRule) int {
  return -1
}

func (rule *Rule) IsRight() bool {
  return rule.right
}
