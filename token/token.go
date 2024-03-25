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

	LT  = "<"
	GT  = ">"
	NEQ = "#"
	LEQ = "<="
	GEQ = ">="

	// 分隔符
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// 关键字
	FUNCTION = "FUNCTION"
	LET      = "LET"

	TRUE   = "TRUE"
	FALSE  = "FALSE"
	IF     = "IF" // 保留
	ELSE   = "ELSE"
	RETURN = "RETURN"
	STRING = "STRING"
	// TODO 新加入
	VAR       = "VAR"
	BEGIN     = "BEGIN"
	CALL      = "CALL"
	CONST     = "CONST"
	DO        = "DO"
	END       = "END"
	ODD       = "ODD"
	PROCEDURE = "PROCEDURE"
	READ      = "READ"
	THEN      = "THEN"
	WHILE     = "WHILE"
	WRITE     = "WRITE"

	// 比较字符
	EQ     = "=="
	NOT_EQ = "!="
)

var keywords = map[string]TokenType{
	"fn":        FUNCTION,
	"let":       LET,
	"true":      TRUE,
	"false":     FALSE,
	"if":        IF,
	"else":      ELSE,
	"return":    RETURN,
	"var":       VAR,
	"begin":     BEGIN,
	"call":      CALL,
	"const":     CONST,
	"do":        DO,
	"end":       END,
	"odd":       ODD,
	"procedure": PROCEDURE,
	"read":      READ,
	"while":     WHILE,
	"write":     WRITE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
