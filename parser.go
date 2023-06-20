package ergolas

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

var (
	ProgramNode          NodeType = "Program"
	ExpressionsNode      NodeType = "Expressions"
	FunctionCallNode     NodeType = "FunctionCall"
	BinaryExpressionNode NodeType = "Binary"
	UnaryExpressionNode  NodeType = "Unary"
	QuotedExpressionNode NodeType = "Quoted"
	PropertyAccessNode   NodeType = "PropertyAccess"
	ParenthesisNode      NodeType = "Parenthesis"
	IdentifierNode       NodeType = "Identifier"
	BlockNode            NodeType = "Block"
	IntegerNode          NodeType = "Integer"
	FloatNode            NodeType = "Float"
	StringNode           NodeType = "String"
	OperatorNode         NodeType = "Operator"
)

type parser struct {
	tokens []Token
	cursor int

	debugStackLevel int
}

func (p *parser) log(msg string, delta int) {
	if Debug {
		log.Printf(`%s%s`, strings.Repeat("  ", p.debugStackLevel), msg)
		p.debugStackLevel += delta
	}
}

func (p *parser) done() bool {
	return p.cursor >= len(p.tokens)
}

func (p *parser) peek() Token {
	return p.tokens[p.cursor]
}

func (p *parser) advance() Token {
	p.cursor++
	return p.tokens[p.cursor-1]
}

func (p *parser) expectValue(value string) error {
	if p.done() {
		return fmt.Errorf(`expected "%s" but got eof`, value)
	}
	if p.peek().Value != value {
		return fmt.Errorf(`expected "%s" but got "%s"`, value, p.peek().Value)
	}
	p.advance()
	return nil
}

func (p *parser) expectType(typ TokenType) (Token, error) {
	if p.done() {
		return Token{}, fmt.Errorf(`expected %v but got eof`, typ)
	}
	if p.peek().Type != typ {
		return Token{}, fmt.Errorf(`expected %v but got %v`, typ, p.peek().Type)
	}
	return p.advance(), nil
}

func (p *parser) advanceLines() {
	for !p.done() && p.peek().Type == NewlineToken {
		p.advance()
	}
}

// parse has grammar
//
//	<Program> ::= <Statements>
func (p *parser) parse() (Node, error) {
	p.log(`enter parse()`, +1)
	defer p.log(`exit parse()`, -1)

	statements, err := p.parseStatements()
	if err != nil {
		return nil, err
	}

	return listNode{ProgramNode, statements}, nil
}

// parseExpressions has grammar
//
//	<Expressions> ::= <Statements>
func (p *parser) parseExpressions() (Node, error) {
	p.log(`enter parse()`, +1)
	defer p.log(`exit parse()`, -1)

	statements, err := p.parseStatements()
	if err != nil {
		return nil, err
	}

	return listNode{ExpressionsNode, statements}, nil
}

// parseStatements has grammar
//
//	<Statements> ::= ( <Expression> ";"? )*
func (p *parser) parseStatements() ([]Node, error) {
	statements := []Node{}

	p.advanceLines()

	for !p.done() && p.peek().Value != "}" {
		stmt, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		statements = append(statements, stmt)

		if !p.done() && p.peek().Value == ";" {
			p.advance()
		}

		p.advanceLines()
	}

	return statements, nil
}

// parseExpression has grammar
//
//	<Expression> ::= <RightBinaryExpression>
func (p *parser) parseExpression() (Node, error) {
	p.log(`enter parseExpression()`, +1)
	defer p.log(`exit parseExpression()`, -1)

	return p.parseRightBinaryExpression()
}

