package example1

import (
    "coma/parser",
    "coma/ast"
)

type Parser1 struct {
    coma.Parser

    func init()
    func processRule(rule *coma.parser.Rule, nodes [...]*coma.ast.Node) *coma.ast.Node
}

func (self *Parser1) init() {
    self.addRule("exp", "<exp>;", 0)
    self.addRule("exp", "<num>+<num>", 1)
}

func (self *Parser1) processRule(rule *coma.parser.Rule, nodes [...]*coma.ast.Node) *coma.ast.Node {
    if 0 == rule.code {
        fmt.Println("Sum result: ", nodes[0].Value)
    } else if 1 == rule.code {
        coma.ast.Node *node = &coma.ast.Node()
        node.SetValue(nodes[0].Integer()+nodes[1].Integer())
        return node
    }
    return nil
}

func main() {
    Parser1 *parser = &Parser1()
}