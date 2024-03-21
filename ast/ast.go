package ast

import (
	"bytes"
	"github.com/fanyeke/monkey/token"
)

// Node AST每个几点都需要实现Node接口
type Node interface {
	// TokenLiteral 方法返回与节点相关联词法单元的字面量,仅用于调试和测试
	// 真实的AST由相互连接的节点组成, 有的实现 Statement 接口,有的实现 Expression 接口
	TokenLiteral() string
	// String 既可以在调试时打印AST节点，也可以用来比较AST节点
	String() string
}

// Statement 实现了Node的结构体
type Statement interface {
	Node
	statementNode()
}

// Expression 实现了Node的结构体
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

// TokenLiteral 返回 Token 的内容, 其实也就是let
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

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {

}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

// String 输出 return 语句的信息
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
