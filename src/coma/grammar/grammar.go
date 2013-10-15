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
  rules.MakeRule(0, "section", "word ':' comment*", nil),
  rules.MakeRule(0, "section", "word ':' '\n'", nil),
  // rules.MakeRule(0, "string", "'\\''!\n*'\\''", nil),
  // rules.MakeRule(0, "string", "'\"'!\n*'\"'", nil),
  // rules.MakeRule(0, "rule", "word '=' ruleItemBeg* ruleItemEnd", nil),
  // rules.MakeRule(0, "ruleBeg", "word '='", nil),
  // rules.MakeRule(0, "ruleItemBeg", "!'.'*';'", nil),
  // rules.MakeRule(0, "ruleItemEnd", "!';'*'.'", nil),
  rules.MakeRule(0, "rule", "word '='%*'.'", nil),
}

func MakeGrammar() *parser.Parser {
  return parser.MakeParserByOptions(options)
}
