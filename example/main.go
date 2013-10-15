package main

import (
  "fmt"
  "reflect"
)

type S struct {
	i int
}

func main() {
	// s := S{}
  fmt.Println("Main test start: ", reflect.ValueOf(S{}).Type())
}
