Coma
====

    @project coma
    @status in development
    @version 0.0.1-pre-pre-pre-alpha

The project is a context free grammar parser.


Algorithm
=========

 1. With the lexer break the input stream into tokens.
 2. Are building parsing rules in the order of the smallest particles.
 3. Consistently apply the rules to the entire input stream forming a collection of non-terminals and non-terminals.
 4. When finished successfully with the processing of the next rules start applying all the rules from the beginning, until there will be no rules that we can do.
 5. On the way out we get a parse tree AST


Grammar constructor example
===========================

```sh
# Grammar file example for JSON file
# Multi comments

lexer: # Lexer tokinizer

  digit = "0123456789".
  name = "abcdefghijklmnopqrstuvwxyz0123456789".

parser: # Parser rules

  digit           = '-'digit ; ^[0-9]+$.
  string          = '\'' %* '\'' ; '"' %* '"'.
  key             = name ; string.
  value           = key ; 'false' ; 'true' ; array.

  object          = '{' [objectKeyvalue [',' objectKeyvalue]*] '}'.
  objectKeyvalue  = key ':' value.
  array           = '[' [value [',' value]*] ']'.
```

Project tools
=============

https://github.com/demdxx/goproj


License MIT
===========

Copyright (C) <2013> Dmitry Ponomarev <demdxx@gmail.com>


Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
