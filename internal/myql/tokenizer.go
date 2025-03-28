package myql

import "github.com/bzick/tokenizer"

const (
	TOK_DOT = iota + 1
	TOK_EQ
	TOK_NEQ
	TOK_GT
	TOK_GTE
	TOK_LT
	TOK_LTE
	TOK_BR_OPEN
	TOK_BR_CLOSE
	TOK_OR
	TOK_AND
	TOK_STR
)

func newQLTokenizer() *tokenizer.Tokenizer {
	t := tokenizer.New()
	t.DefineTokens(TOK_DOT, []string{"."})
	t.DefineTokens(TOK_EQ, []string{"="})
	t.DefineTokens(TOK_NEQ, []string{"!="})
	t.DefineTokens(TOK_GT, []string{">"})
	t.DefineTokens(TOK_GTE, []string{">="})
	t.DefineTokens(TOK_LT, []string{"<"})
	t.DefineTokens(TOK_LTE, []string{"<="})
	t.DefineTokens(TOK_BR_OPEN, []string{"("})
	t.DefineTokens(TOK_BR_CLOSE, []string{")"})
	t.DefineTokens(TOK_OR, []string{"||", "or"})
	t.DefineTokens(TOK_AND, []string{"&&", "and"})
	t.DefineStringToken(TOK_STR, "\"", "\"").SetEscapeSymbol(tokenizer.BackSlash)
	t.AllowKeywordSymbols(tokenizer.Underscore, tokenizer.Numbers)
	return t
}
