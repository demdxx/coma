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

package ast

import (
  "coma/token"
  "fmt"
  "strings"
)

type WrapParams struct {
  Link    interface{}
  Count   int
}

type Node struct {
  Token  *token.Token
  Link    interface{}
  Wraper  WrapParams
  Parent *Node
  Prev   *Node
  Next   *Node
  Child  *Node
}

func escapeString(str string) string {
  v := strings.Replace(str, "\n", "\\n", -1)
  v = strings.Replace(v, "\t", "\\t", -1)
  // v  = strings.Replace(v, " ", "_", -1)
  return v
}

func (n *Node) ToString(childs bool) string {
  var result string = ""
  if nil != n.Token {
    if "\n" == n.Token.Value {
      result += "\\n"
    } else {
      result = n.Token.Value
    }
  }
  if childs && nil != n.Child {
    c := n.Child
    for nil != c {
      s := c.ToString(true)
      if len(s) > 0 {
        result += s + " "
      }
      c = c.Next
    }
  }
  return result
}

func (n *Node) ToTreeString(spaces int, step int) string {
  var result string = ""
  for i := 0; i < spaces; i++ {
    result += " "
  }
  if nil != n.Token {
    if token.Separator == n.Token.Code {
      result += "` `"
    } else {
      // if nil != n.Link {
      //   result += "EXP: " + escapeString(n.Link.ToString()) + " "
      // }
      result += "VAL: " + escapeString(n.Token.Value)
    }
  } else if nil != n.Link {
    result += "EXP: >> "
  }
  // else if nil != n.Link {
  //   result += "EXP: " + escapeString(n.Link.ToString()) + " >>> " + escapeString(n.ToString(true))
  // }
  result += "\n"
  if nil != n.Child {
    c := n.Child
    for nil != c {
      result += c.ToTreeString(spaces+step, step)
      c = c.Next
    }
  }
  return result
}

func (n *Node) ToStringByPrintFunc(step int, pf func (step int, n *Node) string) string {
  result := pf(step, n)
  if "" == result {
    return "<error>"
  }

  result += "\n"
  if nil != n.Child {
    c := n.Child
    for nil != c {
      result += c.ToStringByPrintFunc(step + 1, pf)
      c = c.Next
    }
  }
  return result
}

/**
 * Get the next item. Before it can be to miss a few nodes in the tree
 * @self *Node
 * @param skip count
 * @return Next node
 */
func (n *Node) GetNext(skip int) *Node {
  it := n.Next
  for i := 0; i < skip && nil != it; i++ {
    it = it.Next
  }
  return it
}

/**
 * Add next element to node
 * @self *Node
 * @param it *Node
 * @return Returns the node to which the element was added
 */
func (n *Node) AddNext(it *Node) *Node {
  r := n
  for nil != r.Next {
    r = r.Next
  }
  r.Next = it
  it.Prev = r
  return r
}

/**
 * Add child element to node
 * @self *Node
 * @param it *Node
 * @return Returns the node to which the element was added
 */
func (n *Node) AddChild(it *Node) *Node {
  if nil != n.Child {
    return n.Child.AddNext(it)
  } else {
    n.Child = it
    it.Parent = n
    it.Prev = nil
    it.Next = nil
    it.Child = nil
  }
  return n.Child
}

func (n *Node) Clean() *Node {
  // Clean from nodes without link
  var r *Node
  n2 := n
  n3 := n.Next
  for nil != n2 {
    if nil == n2.Link {
      n2.Delete()
    } else if nil == r {
      r = n2
    }
    n2 = n3
    if nil != n3 {
      n3 = n3.Next
    }
  }
  return r
}

func (n *Node) Replace(nw *Node, linkChilds bool) *Node {
  if nil != nw {
    nw.Link = n.Link
    nw.Prev = n.Prev
    nw.Next = n.Next
    nw.Parent = n.Parent
    if linkChilds {
      nw.Child = n.Child
    }
    if nil != nw.Parent {
      nw.Parent.Child = nw
    } else if nil != nw.Prev {
      nw.Prev.Next = nw
    }
    if nil != nw.Next {
      nw.Next.Prev = nw
    }
  }
  return nw
}

func (n *Node) Delete() *Node {
  if nil != n.Next {
    if nil != n.Prev {
      n.Prev.Next = n.Next
    } else if nil != n.Parent {
      n.Parent.Child = n.Next
    }
    n.Next.Prev = n.Prev
    n.Next.Parent = n.Parent
    return n.Next
  }

  // Null prev or parent links
  if nil != n.Prev {
    n.Prev.Next = nil
  }
  if nil != n.Parent {
    n.Parent.Child = nil
  }
  return n.Prev
}

/**
 * Clone current node
 * @param next
 * @param child
 * @return New Node
 */
func (n* Node) Clone(next, child int) *Node {
  node := &Node{
    Token: n.Token,
    Link: n.Link,
    Parent: nil,
    Prev: nil,
    Next: nil,
    Child: nil,
  }
  if next > 0 && nil != n.Next {
    node.Next = n.Next.Clone(next-1, child)
    node.Next.Prev = node
  }
  if child > 0 && nil != n.Child {
    node.Child = n.Child.Clone(next , child-1)
    node.Child.Parent = node
  }
  return node
}

/**
 * Concat current element which not changed
 * @self *Node
 * @param count nodes with current
 * @param link interface{}
 * @return new node
 */
func (n *Node) Concat(count int, link interface{}) *Node {
  // Get last element
  last := n
  if count > 0 {
    last = n.GetNext(count - 1)
    if nil == last {
      return nil
    }
  }

  // Make new node
  node := &Node{Link: link, Parent: n.Parent, Child: n, Prev: n.Prev, Next: last.Next}

  // Set parent node child item
  if nil != n.Parent {
    n.Parent.Child = node
  } else if nil != n.Prev {
    n.Prev.Next = node
    n.Prev = nil
  }
  if nil != node.Next {
    node.Next.Prev = node
    last.Next = nil
  }
  n.Parent = node

  // Post process node and update links if necessary
  // nw := rule.PostProcess(node)
  // if nil != nw && nw != node {
  //   node.Replace(nw, false)
  //   if nil != nw.Child {
  //     nw.Child.Clean()
  //   }
  //   return nw
  // }
  // node.Child.Clean()
  return node
}

func (n *Node) RollUp(right bool, roll func(*Node) (int /* count */, interface{} /* link */)) int {
  count := 0
  it := n
  if right {
    for it.Next != nil {
      it = it.Next
    }
    for nil != it {
      cn, r := roll(it)
      fmt.Println("RollUp: ", cn, r)
      if nil != r {
        if nil != it.Concat(cn-1, r) {
          count++
        }
      }
      it = it.Prev
    }
  } else {
    // fmt.Println("========================================================")
    for nil != it {
      cn, r := roll(it)
      if nil != r {
        nn := it.Concat(cn-1, r)
        if nil != nn {
          count++
          it = nn
        }
      }
      it = it.Next
    }
  }
  return count
}

func (n *Node) Each(f func(*Node) bool) bool {
  it := n
  for nil != it {
    if !f(it) {
      return false
    }
    it = it.Next
  }
  return true
}
