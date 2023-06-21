package main

import (
	"fmt"
	"log"

	"github.com/alecthomas/repr"
	"github.com/chzyer/readline"
	"github.com/fatih/color"

	"github.com/aziis98/ergolas"
)

// ctx is the main repl evaluation context. This is a global as this is just a small experimental repl and this way I don't need to pass the context thorough every function call.
var ctx = ergolas.NewRootContext()

func init() {
	log.SetFlags(log.Lshortfile | log.Lmsgprefix)
}

func main() {
	color.Set(color.Italic)
	fmt.Println()
	fmt.Println("My Lang REPL")
	fmt.Println("Type 'exit <number>' or press Ctrl+C to quit.")
	fmt.Println()
	color.Unset()

	rl, err := readline.New(color.YellowString("> "))
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		processInput(line)
	}
}

func processInput(input string) {
	tokens, err := ergolas.Tokenize(input)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	node, err := ergolas.ParseExpressions(tokens)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	color.Set(color.FgBlue)
	fmt.Println("---< AST >---")
	ergolas.PrintAST(node)
	color.Unset()

	color.Set(color.FgWhite)
	fmt.Println("---< Output >---")
	result, err := ergolas.EvaluateWith(node, ctx)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	color.Unset()

	// color.Set(color.FgHiYellow)
	// fmt.Println("---< New Context >---")
	// fmt.Println(repr.String(ctx, repr.Indent("  ")))
	// color.Unset()

	color.Set(color.FgGreen)
	fmt.Println("---< Result >---")
	fmt.Println(repr.String(result, repr.Indent("  ")))
	color.Unset()

}
