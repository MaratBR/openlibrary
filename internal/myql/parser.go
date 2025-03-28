package myql

import (
	"fmt"

	"github.com/bzick/tokenizer"
)

type Expr interface {
}

type OPType int

const (
	BOP_UNKNOWN OPType = iota
	BOP_OR
	BOP_AND
	BOP_EQ
	BOP_NEQ
	BOP_GT
	BOP_GTE
	BOP_LT
	BOP_LTE
)

type BinaryOP struct {
	Type  OPType
	Left  Expr
	Right Expr
}

func createBOP(typ tokenizer.TokenKey, left, right Expr) *BinaryOP {
	var bopType OPType = BOP_UNKNOWN

	switch typ {
	case TOK_AND:
		bopType = BOP_AND
	case TOK_OR:
		bopType = BOP_OR
	case TOK_EQ:
		bopType = BOP_EQ
	case TOK_NEQ:
		bopType = BOP_NEQ
	case TOK_GT:
		bopType = BOP_GT
	case TOK_GTE:
		bopType = BOP_GTE
	case TOK_LT:
		bopType = BOP_LT
	case TOK_LTE:
		bopType = BOP_LTE
	}

	return &BinaryOP{
		Type:  bopType,
		Left:  left,
		Right: right,
	}
}

type Identifier struct {
	Path []string
	Prop string
}

type Parser struct {
	t    *tokenizer.Tokenizer
	s    *tokenizer.Stream
	expr Expr
	err  error
}

func NewParser() *Parser {
	t := newQLTokenizer()

	return &Parser{
		t: t,
	}
}

func (p *Parser) Parse(s string) (Expr, error) {
	p.s = p.t.ParseString(s)
	defer func() {
		p.s.Close()
		p.s = nil
	}()

	for p.s.IsValid() {
		if p.s.CurrentToken().Is(TOK_AND, TOK_OR, TOK_GT, TOK_GTE, TOK_LT, TOK_LTE) {
			opType := p.s.CurrentToken().Key()
			p.s.GoNext()
			left := p.expr

			p.parseExpr()
			if p.err != nil {
				return nil, p.err
			}

			p.expr = createBOP(opType, left, p.expr)
		} else {
			p.parseExpr()
			if p.err != nil {
				return nil, p.err
			}
		}

	}
}

func (p *Parser) parseExpr() {
	tok := p.s.CurrentToken()

	switch tok.Key() {
	case tokenizer.TokenFloat:
		p.parseFloat()
	case tokenizer.TokenInteger:
		p.parseInteger()
	case tokenizer.TokenKeyword:
		p.parseKeyword()
	case TOK_BR_OPEN:
		p.parseBrExpr()
	default:
		p.err = fmt.Errorf("unexpected token %s", tok.ValueString())
	}
}

func (p *Parser) parseKeyword() {
	tok := p.s.CurrentToken()

	if tok.Key() != tokenizer.TokenKeyword {
		panic("expected keyword")
	}

	iden := Identifier{
		Prop: tok.ValueString(),
	}

	p.s.GoNext()

	for p.s.CurrentToken().Is(TOK_DOT) {
		p.s.GoNext()

		if p.s.CurrentToken().Is(tokenizer.TokenKeyword) {
			idenStr := p.s.CurrentToken().ValueString()
			iden.Path = append(iden.Path, iden.Prop)
			iden.Prop = idenStr
		} else {
			p.err = fmt.Errorf("expected identifier, got %s", p.s.CurrentToken().ValueString())
		}
	}
}

func (p *Parser) parseInteger() {
	if p.s.CurrentToken().Key() != tokenizer.TokenInteger {
		panic("expected integer")
	}

	p.expr = p.s.CurrentToken().ValueInt64()
}

func (p *Parser) parseFloat() {
	if p.s.CurrentToken().Key() != tokenizer.TokenFloat {
		panic("expected integer")
	}

	p.expr = p.s.CurrentToken().ValueFloat64()
}

func (p *Parser) parseBrExpr() {
	if p.s.CurrentToken().Key() != TOK_BR_OPEN {
		panic("expected '('")
	}

	p.s.GoNext()

	p.parseExpr(p.s)

	if p.err != nil {
		return
	}

	if p.s.CurrentToken().Key() != TOK_BR_CLOSE {
		p.err = fmt.Errorf("expected ')', got %s", string(p.s.CurrentToken().Value()))
	}

	p.s.GoNext()

	return
}
