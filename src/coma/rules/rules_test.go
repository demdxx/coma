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
  // "fmt"
  "testing"
)

func TestRules(t *testing.T) {
  ruleN := MakeRule(-1, "number", "^[0-9]+$", nil)
  ruleA := MakeAlias(0, "exp", "number", nil)
  rule1 := MakeRule(0, "exp", "exp '+' exp", nil)
  rule2 := MakeRule(0, "exp", "[exp '!']* '-' [number ['+' number]*]", nil)

  t.Logf("Rule N: %s", ruleN.ToString(true))
  t.Logf("Rule A: %s", ruleA.ToString(true))
  t.Logf("Rule 1: %s", rule1.ToString(true))
  t.Logf("Rule 2: %s", rule2.ToString(true))
}
