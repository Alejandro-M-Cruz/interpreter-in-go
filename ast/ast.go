package ast

import (
	"bytes"
	"example.com/writing-an-interpreter/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteRune(';')
	return out.String()
}

func (ls *LetStatement) statementNode() {

}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) String() string {
	return i.Value
}

func (i *Identifier) expressionNode() {

}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteRune(';')
	return out.String()
}

func (rs *ReturnStatement) statementNode() {
}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

func (es *ExpressionStatement) statementNode() {
}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

func (i *IntegerLiteral) expressionNode() {
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) String() string {
	return "(" + pe.Operator + pe.Right.String() + ")"
}

func (pe *PrefixExpression) expressionNode() {
}

type InfixExpression struct {
	Token    token.Token // The operator token
	Operator string
	Left     Expression
	Right    Expression
}

func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteRune('(')
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteRune(')')

	return out.String()
}

func (ie *InfixExpression) expressionNode() {
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return b.Token.Literal
}

func (b *Boolean) expressionNode() {
}

type Null struct {
	Token token.Token
}

func (n *Null) TokenLiteral() string {
	return n.Token.Literal
}

func (n *Null) String() string {
	return n.Token.Literal
}

func (n *Null) expressionNode() {
}

type IfExpression struct {
	Token       token.Token // the if token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

func (sl *StringLiteral) String() string {
	return sl.Value
}

func (sl *StringLiteral) expressionNode() {
}

func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" { ")
	out.WriteString(ie.Consequence.String())
	out.WriteString(" }")

	if ie.Alternative != nil {
		out.WriteString(" else")
		out.WriteString(" { ")
		out.WriteString(ie.Alternative.String())
		out.WriteString(" }")
	}

	return out.String()
}

func (ie *IfExpression) expressionNode() {
}

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (be *BlockStatement) TokenLiteral() string {
	return be.Token.Literal
}

func (be *BlockStatement) String() string {
	var statements []string

	for _, s := range be.Statements {
		statements = append(statements, s.String())
	}

	return strings.Join(statements, " ")
}

func (be *BlockStatement) statementNode() {
}

type FunctionLiteral struct {
	Token      token.Token // the fn token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	var params []string
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral() + "(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") { ")
	out.WriteString(fl.Body.String())
	out.WriteString(" }")

	return out.String()
}

func (fl *FunctionLiteral) expressionNode() {
}

type CallExpression struct {
	Token     token.Token // the ( token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *CallExpression) String() string {
	var out bytes.Buffer
	var arguments []string
	for _, a := range ce.Arguments {
		arguments = append(arguments, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteRune('(')
	out.WriteString(strings.Join(arguments, ", "))
	out.WriteRune(')')

	return out.String()
}

func (ce *CallExpression) expressionNode() {
}

type ArrayLiteral struct {
	Token    token.Token // the [ token
	Elements []Expression
}

func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}

func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	var elements []string

	for _, e := range al.Elements {
		elements = append(elements, e.String())
	}

	out.WriteRune('[')
	out.WriteString(strings.Join(elements, ", "))
	out.WriteRune(']')

	return out.String()
}

func (al *ArrayLiteral) expressionNode() {
}

type IndexExpression struct {
	Token token.Token // the [ token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteRune('(')
	out.WriteString(ie.Left.String())
	out.WriteRune('[')
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

func (ie *IndexExpression) expressionNode() {
}

type MapLiteral struct {
	Token token.Token // the { token
	Pairs map[Expression]Expression
}

func (ml *MapLiteral) TokenLiteral() string {
	return ml.Token.Literal
}

func (ml *MapLiteral) String() string {
	var out bytes.Buffer
	var pairs []string

	for k, v := range ml.Pairs {
		pairs = append(pairs, k.String()+": "+v.String())
	}

	out.WriteRune('{')
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteRune('}')

	return out.String()
}

func (ml *MapLiteral) expressionNode() {
}
