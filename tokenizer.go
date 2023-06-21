package ergolas

import (
	"fmt"
	"regexp"
	"strings"
)

type TokenType string

type Token struct {
	Type     TokenType
	Value    string
	Location int
}

func computeLineColumn(source string, index int) (line, column int) {
	lines := strings.Split(source, "\n")
	totalChars := 0

	for i, line := range lines {
		lineLength := len(line) + 1
		if index < totalChars+lineLength {
			lineIndex := index - totalChars
			return i + 1, lineIndex + 1
		}
		totalChars += lineLength
	}

	panic("character index out of range")
}

type TokenizeError struct {
	Source   *string
	Location int
	Message  string
}

func (e TokenizeError) Error() string {
	line, col := computeLineColumn(*e.Source, e.Location)
	return fmt.Sprintf(`[%d:%d] %s`, line, col, e.Message)
}

type rule struct {
	Type   TokenType
	Regex  *regexp.Regexp
	Ignore bool
}

var (
	FloatToken       TokenType = "Float"
	IntegerToken     TokenType = "Integer"
	StringToken      TokenType = "String"
	QuoteToken       TokenType = "Quote"
	UnquoteToken     TokenType = "Unquote"
	LOperatorToken   TokenType = "LOperator"
	ROperatorToken   TokenType = "ROperator"
	PunctuationToken TokenType = "Punctuation"
	IdentifierToken  TokenType = "Identifier"
	CommentToken     TokenType = "Comment"
	WhitespaceToken  TokenType = "Whitespace"
	NewlineToken     TokenType = "Newline"
)

var rules = []rule{
	{Type: FloatToken,
		Regex: regexp.MustCompile(`^[0-9]+\.[0-9]+`)},
	{Type: IntegerToken,
		Regex: regexp.MustCompile(`^[0-9]+`)},
	{Type: StringToken,
		Regex: regexp.MustCompile(`^"(\\.|[^"])*"`)},
	{Type: ROperatorToken, // The operators ":=", "::", "<-", "->" and "|>" are right associative
		Regex: regexp.MustCompile(`^(\:\=|\:\:)`)},
	{Type: QuoteToken,
		Regex: regexp.MustCompile(`^:`)},
	{Type: UnquoteToken,
		Regex: regexp.MustCompile(`^\$`)},
	{Type: LOperatorToken,
		Regex: regexp.MustCompile(`^[\+\-\*\/\%\=\<\>\!\&\|\^]+`)},
	{Type: PunctuationToken,
		Regex: regexp.MustCompile(`^[\.\,\;\(\)\[\]\{\}]`)},
	{Type: IdentifierToken,
		Regex: regexp.MustCompile(`^[a-zA-Z\-\_\$][a-zA-Z0-9\-\_\$]*`)},
	{Type: NewlineToken,
		Regex: regexp.MustCompile(`^\n\s*`)},
	{Type: CommentToken, Ignore: true,
		Regex: regexp.MustCompile(`^#.*`)},
	{Type: WhitespaceToken, Ignore: true,
		Regex: regexp.MustCompile(`^[ \t]+`)},
}

func matchRules(source string) (*Token, bool) {
	for _, rule := range rules {
		match := rule.Regex.FindString(source)
		if match != "" {
			return &Token{Type: rule.Type, Value: match}, rule.Ignore
		}
	}

	return nil, true
}

func Tokenize(source string) ([]Token, error) {
	cursor := 0
	tokens := []Token{}

	for cursor < len(source) {
		remaining := source[cursor:]

		t, ignore := matchRules(remaining)
		if t == nil {
			return nil, TokenizeError{&source, cursor, "unexpected character"}
		}

		cursor += len(t.Value)
		if !ignore {
			tokens = append(tokens, *t)
		}
	}

	return tokens, nil
}
