# monkey 语言解释器
## 1 项目简介
monkey是一个由go语言实现的一个解释器，初衷是完成课程作业，主要的逻辑在主分支，作业的实现在lab分支中，lab3/4为完全体。
效果展示：
```
test case 10:
((8 * (5 - 3)) + 2) / 2

(((8 * (5 - 3)) + 2) / 2)
9
```
## 2 实验环境

- 操作系统：Windows11
- 开发语言：Go 1.21.3
- 开发平台：Goland

## 3 实验内容

### 3.1 项目目录

因为项目本身比较庞大，实现了很多拓展功能，这里只列出了跟实验相关的主要文件

```
│  input.txt
│  json_out.txt
│  main.go
│  txt_out.txt
├─ast
│      ast.go
├─evaluator
│      evaluator.go
├─input
│      input.json
├─lexer
│      lexer.go
├─object
│      environment.go
│      object.go
├─parser
│      parser.go
├─repl
│      repl.go
└─token
        token.go
```

- `token`: 定义关键字和各类符号
- `lexer`: 词法分析
- `input`: 输入文件夹
- `repl`: 实现命令行输入
- `ast`: 语法分析树
- `evaluator`: 解析语法树求值
- `object`: 包装基本类型
- `parser`: 语法分析

### 3.2 快速使用

运行程序即自动执行对文本文件和JSON文件输入的解析, 然后进入命令行输入模式.

**使用**

使用go run main.go或者运行已经编译好的二进制文件./main

```
go run main.go
./main
```

**文本文件输入**

主目录下的input.txt支持输入一个源程序，并将结果输出到主目录下的txt_out.txt中

> 因为Go对JSON兼容更加简单，处理文本文件可能需要解决类似于TCP传输的粘包问题，所以使用JSON进行批量输入

**JSON文件输入**

input文件夹下的input.json支持输入多个源程序，并将结果输出到主目录下的json_out.txt中

**命令行输入**

输入一行回车将立即分析这一行的标识符数量

### 3.3 token 关键字定义

在`token/token.go`中定义有PL/0规则下的标识符和关键字，这里列出部分：

```go
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// 标识符 + 字面量
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 1343456
	// 运算符
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	// 关键字
	TRUE      = "TRUE"
	FALSE     = "FALSE"
	BEGIN     = "BEGIN"
	CALL      = "CALL"
	// [...]
)
```

其中的`ILLEGAL` 、`EOF`、 `IDNET`、 `INT`是辅助词法分析和显示的常量。

定义了一个结构体`Token`包装一个字符作为整体，`Token`中包含两个字段，分别是`Type`代表类型，`Literal`作为字面值：

```go
type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}
```

另外，定义了一个map维护实际关键字和类型的映射关系，一个LookupIdent函数判断是不是关键字。

### 3.4 lexer 字符解析器

`lexer`的核心结构体为：

```go
type Lexer struct {
	input        string
	position     int  // 输入字符的当前位置
	readPosition int  // 输入字符的当前位置下一个，也就是下一个读取的位置
	ch           byte // 当前正在读取的字符
}
```

其中`input`为输入的全部字符，`position`和`readPosition`实际上相当于是一个双指针，为的是方便往前看下一个字符是什么，从而做出一些例如判断“!=”这种符号需要；`ch`记录当前的字符。

`lexer`的核心方式是往下读取一个token，逻辑实现起来可以分为三部分：

1. 跳过空格
2. 枚举可能的单个符号，对于双字符的符号特别判断
3. default分支判断整个token是整数/标识符/关键字

default分支基本：

