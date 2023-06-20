package ergolas

import (
	"fmt"
	"math"
	"os"
)

type Context struct {
	Parent   *Context
	Bindings map[string]any
}

func (ctx *Context) GetKey(name string) (any, error) {
	value, ok := ctx.Bindings[name]
	if !ok {
		if ctx.Parent != nil {
			return ctx.Parent.GetKey(name)
		} else {
			return nil, fmt.Errorf(`unbound variable "%s"`, name)
		}
	}

	return value, nil
}

func NewRootContext() *Context {
	return &Context{nil, map[string]any{
		"exit": func(args ...any) (any, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf(`expected 1 argument, got %v`, len(args))
			}

			nExitCode, ok := args[0].(int)
			if !ok {
				return nil, fmt.Errorf(`expected integer but got %v`, nExitCode)
			}

			os.Exit(nExitCode)
			return nil, nil
		},
		"println": func(args ...any) (any, error) {
			for _, arg := range args {
				fmt.Print(arg)
			}
			fmt.Println()

			return nil, nil
		},
		"true":  true,
		"false": false,
	}}
}

func isTruthy(v any) bool {
	if b, ok := v.(bool); ok {
		return b
	}

	return v != nil
}

func eval(node Node, ctx *Context) (any, error) {
	switch node.Type() {
	case ProgramNode:
		for _, n := range node.Children() {
			_, err := eval(n, ctx)
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	case ExpressionsNode:
		var lastResult any

		for _, n := range node.Children() {
			var err error
			if lastResult, err = eval(n, ctx); err != nil {
				return nil, err
			}
		}

		return lastResult, nil
	case FunctionCallNode:
		calleeAst := node.Children()[0]
		argsAst := node.Children()[1:]

		vCallee, err := eval(calleeAst, ctx)
		if err != nil {
			return nil, err
		}

		vArgs := []any{}
		for _, argAst := range argsAst {
			vArg, err := eval(argAst, ctx)
			if err != nil {
				return nil, err
			}

			vArgs = append(vArgs, vArg)
		}

		fn, ok := vCallee.(func(args ...any) (any, error))
		if !ok {
			return nil, fmt.Errorf(`not a function: %v`, vCallee)
		}

		return fn(vArgs...)
	case BinaryExpressionNode:
		lhs := node.Children()[0]
		op := node.Children()[1].Metadata()["Value"].(string)
		rhs := node.Children()[2]

		if op == ":=" {
			if lhs.Type() != IdentifierNode {
				return nil, fmt.Errorf(`expected identifier on left side of assignment`)
			}

			name := lhs.Metadata()["Value"].(string)

			vRhs, err := eval(rhs, ctx)
			if err != nil {
				return nil, err
			}

			ctx.Bindings[name] = vRhs
			return nil, nil
		}
		if op == "&&" {
			vLhs, err := eval(lhs, ctx)
			if err != nil {
				return nil, err
			}

			if !isTruthy(vLhs) {
				return vLhs, nil
			}

			vRhs, err := eval(rhs, ctx)
			if err != nil {
				return nil, err
			}

			return vRhs, nil
		}
		if op == "||" {
			vLhs, err := eval(lhs, ctx)
			if err != nil {
				return nil, err
			}

			if isTruthy(vLhs) {
				return vLhs, nil
			}

			vRhs, err := eval(rhs, ctx)
			if err != nil {
				return nil, err
			}

			return vRhs, nil
		}

		vLhs, err := eval(lhs, ctx)
		if err != nil {
			return nil, err
		}
		vRhs, err := eval(rhs, ctx)
		if err != nil {
			return nil, err
		}

		switch op {
		case "+":
			if nLhs, ok := vLhs.(int64); ok {
				if nRhs, ok := vRhs.(int64); ok {
					return nLhs + nRhs, nil
				}
			}
			if nLhs, ok := vLhs.(float64); ok {
				if nRhs, ok := vRhs.(float64); ok {
					return nLhs + nRhs, nil
				}
			}
			if nLhs, ok := vLhs.(string); ok {
				if nRhs, ok := vRhs.(string); ok {
					return nLhs + nRhs, nil
				}
			}

			return nil, fmt.Errorf(`cannot apply operator "+" to types %T and %T`, vLhs, vRhs)
		case "-":
			if nLhs, ok := vLhs.(int64); ok {
				if nRhs, ok := vRhs.(int64); ok {
					return nLhs - nRhs, nil
				}
			}
			if nLhs, ok := vLhs.(float64); ok {
				if nRhs, ok := vRhs.(float64); ok {
					return nLhs - nRhs, nil
				}
			}

			return nil, fmt.Errorf(`cannot apply operator "-" to types %T and %T`, vLhs, vRhs)
		case "*":
			if nLhs, ok := vLhs.(int64); ok {
				if nRhs, ok := vRhs.(int64); ok {
					return nLhs * nRhs, nil
				}
			}
			if nLhs, ok := vLhs.(float64); ok {
				if nRhs, ok := vRhs.(float64); ok {
					return nLhs * nRhs, nil
				}
			}

			return nil, fmt.Errorf(`cannot apply operator "*" to types %T and %T`, vLhs, vRhs)
		case "/":
			if nLhs, ok := vLhs.(int64); ok {
				if nRhs, ok := vRhs.(int64); ok {
					return nLhs / nRhs, nil
				}
			}
			if nLhs, ok := vLhs.(float64); ok {
				if nRhs, ok := vRhs.(float64); ok {
					return nLhs / nRhs, nil
				}
			}

			return nil, fmt.Errorf(`cannot apply operator "/" to types %T and %T`, vLhs, vRhs)
		case "%":
			if nLhs, ok := vLhs.(int64); ok {
				if nRhs, ok := vRhs.(int64); ok {
					return nLhs % nRhs, nil
				}
			}
			if nLhs, ok := vLhs.(float64); ok {
				if nRhs, ok := vRhs.(float64); ok {
					return math.Mod(nLhs, nRhs), nil
				}
			}

			return nil, fmt.Errorf(`cannot apply operator "%%" to types %T and %T`, vLhs, vRhs)
		}

	case QuotedExpressionNode:
		return node, nil

	case PropertyAccessNode:
		return nil, fmt.Errorf(`not implemented`)

	case ParenthesisNode:
		return eval(node.Children()[0], ctx)

	case IdentifierNode:
		name := node.Metadata()["Value"].(string)
		return ctx.GetKey(name)

	case BlockNode:
		return nil, fmt.Errorf(`not implemented`)

	case IntegerNode:
		return node.Metadata()["Value"], nil

	case FloatNode:
		return node.Metadata()["Value"], nil

	case StringNode:
		return node.Metadata()["Value"], nil
	}

	return nil, fmt.Errorf(`unexpected node %T`, node)
}

func Evaluate(node Node) (any, error) {
	return eval(node, NewRootContext())
}

func EvaluateWith(node Node, ctx *Context) (any, error) {
	return eval(node, ctx)
}
