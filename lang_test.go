package ergolas_test

import (
	"fmt"
	"log"

	"github.com/aziis98/ergolas"
)

func ExampleParse_function_call_precedence() {
	tokens, err := ergolas.Tokenize(`f x + f y`)
	if err != nil {
		log.Fatal(err)
	}

	node, err := ergolas.ParseExpression(tokens)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ast:")
	ergolas.PrintAST(node)

	// Output:
	// Ast:
	// - FunctionCall
	//   - Identifier { Value: "f" }
	//   - Binary
	//     - Identifier { Value: "x" }
	//     - Operator { Value: "+" }
	//     - Identifier { Value: "f" }
	//   - Identifier { Value: "y" }

}

func ExampleParse_quasiquote() {
	tokens, err := ergolas.Tokenize(`:(1 + 2 + $(2 * 2))`)
	if err != nil {
		log.Fatal(err)
	}

	node, err := ergolas.ParseExpression(tokens)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ast:")
	ergolas.PrintAST(node)

	// Output:
	// Ast:
	// - Quoted
	//   - Parenthesis
	//     - Binary
	//       - Binary
	//         - Integer { Value: "1" }
	//         - Operator { Value: "+" }
	//         - Integer { Value: "2" }
	//       - Operator { Value: "+" }
	//       - Unquote
	//         - Parenthesis
	//           - Binary
	//             - Integer { Value: "2" }
	//             - Operator { Value: "*" }
	//             - Integer { Value: "2" }
}

func ExampleParse_complex_expression() {
	tokens, err := ergolas.Tokenize(`foo (bar 1 2 3 "hi") (baz (3 * 4.0 + (2 ^ 3)) :symbol)`)
	if err != nil {
		log.Fatal(err)
	}

	node, err := ergolas.Parse(tokens)
	if err != nil {
		log.Fatal(err)
	}

	ergolas.PrintAST(node)

	// Output:
	// - Program
	//   - FunctionCall
	//     - Identifier { Value: "foo" }
	//     - Parenthesis
	//       - FunctionCall
	//         - Identifier { Value: "bar" }
	//         - Integer { Value: "1" }
	//         - Integer { Value: "2" }
	//         - Integer { Value: "3" }
	//         - String { Value: "hi" }
	//     - Parenthesis
	//       - FunctionCall
	//         - Identifier { Value: "baz" }
	//         - Parenthesis
	//           - Binary
	//             - Binary
	//               - Integer { Value: "3" }
	//               - Operator { Value: "*" }
	//               - Float { Value: "4" }
	//             - Operator { Value: "+" }
	//             - Parenthesis
	//               - Binary
	//                 - Integer { Value: "2" }
	//                 - Operator { Value: "^" }
	//                 - Integer { Value: "3" }
	//         - Quoted
	//           - Identifier { Value: "symbol" }
}

func ExampleParse_inline_if_blocks() {
	tokens, err := ergolas.Tokenize(`if { ans == 42 } { println "Yep" } { println "Nope" }`)
	if err != nil {
		log.Fatal(err)
	}

	node, err := ergolas.Parse(tokens)
	if err != nil {
		log.Fatal(err)
	}

	ergolas.PrintAST(node)

	// Output:
	// - Program
	//   - FunctionCall
	//     - Identifier { Value: "if" }
	//     - Block
	//       - Binary
	//         - Identifier { Value: "ans" }
	//         - Operator { Value: "==" }
	//         - Integer { Value: "42" }
	//     - Block
	//       - FunctionCall
	//         - Identifier { Value: "println" }
	//         - String { Value: "Yep" }
	//     - Block
	//       - FunctionCall
	//         - Identifier { Value: "println" }
	//         - String { Value: "Nope" }
}

