package parser

import (
	"fmt"
	"github.com/fanyeke/monkey/ast"
	"github.com/fanyeke/monkey/lexer"
	"github.com/fanyeke/monkey/token"
	"strconv"
)

// 这些常量是用来区分运算符优先级的
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

// precedences 中缀解析需要的一些符号和优先级的映射表
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

// 定义两种类型的函数: 前缀解析函数和中缀解析函数
// 两个函数均返回 ast.Expression
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)
type Parser struct {
	l      *lexer.Lexer // Parser 内联了 lexer.Lexer , Lexer 持有着输入的字符串
	errors []string

	// curToken 和 peekToken 的性质与Lexer中的当前字符和下一个字符相同, 但是它们指向的是当前词法单元和下一个词法单元
	// 原因是有可能 curToken 没有提供足够的信息, 需要下一个词法单元 peekToken 来提供
	// 例如: 读到了 5 ,这时就需要确定是一行的末尾还是算数表达式的开头
	curToken  token.Token
	peekToken token.Token

	// 使用 prefixParseFns 和 infixParseFns 两个 map 来保存映射函数, 每一种 token.TokenType 会对应一种处理函数
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	// 初始化解析函数的映射map
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	// 注册前缀解析函数
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	// 初始化中缀解析函数的映射
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	// 注册布尔值解析函数
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	// 注册括号的解析函数
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	// 注册花括号的解析函数
	p.registerPrefix(token.IF, p.parseIfExpression)
	// 注册函数解析函数
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	// 注册中缀解析函数
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	// 注册表达式解析函数
	p.registerInfix(token.LPAREN, p.parseCallExpression)
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
		stmt := p.parseStatement()
		// 如果不是 nil 就插入到根节点的 program.Statements 切片中
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// ParseStatement 词法单元类型分支, 返回 ast.Statement 节点
func (p *Parser) parseStatement() ast.Statement {
	// 词法单元的类型如果在分支中就执行对应函数, 否则返回nil
	switch p.curToken.Type {
	case token.LET:
		// 如果是 LET 类型就执行
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:

		return p.parseExpressionStatement()
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

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	// 直到遇到分号
	for p.curTokenIs(token.SEMICOLON) {
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
		p.peekError(t)
		return false
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{
		Token: p.curToken,
	}
	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// registerPrefix 注册前缀处理函数的映射
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfix 注册中缀处理函数的映射
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// parseExpressionStatement
func (p *Parser) parseExpressionStatement() ast.Statement {
	defer untrace(trace("parseExpressionStatement"))
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// parseExpression 检查前缀位置是否有对应类型的解析函数, 传入参数为优先级
func (p *Parser) parseExpression(precedence int) ast.Expression {
	defer untrace(trace("parseExpression"))

	// 获取前缀解析函数
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	// 没有到语句结束, 并且当前字符优先级低于下一个字符优先级
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		// 获取中缀对应的解析函数, 要知道中缀函数映射的map中条件是类型, 也就是说只有当前语法单元是我们所定义的如"<"">"等才会调用函数
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		// 更新leftExp为中缀解析的结果
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

// parseIntegerLiteral 解析整数
func (p *Parser) parseIntegerLiteral() ast.Expression {
	defer untrace(trace("parseIntegerLiteral"))
	lit := &ast.IntegerLiteral{
		Token: p.curToken,
	}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("cound not parse %q as integer.", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value

	return lit
}

// noPrefixParseFnError 没有注册前缀解析函数
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefox parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// parsePrefixExpression 遇到 "!"和"-"执行此函数, 将其写入expression中, 并且加上其所对应的优先级
func (p *Parser) parsePrefixExpression() ast.Expression {
	defer untrace(trace("parsePrefixExpression"))
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// peekPrecedence 返回 peekToken 也就是下一个字符单元的优先级, 如果没有设置, 则默认返回 LOWEST
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

// curPrecedence 返回 curToken 也就是当前字符单元的优先级, 同样如果没有设置, 返回最低优先级 LOWEST
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

// parseInfixExpression 中缀解析函数
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	defer untrace(trace("parseInfixExpression"))
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// parseBoolean 布尔值解析函数
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
	}
}

// parseGroupedExpression 括号解析函数
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

// parseIfExpression if语句解析
func (p *Parser) parseIfExpression() ast.Expression {
	// 如果p.peekToken不是预期的类型，那么expectPeek会向语法分析器添加错误；如果是预期的类型，则expectPeek将通过调用nextToken方法来前移词法单元。
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}
	return expression
}

// parseBlockStatement 花括号区块解析
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

// parseFunctionLiteral 函数区块解析
func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	lit.Body = p.parseBlockStatement()
	return lit
}

// parseFunctionParameters 解析函数的入参
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var identifiers []*ast.Identifier
	// 如果下一个词法单元是")"则跳过并返回
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}
	// 跳过当前"(", 开始读取第一个参数
	p.nextToken()
	// 将第一个参数初始化为 Identifier 并加入切片中
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)
	// 如果下个词法单元是","就一直循环
	for p.peekTokenIs(token.COMMA) {
		// 跳过两个词法单元, 到下一个参数
		p.nextToken()
		p.nextToken()
		// 将这个参数加入到返回切片中
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}
	// 如果还没有遇到")"就是出错了, 返回nil
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	// 这里肯定是遇到了")",参数构建完毕, 返回切片
	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArgusments()
	return exp
}

func (p *Parser) parseCallArgusments() []ast.Expression {
	var args []ast.Expression

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return args
}
