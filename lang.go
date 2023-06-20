package ergolas

import (
	"fmt"
	"strings"
)

var Debug = false

type NodeType string

var ErrorNodeType NodeType = "ErrorNode"

type NodeMetadata map[string]any

func (nm NodeMetadata) String() string {
	sb := &strings.Builder{}
	i := len(nm)
	for k, v := range nm {
		fmt.Fprintf(sb, `%s: "%v"`, k, v)
		i--
		if i > 0 {
			fmt.Fprint(sb, ", ")
		}
	}
	return sb.String()
}

type Node interface {
	Type() NodeType
	Children() []Node
	Metadata() NodeMetadata
}

type listNode struct {
	typ      NodeType
	children []Node
}

func (n listNode) Type() NodeType {
	return n.typ
}

func (n listNode) Children() []Node {
	return n.children
}

func (n listNode) Metadata() NodeMetadata {
	return NodeMetadata{}
}

type leafNode struct {
	typ   NodeType
	value any
}

func (n leafNode) Type() NodeType {
	return n.typ
}

func (n leafNode) Children() []Node {
	return nil
}

func (n leafNode) Metadata() NodeMetadata {
	return NodeMetadata{"Value": n.value}
}

func PrintAST(node Node) {
	printAST(node, 0)
}

func printAST(node Node, depth int) {
	indent := strings.Repeat("  ", depth)
	fmt.Printf("%s- %s", indent, node.Type())

	meta := node.Metadata()
	if len(meta) > 0 {
		fmt.Printf(" { %s }\n", meta)
	} else {
		fmt.Printf("\n")
	}

	for _, n := range node.Children() {
		printAST(n, depth+1)
	}
}

func Parse(tokens []Token) (Node, error) {
	p := &parser{tokens: tokens, cursor: 0}
	return p.parse()
}

func ParseExpression(tokens []Token) (Node, error) {
	p := &parser{tokens: tokens, cursor: 0}
	return p.parseExpression()
}

func ParseExpressions(tokens []Token) (Node, error) {
	p := &parser{tokens: tokens, cursor: 0}
	return p.parseExpressions()
}
