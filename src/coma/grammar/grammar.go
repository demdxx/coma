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
  "coma/rules"
  "coma/parser"
  "coma/lexer"
  "coma/token"
)

/**
# TODO:
# Here we can init context
# CONTEXT_VAR = value;

lexer:

  digit = "0123456789"
  word = "abcdefghijklmnopqrstuvwxyz0123456789_";

parser:

  string = '\'' !['\n']* '\'' ; '"' !['\n']* '"'.
  comment = '#' %* '\n'.
  section = NEWLINE?'lexer:' ; NEWLINE?'parser:'.
  rule = NEWLINE? word '=' !comment+ ';'.

JSON BNF:

element
  = '[' tag-name ',' attributes ',' element-list ']'
  | '[' tag-name ',' attributes ']'
  | '[' tag-name ',' element-list ']'
  | '[' tag-name ']'
  | string
  ;
tag-name
  = string
  ;
attributes
  = '{' attribute-list '}'
  | '{' '}'
  ;
attribute-list
  = attribute ',' attribute-list
  | attribute
  ;
attribute
  = attribute-name ':' attribute-value
  ;
attribute-name
  = string
  ;
attribute-value
  = string
  | number
  | 'true'
  | 'false'
  | 'null'
  ;
element-list
  = element ',' element-list
  | element
  ;

*/

var options = []interface{}{
  // Lexer options
  lexer.Option{lexer.Option_SeparateAlpha, ".,:;`~!?@#$%%^&*-+=_/\\|(){}[]<>\"' \t\n"},
  lexer.Option{lexer.Option_AbsSeparateAlpha, " \t\r"},

  lexer.MakeLex(token.Word, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", lexer.PrepareAs_Alpha),

  // Parser options
  rules.MakeRule(-1, "comment", "'#'%*'\n'", nil),
  rules.MakeRule(-1, "comment", "comment comment+", nil),
  rules.MakeRule(0, "word", "^[a-zA-Z0-9]+$", nil),
  rules.MakeRule(0, "string", "'\\''%*'\\''", nil),
  rules.MakeRule(0, "string", "'\"'%*'\"'", nil),
  rules.MakeRule(0, "section", "word ':' comment*", nil),
  rules.MakeRule(0, "section", "word ':' '\n'", nil),
  // rules.MakeRule(0, "ruleItemBeg", "!'.'*';'", nil),
  // rules.MakeRule(0, "ruleItemEnd", "!';'*'.'", nil),
  // rules.MakeRule(0, "rule", "word '=' ruleItemBeg* ruleItemEnd", nil),
  // rules.MakeRule(0, "ruleBeg", "word '='", nil),
  rules.MakeRule(0, "rule", "word '='%*'.'[['\n'][ ]]*", nil),
  // rules.MakeRule(0, "rule", "word '=' [!';'+ ';']* '.'", nil),
  // rules.MakeRule(0, "rule", "rule'\n'+", nil),
  // rules.MakeRule(0, "rules", "rule+", nil),
  rules.MakeRule(0, "section_block", "section '\n'* rule+", nil),
}

func MakeGrammar() *parser.Parser {
  return parser.MakeParserByOptions(options)
}