func ExampleParse_multiline_if_blocks() {
	tokens, err := ergolas.Tokenize(`
		if { ans == 42 } { 
			println "Yep"
		} {
			println "Nope"
		}
	`)
	if err != nil {
		log.Fatal(err)
	}

	node, err := ergolas.Parse(tokens)
	if err != nil {
		log.Fatal(err)
	}

	ergolas.PrintAST(node)

	// Output:
	// - Program
	//   - FunctionCall
	//     - Identifier { Value: "if" }
	//     - Block
	//       - Binary
	//         - Identifier { Value: "ans" }
	//         - Operator { Value: "==" }
	//         - Integer { Value: "42" }
	//     - Block
	//       - FunctionCall
	//         - Identifier { Value: "println" }
	//         - String { Value: "Yep" }
	//     - Block
	//       - FunctionCall
	//         - Identifier { Value: "println" }
	//         - String { Value: "Nope" }

}

func ExampleEvaluate_arithmetic() {
	tokens, err := ergolas.Tokenize(`1 + 2 * 3`)
	if err != nil {
		log.Fatal(err)
	}

	node, err := ergolas.ParseExpression(tokens)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ast:")
	ergolas.PrintAST(node)

	result, err := ergolas.Evaluate(node)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Value:")
	fmt.Printf("%v\n", result)

	// Output:
	// Ast:
	// - Binary
	//   - Binary
	//     - Integer { Value: "1" }
	//     - Operator { Value: "+" }
	//     - Integer { Value: "2" }
	//   - Operator { Value: "*" }
	//   - Integer { Value: "3" }
	// Value:
	// 9
}

func ExampleEvaluate_conditionals() {
	tokens, err := ergolas.Tokenize(`
		println ">>> " (false && false) " ~ false"
		println ">>> " ("hi" && false) " ~ false"
		println ">>> " (3.0 && false) " ~ false"
		println ">>> " (false && true) " ~ false"
		println ">>> " (true && :example) " ~ :example"
		println ">>> " (true && true) " ~ true"
	`)
	if err != nil {
		log.Fatal(err)
	}

	node, err := ergolas.Parse(tokens)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ast:")
	ergolas.PrintAST(node)

	fmt.Println("Value:")
	if _, err := ergolas.Evaluate(node); err != nil {
		log.Fatal(err)
	}

	// Output:
	// Ast:
	// - Program
	//   - FunctionCall
	//     - Identifier { Value: "println" }
	//     - String { Value: ">>> " }
	//     - Parenthesis
	//       - Binary
	//         - Identifier { Value: "false" }
	//         - Operator { Value: "&&" }
	//         - Identifier { Value: "false" }
	//     - String { Value: " ~ false" }
	//   - FunctionCall
	//     - Identifier { Value: "println" }
	//     - String { Value: ">>> " }
	//     - Parenthesis
	//       - Binary
	//         - String { Value: "hi" }
	//         - Operator { Value: "&&" }
	//         - Identifier { Value: "false" }
	//     - String { Value: " ~ false" }
	//   - FunctionCall
	//     - Identifier { Value: "println" }
	//     - String { Value: ">>> " }
	//     - Parenthesis
	//       - Binary
	//         - Float { Value: "3" }
	//         - Operator { Value: "&&" }
	//         - Identifier { Value: "false" }
	//     - String { Value: " ~ false" }
	//   - FunctionCall
	//     - Identifier { Value: "println" }
	//     - String { Value: ">>> " }
	//     - Parenthesis
	//       - Binary
	//         - Identifier { Value: "false" }
	//         - Operator { Value: "&&" }
	//         - Identifier { Value: "true" }
	//     - String { Value: " ~ false" }
	//   - FunctionCall
	//     - Identifier { Value: "println" }
	//     - String { Value: ">>> " }
	//     - Parenthesis
	//       - Binary
	//         - Identifier { Value: "true" }
	//         - Operator { Value: "&&" }
	//         - Quoted
	//           - Identifier { Value: "example" }
	//     - String { Value: " ~ :example" }
	//   - FunctionCall
	//     - Identifier { Value: "println" }
	//     - String { Value: ">>> " }
	//     - Parenthesis
	//       - Binary
	//         - Identifier { Value: "true" }
	//         - Operator { Value: "&&" }
	//         - Identifier { Value: "true" }
	//     - String { Value: " ~ true" }
	// Value:
	// >>> false ~ false
	// >>> false ~ false
	// >>> false ~ false
	// >>> false ~ false
	// >>> {Quoted [{Identifier example}]} ~ :example
	// >>> true ~ true

}
