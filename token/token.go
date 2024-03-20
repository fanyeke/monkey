package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// 标识符 + 字面量
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 1343456

	// 运算法
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	// 分隔符
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// 关键字
	TRUE      = "TRUE"
	FALSE     = "FALSE"
	BEGIN     = "BEGIN"
	CALL      = "CALL"
	CONST     = "CONST"
	DO        = "DO"
	END       = "END"
	IF        = "IF"
	ODD       = "ODD"
	PROCEDURE = "RPOCEDURE"
	READ      = "READ"
	THEN      = "THEN"
	VAR       = "VAR"
	WHILE     = "WHILE"
	WRITE     = "WRITE"

	// 比较字符
	EQ     = "=="
	NOT_EQ = "!="
)

var keywords = map[string]TokenType{
	"true":      TRUE,
	"false":     FALSE,
	"if":        IF,
	"begin":     BEGIN,
	"call":      CALL,
	"do":        DO,
	"const":     CONST,
	"end":       END,
	"odd":       ODD,
	"procedure": PROCEDURE,
	"read":      READ,
	"var":       VAR,
	"while":     WHILE,
	"write":     WRITE,
}

// LookupIdent 判断是不是关键字
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
