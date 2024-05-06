package ast

import (
	"bytes"
	"github.com/fanyeke/monkey/token"
	"strings"
)

/*
ast是定义一些需要挂在语法解析树上的单元
首先可以先来看两个接口 Statement 和 Expression,
他们两个就是为了规范一些节点的实现, 他们也把节点分为了两类节点:
一类是 Statement 节点, 可以理解为是语句, 我们可以看一些他的实现, 例如 LetStatement(VarStatement) 就是let(var) x = 1这一条语句的节点
另一类是 Expression 节点, 可以理解为是表达式, 不仅指算数表达式, "x"也是一个表达式(Identifier),
因此AST语法解析树的解构就明晰了: 每一次输入会产生一颗AST语法解析树, 叶子节点都是表达式, 语句是关联叶子节点的整体
*/

// Node AST每个几点都需要实现Node接口
type Node interface {
	// TokenLiteral 方法返回与节点相关联词法单元的字面量,仅用于调试和测试
	// 真实的AST由相互连接的节点组成, 有的实现 Statement 接口,有的实现 Expression 接口
	TokenLiteral() string
	// String 既可以在调试时打印AST节点，也可以用来比较AST节点
	String() string
}

// Statement 语句实现了Node的结构体
type Statement interface {
	Node
	statementNode()
}

// Expression 表达式实现了Node的结构体
type Expression interface {
	Node
	expressionNode()
}

// Program 存储的是 Statement 切片， 而它又实现了 Node 结构体，
// 因此这是一个保存信息的组
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// String 输出 Program 中每个 Node 的信息
func (p *Program) String() string {
	var out bytes.Buffer
	// 拼接信息即可
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type LetStatement struct {
	Token token.Token // token.LET 词法单元
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {

}

// TokenLiteral 返回 Token 的内容, 其实也就是 let
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// String 输入 LetStatement 语句的信息
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	// 写入 Token 的内容,即let
	out.WriteString(ls.TokenLiteral() + " ")
	// 写入"="之前的信息
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	// 写入"="后的内容
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

// Identifier 标识符
type Identifier struct {
	Token token.Token // token.IDENT 词法单元
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) String() string {
	return i.Value
}

// ReturnStatement return 表达式
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {

}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	// 写入 return
	out.WriteString(rs.TokenLiteral() + " ")
	// 如果 return 后还有词法单元, 则进行拼接
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatement 语句?
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// IntegerLiteral  整数字面值
type IntegerLiteral struct {
	Token token.Token
	//该字段将包含整数字面量在源代码中的实际值, 在构建时, 需要把字符串转换为int64类型
	Value int64
}

func (il *IntegerLiteral) expressionNode() {

}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

// PrefixExpression 前缀解析结构体
type PrefixExpression struct {
	Token token.Token // 前缀词法单元,如"!","-"
	// Operator是包含"-"或"!"的字符串；Right字段包含运算符右边的表达式。
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {

}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// InfixExpression 中缀解析结构体
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression

	// 中间代码表达式
	Quads   string
	TempVar string
}

func (ie *InfixExpression) expressionNode() {

}
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// Boolean 布尔值字面量
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

// IfExpression if语句解析结构体
type IfExpression struct {
	Token       token.Token     // if
	Condition   Expression      // 条件
	Consequence *BlockStatement // 结果
	Alternative *BlockStatement // else
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

// BlockStatement "{}"区块语句
type BlockStatement struct {
	Token      token.Token // "{"词法单元
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// FunctionLiteral 函数子面值解析
type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

// CallExpression 表达式字面值解析
type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// StringLiteral 字符串节点
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer
	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}