// parseRightBinaryExpression has grammar
//
//	<RightBinaryExpression> ::= <IntermediateExpression> ( ROperator <RightBinaryExpression>)?
func (p *parser) parseRightBinaryExpression() (Node, error) {
	p.log(`enter parseRightBinaryExpression()`, +1)
	defer p.log(`exit parseRightBinaryExpression()`, -1)

	lhs, err := p.parseIntermediate()
	if err != nil {
		return nil, err
	}

	if !p.done() && p.peek().Type == ROperatorToken {
		t := p.advance()
		rhs, err := p.parseRightBinaryExpression()
		if err != nil {
			return nil, err
		}

		return listNode{BinaryExpressionNode, []Node{lhs, leafNode{OperatorNode, t.Value}, rhs}}, nil
	}

	return lhs, nil
}

var functionCallTerminators = map[string]struct{}{
	";": {}, ")": {}, "]": {}, "}": {}, ":=": {}, "::": {}, "<-": {}, "->": {}, "|>": {},
}

func isFunctionCallTerminator(t Token) bool {
	_, ok := functionCallTerminators[t.Value]
	return t.Type == NewlineToken || ok
}

// parseIntermediate has grammar
//
//	<IntermediateExpression> ::= <PropertyOrValue> ( LOperator <LeftBinaryExpression> )?
//	                           | <PropertyOrValue> ( <PropertyOrValue> LOperator <LeftBinaryExpression> )+
func (p *parser) parseIntermediate() (Node, error) {
	p.log(`enter parseIntermediate()`, +1)
	defer p.log(`exit parseIntermediate()`, -1)

	node, err := p.parsePropertyOrValue()
	if err != nil {
		return nil, err
	}

	if !p.done() && p.peek().Type == LOperatorToken {
		return p.parseLeftBinaryExpression(node)
	}

	nodes := []Node{node}

	for !p.done() && !isFunctionCallTerminator(p.peek()) {
		argBase, err := p.parsePropertyOrValue()
		if err != nil {
			return nil, err
		}
		arg, err := p.parseLeftBinaryExpression(argBase)
		if err != nil {
			return nil, err
		}

		if !p.done() && p.peek().Value == "," {
			p.advance()
		}

		nodes = append(nodes, arg)
	}

	if len(nodes) > 1 {
		return listNode{FunctionCallNode, nodes}, nil
	}

	return node, nil
}

// parseLeftBinaryExpression has grammar
//
//	<LeftBinaryExpression> ::= (LOperator <PropertyOrValue>)*
func (p *parser) parseLeftBinaryExpression(lhs Node) (Node, error) {
	p.log(`enter parseLeftBinaryExpression()`, +1)
	defer p.log(`exit parseLeftBinaryExpression()`, -1)

	for !p.done() && p.peek().Type == LOperatorToken {
		t := p.advance()
		rhs, err := p.parsePropertyOrValue()
		if err != nil {
			return nil, err
		}

		lhs = listNode{BinaryExpressionNode, []Node{lhs, leafNode{OperatorNode, t.Value}, rhs}}
	}

	return lhs, nil
}

// parsePropertyOrValue has grammar
//
//	<PropertyOrValue> ::= <Value> ( "." <Identifier> )*
func (p *parser) parsePropertyOrValue() (Node, error) {
	p.log(`enter parsePropertyOrValue()`, +1)
	defer p.log(`exit parsePropertyOrValue()`, -1)

	node, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	for !p.done() && p.peek().Value == "." {
		p.expectValue(".")
		t, err := p.expectType(IdentifierToken)
		if err != nil {
			return nil, err
		}

		node = &listNode{PropertyAccessNode, []Node{
			node,
			leafNode{IdentifierNode, t.Value},
		}}
	}

	return node, nil
}

