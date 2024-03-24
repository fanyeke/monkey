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

	LT      = "<"
	GT      = ">"
	LEQ     = "<="
	GEQ     = ">="
	BECOMES = ":="
	REQ     = "#"
	// 分隔符
	COMMA     = ","
	SEMICOLON = ";"
	PERIOD    = "."

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// 关键字
	BEGIN     = "BEGIN"
	CALL      = "CALL"
	CONST     = "CONST"
	DO        = "DO"
	END       = "END"
	IF        = "IF"
	ODD       = "ODD"
	PROCEDURE = "PROCEDURE"
	READ      = "READ"
	THEN      = "THEN"
	VAR       = "VAR"
	WHILE     = "WHILE"
	WRITE     = "WRITE"

	// 比较字符
	EQ     = "=="
	NOT_EQ = "!="

	STRING = "STRING"
)

var keywords = map[string]TokenType{
	"if":        IF,
	"begin":     BEGIN,
	"call":      CALL,
	"const":     CONST,
	"do":        DO,
	"end":       END,
	"odd":       ODD,
	"procedure": PROCEDURE,
	"read":      READ,
	"then":      THEN,
	"while":     WHILE,
	"write":     WRITE,
	"var":       VAR,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
