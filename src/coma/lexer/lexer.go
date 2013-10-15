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

package lexer

import (
  "coma/token"
  // "fmt"
  "io"
  "bytes"
  "strings"
)

const (
  Option_SeparateAlpha = iota
  Option_AbsSeparateAlpha
)

const (
  state_Error = iota
  state_Ok
)

type state int

/** ***************************************************************************
 * Structures
 ******************************************************************************/

type Option struct {
  OptType int
  Value   string
}

//
// Lexer
//

type ILexer interface {
  SetStream(r io.ReadSeeker)
  AddLex(lex *LexBase) *Lexer
  AddStdLex(lex int, lit string, prepare int) *Lexer
  AddCustomLex(lex int, test func(*Lexer, string) bool) *Lexer
  NextToken() (int, string)
}

type Lexer struct {
  ILexer
  input io.ReadSeeker

  // decode info
  lexemes []LexBase
  sep     string // Separators
  absSep  string // Absolute separators

  // info lines
  ch         rune // current character
  rdOffset   int  // reading offset (position after current character)
  lineOffset int  // current line offset

  // public state - ok to modify
  ErrorCount int // number of errors encountered
}

/** ***************************************************************************
 * Lexer declaration
 ******************************************************************************/

func MakeLexer(r io.ReadSeeker, l []LexBase) *Lexer {
  lex := new(Lexer)
  lex.init(r, l)
  return lex
}

func MakeLexerByOptions(opt []interface{}) *Lexer {
  lex := new(Lexer)
  lex.initByOptions(opt)
  return lex
}

func (l *Lexer) init(r io.ReadSeeker, lexs []LexBase) {
  if nil != r {
    l.input = r
  }
  if nil != lexs {
    l.lexemes = make([]LexBase, len(lexs))
    copy(l.lexemes, lexs)
  }
}

func (l *Lexer) initByOptions(opt []interface{}) {
  if nil != opt {
    for _, it := range opt {
      switch it.(type) {
      case *Lex:
        l.AddLex(it.(LexBase))
      case *LexCustom:
        l.AddLex(it.(LexBase))
      case *Option:
        switch it.(*Option).OptType {
        case Option_SeparateAlpha:
          l.SetSeparateAlpha(it.(*Option).Value)
        case Option_AbsSeparateAlpha:
          l.SetAbsSeparateAlpha(it.(*Option).Value)
        default:
          // TODO: Error...
        }
      case Option:
        switch it.(Option).OptType {
        case Option_SeparateAlpha:
          l.SetSeparateAlpha(it.(Option).Value)
        case Option_AbsSeparateAlpha:
          l.SetAbsSeparateAlpha(it.(Option).Value)
        default:
          // TODO: Error...
        }
      default:
        // TODO: Error...
      }
    }
  }
}

func (l *Lexer) SetSeparateAlpha(separators string) {
  l.sep = separators
}

func (l *Lexer) SetAbsSeparateAlpha(separators string) {
  l.absSep = separators
}

/** ***************************************************************************
 * Lexer actions
 ******************************************************************************/

func (l *Lexer) SetStream(r io.ReadSeeker) {
  l.input = r
  l.rdOffset = 0
  l.lineOffset = 0
  l.ErrorCount = 0
}

func (l *Lexer) AddLex(lex LexBase) *Lexer {
  if nil == l.lexemes {
    l.lexemes = make([]LexBase, 1)
    l.lexemes[0] = lex
  } else {
    l.lexemes = append(l.lexemes, lex)
  }
  return l
}

func (l *Lexer) AddStdLex(code int, lit string, prepare int) *Lexer {
  t := &Lex{code: code, lit: lit, prepare: prepare}
  return l.AddLex(t.LexBase)
}

func (l *Lexer) AddCustomLex(code int, test func(*Lexer, string) bool) *Lexer {
  t := &LexCustom{code: code, testFunc: test}
  return l.AddLex(t.LexBase)
}

/** ***************************************************************************
 * Next lex process helpers
 ******************************************************************************/

func (l Lexer) IsSeparator(b byte) bool {
  return -1 != strings.Index(l.sep, string(b))
}

func (l Lexer) IsAbsoluteSeparator(b byte) bool {
  return -1 != strings.Index(l.absSep, string(b))
}

