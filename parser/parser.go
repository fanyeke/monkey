package parser

import (
	"fmt"
	"github.com/fanyeke/monkey/ast"
	"github.com/fanyeke/monkey/lexer"
	"github.com/fanyeke/monkey/token"
)

type Parser struct {
	l *lexer.Lexer // Parser 内联了 lexer.Lexer , Lexer 持有着输入的字符串
	// curToken 和 peekToken 的性质与Lexer中的当前字符和下一个字符相同, 但是它们指向的是当前词法单元和下一个词法单元
	// 原因是有可能 curToken 没有提供足够的信息, 需要下一个词法单元 peekToken 来提供
	// 例如: 读到了 5 ,这时就需要确定是一行的末尾还是算数表达式的开头
	curToken  token.Token
	peekToken token.Token

	errors []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	// 向后读取两个词法单元, 根据这两个词法单元设置 curToken 和 peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s install", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram 普拉特解析方法入口
func (p *Parser) ParseProgram() *ast.Program {
	// 构造AST的根节点
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// 遍历输入每个词法单元, 直到文件结尾 token.EOF
	for p.curToken.Type != token.EOF {
		stmt := p.ParseStatement()
		// 如果不是 nil 就插入到根节点的 program.Statements 切片中
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// ParseStatement 词法单元类型分支, 返回 ast.Statement 节点
func (p *Parser) ParseStatement() ast.Statement {
	// 词法单元的类型如果在分支中就执行对应函数, 否则返回nil
	switch p.curToken.Type {
	case token.LET:
		// 如果是 LET 类型就执行
		return p.parseLetStatement()
	default:
		return nil
	}
}

// parseLetStatement 解析 LET 主要函数
func (p *Parser) parseLetStatement() *ast.LetStatement {
	// 把当前词法单元的 token 存入 stmt 中
	stmt := &ast.LetStatement{Token: p.curToken}

	// 根据下一个词法单元判断是不是符合 LET 语句
	// 如果下一个词法单元不是一个标识符, 说明不是想要的元素, 直接返回, 词法单元指针继续前进
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	// 值得注意的是上面已经进行了指针的移动, 那么此时的Token就已经是标识符了
	// 下一个词法单元是标识符, 那么往 Statement 节点的 Name 中当前 **标识符** 存入当前 Token 的内容
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// 如果下一个词法类型不是"=", 也就不是想要的元素, 会直接返回, 词法指针继续移动
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: 跳过对表达式的处理，直到遇见分号

	// 直到遇到分号
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// peekTokenIs 检查下一个词法单元类型是否是指定类型
// 需要频繁用到, 从因此进行了抽象
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek 如果下一个词法单元是指定类型就进行跳过
// 断言函数: 主要目的是通过检查下一个语法单元的类型, 确保语法单元的正确性
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		// 注意这里会移动词法单元指针
		p.nextToken()
		return true
	} else {
		return false
	}
}