![image](https://github.com/fanyeke/monkey/assets/113885744/d0077a2f-87c9-488d-a1f3-f0d06d59a2ed)


```
// NextToken 读取一下个token,可以理解为把读取的单个字符加工包装上类型
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	// 跳过空格,换行等
	l.skipWhitespace()

	// case中枚举单个的token,例如+-=*(){}等
	// default中区分关键字/标识符,鉴别数字和错误处理
	switch l.ch {
	// 运算符
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{
				Type:    token.EQ,
				Literal: literal,
			}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{
				Type:    token.NOT_EQ,
				Literal: literal,
			}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	// ... 此处省略一些单个符号的判断
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) && !isDigit(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readIdentifier()
			// 字母开头标识符非法
			for _, v := range tok.Literal {
				if !isDigit(byte(v)) {
					tok.Type = token.ILLEGAL
					return tok
				}
			}
			// 全是数字则为整数
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

```

### 3.5 input 文件输入

读取输入用的是go语言的标准库，输入格式按照下面的进行输入或者修改主目录下的input.txt也可以进行单个程序的解析。

```
{
  "input": [
    "Const num=100;\nVar a1,b2;\nBegin\nRead(A1);\nb2:=a1+num;\nwrite(A1,B2);\nEnd.",
    "_nums;\nn_nums;\nnums_;",
    "555yy\n444;\nyy6",
    "a 3 b 1 a;",
    "",
    "begin \n    nihao\n    wobuzhidaowoyaoshuoshenme\n    danshi\n    fengkuangxinqisi\n    v\n    wo\n    wushikuai\n    rangwochidunhaode\n    xiexie\nend",
    "int main(int argc, char* argv[]){\n    \n    return 0;\n}"
  ]
}
```

主程序和repl中实现的是文件的读写、标识符的计数和命令输入等功能，并非本次实验的核心功能。
### 3.6 ast 语法分析树

本项目是在实验1和实验2的基础上进行拓展的，因此不对词法分析做介绍，直接介绍语法分析。

语法分析使用的是普拉特解析方法，这种方法要求我们建立出一颗语法分洗树，例如(1+2)*3这个表达式，我们需要建立出这样一棵语法分析树（省略了括号的表示）：

![image](https://github.com/fanyeke/monkey/assets/113885744/328510b1-31ec-40fb-9e4d-83cce22d113b)


目的是为了确定语句的指向关系和执行关系，我们可以先对ast进行分析，首先对于一棵树，我们需要对树的节点进行定义，不同的节点有不同的数据结构，比如：对于var关键字，它就需要至少left表示变量和right表示表达式；对于运算符，left和right都是表达式。针对于不同的节点，我们可以抽象出来两类节点接口，分别是语句和表达式，语句不会产生值，表达式会产生值，由此我们声明两个接口

```go
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

```

在ast树上的节点都是基于这两个接口进行拓展实现的，以中缀解析节点为例，结构体各部分的含义如注释所示，实现了三个方法：expressionNode空方法仅用于实现接口（空方法是不会被调用的，在go中目的是符合接口可以方便断言），TokenLiteral用于调试，String方便输出和比较。此外还有很多类型，就不一一列举了。

```go
// InfixExpression 中缀解析结构体
type InfixExpression struct {
	Token    token.Token// token值
	Left     Expression// 指向左表达式
	Operator string// 表达式符号
	Right    Expression// 指向右表达式
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

```

### 3.7 语法解析

了解了语法树的构成以后，语法解析的核心叫比较好理解了：通过判断不同类型，选择不同的节点并加入语法树，最终形成一棵方便于求值的ast语法树。

语法解析还有一个值得注意的部分表达式例如算数运算，由于有不同的优先级，读取的顺序需要高优先级节点是低优先级节点的祖辈，这样才能有正确的计算顺序，因此我们需要针对不同的解析方式设计不同的方法，分为前缀解析和中缀解析两类，还要有配套的优先级判定方法。

```go
// 设定优先级
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       //  array[index]
)

```

```go
// Parser token解析器-词法单元解析器
// Parser 相当于 Lexer 字符解析器的再一层等装, 目的在于模块化层次化项目
// Parser 包含这样几部分: 1. Lexer 2. 前缀中缀解析映射函数 3. 当前和下一个token 4. 错误信息
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
```

```go
// ParseProgram 普拉特解析方法入口
func (p *Parser) ParseProgram() *ast.Program {
	// 构造AST的根节点
	program := &ast.Program{}
	program.Statements = []ast.Statement{} // 因为ast没有设计New方法, 这里对ast进行初始化

	// 遍历输入每个词法单元, 直到文件结尾 token.EOF
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()// 这里对应这边不同的类型应该执行的函数
		// 如果不是 nil 就插入到根节点的 program.Statements 切片中
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

```

入口方法主要在 p.parseStatement()中会把不同的token类型分到不同的函数当中去，这里我们还是来介绍表达式的建立逻辑（以+-*/为例）：

```go
// 函数运行到这里代表着当前字符已经是+-*/
// 函数传入的是加号左边的全部表达式
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
    defer untrace(trace("parseInfixExpression"))// 追踪语句可以忽略
    expression := &ast.InfixExpression{// 新建一个节点，把left的表达式挂上
       Token:    p.curToken,
       Operator: p.curToken.Literal,
       Left:     left,
    }

    precedence := p.curPrecedence()// 获得当前符号的优先级
    p.nextToken()
    expression.Right = p.parseExpression(precedence)// 可以理解为递归获得右边的表达式，并将其挂在right上

    return expression// 最后返回这个节点
}
```

最后以(1+2)*3的解析为例， 可以看一下整个的解析过程：

![image](https://github.com/fanyeke/monkey/assets/113885744/20b847b1-75d7-49d7-8cf6-d25ca1187921)


### 3.8 求值

树构建成功之后解析就比较简单了，用dfs遍历下去然后一步一步返回即可，核心函数是evaluator.go中的Eval：

```go
// Eval 树遍历解释器, Eval 将 ast.Node 作为输入并返回一个 object.Object
// 不同ast节点的求值方式不同
func Eval(node ast.Node, env *object.Environment) object.Object {
    switch node := node.(type) {
    /*
       1. 语句和表达式的本质都是沿着ast树往下递归
    */
    case *ast.Program:
       // 传入语句下挂的每一项
       return evalProgram(node, env) // 语句
    case *ast.ExpressionStatement:
       return Eval(node.Expression, env) // 表达式
    case *ast.IntegerLiteral:
       return &object.Integer{Value: node.Value} // 整数字面量
    case *ast.Boolean:
       return nativeBoolToBooleanObject(node.Value) // 布尔字面量
    case *ast.PrefixExpression: // 前缀表达式
       right := Eval(node.Right, env) // 前缀表达式的右部, 只可能是整数或者布尔, 获得右部以开始解析
       if isError(right) {
          return right
       }
       return evalPrefixExpression(node.Operator, right)
    case *ast.InfixExpression: // 中缀表达式
       left := Eval(node.Left, env)
       if isError(left) {
          return left
       }
       right := Eval(node.Right, env)
       if isError(right) {
          return right
       }
       return evalInfixExpression(node.Operator, left, right)
    // [...]
    return nil
}
```

### 3.9 其他

**标识符存储变量**

实现起来方法比较简单，初始化一个map以标识符为key，值为value即可，实现位于object/environment.go中，因为变量可能出现在ast树上的任何位置，所有采用把存储变量的env作为参数传递到树中，当然设置一个全局变量也是可以的。

**内置函数**

实现了一些内置函数，实现逻辑主要是利用一个映射map，检测到标识符会先看看是不是是不是内置函数，如果是就执行函数的逻辑，值得一提的是代码中变量名的判定在内置函数之前，也就是说我们设定一个内置函数len，依旧可以重新定义一个名为len的变量，在此次运行中将不会调用len函数。
### 3.10 测试样例

```
test case 1:
var a=1;
var b=2;
(a+15)*b;

var a = 1;var b = 2;((a + 15) * b)
32
-------------

test case 2:
1+2+3+4+5;

((((1 + 2) + 3) + 4) + 5)
15
-------------

test case 3:
(1+2)*3;

((1 + 2) * 3)
9
-------------

test case 4:
(3+(1-2))-9/3

((3 + (1 - 2)) - (9 / 3))
-1
-------------

test case 5:
var x=1;
var y=2;
var z=3;(x+y)*z;

var x = 1;var y = 2;var z = 3;((x + y) * z)
9
-------------

test case 6:
(20 / (4 + (2 * 3))) * 2

((20 / (4 + (2 * 3))) * 2)
4
-------------

test case 7:
-1-2

((-1) - 2)
-3
-------------

test case 8:
2+3*5

(2 + (3 * 5))
17
-------------

test case 9:
((3 + 4) * 2) / 5

(((3 + 4) * 2) / 5)
2
-------------

test case 10:
((8 * (5 - 3)) + 2) / 2

(((8 * (5 - 3)) + 2) / 2)
9
-------------
```

|      |
| ---- |
