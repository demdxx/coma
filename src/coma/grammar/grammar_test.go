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