// parseValue has grammar
//
//	<Value> ::= <ParensExpression>
//	          | <BlockExpression>
//	          | <Identifier>
//	          | <Integer>
//	          | <Float>
//	          | <String>
//	          | <QuotedExpression>
func (p *parser) parseValue() (Node, error) {
	p.log(`enter parseValue()`, +1)
	defer p.log(`exit parseValue()`, -1)

	if n, err := p.parseParens(); err == nil {
		return n, nil
	}
	if n, err := p.parseBlock(); err == nil {
		return n, nil
	}
	if n, err := p.parseIdentifier(); err == nil {
		return n, nil
	}
	if n, err := p.parseInteger(); err == nil {
		return n, nil
	}
	if n, err := p.parseFloat(); err == nil {
		return n, nil
	}
	if n, err := p.parseString(); err == nil {
		return n, nil
	}
	if n, err := p.parseQuoted(); err == nil {
		return n, nil
	}

	return nil, fmt.Errorf(`expected value but got %s`, p.peek().Type)
}

// parseParens has grammar
//
//	<ParensExpression> ::= "(" <Expression> ")"
func (p *parser) parseParens() (Node, error) {
	p.log(`enter parseParens()`, +1)
	defer p.log(`exit parseParens()`, -1)

	if err := p.expectValue(`(`); err != nil {
		return nil, err
	}

	inner, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if err := p.expectValue(`)`); err != nil {
		return nil, err
	}

	return listNode{ParenthesisNode, []Node{inner}}, nil
}

// parseBlock has grammar
//
//	<Block> ::= "{" <Statements> "}"
func (p *parser) parseBlock() (Node, error) {
	p.log(`enter parseBlock()`, +1)
	defer p.log(`exit parseBlock()`, -1)

	if err := p.expectValue(`{`); err != nil {
		return nil, err
	}

	p.advanceLines()

	statements, err := p.parseStatements()
	if err != nil {
		return nil, err
	}

	p.advanceLines()

	if err := p.expectValue(`}`); err != nil {
		return nil, err
	}

	return listNode{BlockNode, statements}, nil
}

// parseQuoted has grammar
//
//	<QuotedExpression> ::= ":" <PropertyOrValue>
func (p *parser) parseQuoted() (Node, error) {
	p.log(`enter parseQuoted()`, +1)
	defer p.log(`exit parseQuoted()`, -1)

	if _, err := p.expectType(QuoteToken); err != nil {
		return nil, err
	}

	inner, err := p.parsePropertyOrValue()
	if err != nil {
		return nil, err
	}

	return listNode{QuotedExpressionNode, []Node{inner}}, nil
}

// parseInteger has grammar
//
//	<Integer> ::= Integer
func (p *parser) parseInteger() (Node, error) {
	p.log(`enter parseInteger()`, +1)
	defer p.log(`exit parseInteger()`, -1)

	t, err := p.expectType(IntegerToken)
	if err != nil {
		return nil, err
	}

	value, err := strconv.ParseInt(t.Value, 10, 64)
	if err != nil {
		return nil, err
	}

	return leafNode{IntegerNode, value}, nil
}

// parseIdentifier has grammar
//
//	<Identifier> ::= Identifier
func (p *parser) parseIdentifier() (Node, error) {
	p.log(`enter parseIdentifier()`, +1)
	defer p.log(`exit parseIdentifier()`, -1)

	t, err := p.expectType(IdentifierToken)
	if err != nil {
		return nil, err
	}

	return leafNode{IdentifierNode, t.Value}, nil
}

// parseFloat has grammar
//
//	<Float> ::= Float
func (p *parser) parseFloat() (Node, error) {
	p.log(`enter parseFloat()`, +1)
	defer p.log(`exit parseFloat()`, -1)

	t, err := p.expectType(FloatToken)
	if err != nil {
		return nil, err
	}

	value, err := strconv.ParseFloat(t.Value, 64)
	if err != nil {
		return nil, err
	}

	return leafNode{FloatNode, value}, nil
}

// parseString has grammar
//
//	<String> ::= String
func (p *parser) parseString() (Node, error) {
	p.log(`enter parseString()`, +1)
	defer p.log(`exit parseString()`, -1)

	t, err := p.expectType(StringToken)
	if err != nil {
		return nil, err
	}

	value := t.Value[1 : len(t.Value)-1] // TODO: fix escaped characters
	return leafNode{StringNode, value}, nil
}
