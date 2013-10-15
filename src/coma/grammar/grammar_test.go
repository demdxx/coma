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

package grammar

import (
  "coma/ast"
  "io"
  "testing"
  "os"
  "fmt"
)

func TestParser(t *testing.T) {
	file, err := os.Open("example.grm")
	if err != nil {
	  return
	}
	defer file.Close()

	var reader io.ReadSeeker = io.ReadSeeker(file)

  grammar := MakeGrammar()
  grammar.SetStream(reader)

  // Run parser
  grammar.Parse(func(n *ast.Node) bool {
    fmt.Println(n.ToString(true))
    return true
  })
}