/** ***************************************************************************
 * Next token process
 ******************************************************************************/

const (
  word_word = iota
  word_sep
  word_absep
  word_eof
  word_error
)

func (l *Lexer) nextByte() (byte, error) {
  buff := make([]byte, 1)
  _, error := l.input.Read(buff)
  if nil != error {
    return 0x00, error
  }
  l.rdOffset++
  l.ch = rune(buff[0])
  if '\n' == l.ch {
    l.lineOffset++
  }
  return byte(l.ch), error
}

func (l *Lexer) moveBack(offset int64) {
  if '\n' == l.ch {
    l.lineOffset--
    l.ch = 0
  }
  l.rdOffset -= int(offset)
  l.input.Seek(-offset, 1)
}

func (l *Lexer) nextWord() (int, string) {
  // result := ""
  var result bytes.Buffer
  for {
    b, error := l.nextByte()
    if nil != error {
      if io.EOF == error {
        if result.Len() > 0 {
          return word_word, result.String()
        }
        break
      }
      return word_error, "io Error"
    }
    if l.IsAbsoluteSeparator(b) {
      if result.Len() > 0 {
        l.moveBack(1) // Move back 1 byte
        return word_word, result.String()
      }

      // All separators
      result.WriteByte(b)
      for {
        c, err := l.nextByte()
        if nil != err || !l.IsAbsoluteSeparator(c) {
          if nil == err {
            l.moveBack(1)
          }
          break
        }
        result.WriteByte(c)
      }
      return word_absep, result.String()
    } else if l.IsSeparator(b) {
      if result.Len() < 1 {
        result.WriteByte(b)
      } else {
        // fmt.Println("...", result)
        l.moveBack(1) // Move back 1 byte
      }
      return word_sep, result.String()
    } else {
      result.WriteByte(b)
    }
  }
  return word_eof, "EOF"
}

/**
 * Get next token
 * @owner Lexer
 * @return tokenType, *token, tokenWord
 */
func (l *Lexer) NextToken() (int, *LexBase, string) {
  prevword := ""
  for {
    state, word := l.nextWord()
    switch state {
    case word_eof:
      if len(prevword) > 0 {
        if 1 == len(prevword) {
          return token.Simbol, nil, prevword
        } else {
          l.moveBack(int64(len(prevword)) - 1) // Move back
          return token.Simbol, nil, string(prevword[0])
        }
      }
      return token.EOF, nil, "EOF"
    case word_error:
      return token.Error, nil, "Error"
    case word_word:
      if len(prevword) > 0 {
        if 1 == len(prevword) {
          l.moveBack(int64(len(word)))
          return token.Simbol, nil, prevword
        } else {
          // fmt.Println("+++++1", token.Simbol, ":", "`"+prevword+word+"`", len(prevword) + len(word))
          l.moveBack(int64(len(prevword) + len(word) - 1)) // Move back
          return token.Simbol, nil, string(prevword[0])
        }
      }
      if lex, w := l.test(word); nil != lex {
        return (*lex).GetCode(), lex, w
      }
      break
    case word_sep:
      prevword += word
      if lex, w := l.test(prevword); nil != lex {
        return (*lex).GetCode(), lex, w
      }
    case word_absep:
      if len(prevword) > 0 {
        if 1 == len(prevword) {
          l.moveBack(1) // Move back
          return token.Simbol, nil, prevword
        } else {
          // fmt.Println("+++++2", token.Simbol, ":", "`"+prevword+"`", len(prevword))
          l.moveBack(int64(len(prevword))) // Move back
          return token.Simbol, nil, string(prevword[0])
        }
      }
      return token.Separator, nil, word
    default:
      goto error
    }
    // fmt.Println("+++", state, ":", "`"+word+"`")
  }
error:
  return token.Error, nil, ""
}

/** ***************************************************************************
 * Test lexer
 ******************************************************************************/

/**
 * Test word by tokens
 * @owner Lexer
 * @param word
 * @return bool
 */
func (l *Lexer) test(word string) (*LexBase, string) {
  // fmt.Println("test", "`"+word+"`", len(l.tokens))
  for _, lex := range l.lexemes {
    // fmt.Println("test at", lex)
    if lex.test(l, word) {
      return &lex, word
    }
  }
  return nil, ""
}
