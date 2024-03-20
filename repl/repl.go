package repl

import (
	"bufio"
	"fmt"
	"github.com/fanyeke/monkey/lexer"
	"github.com/fanyeke/monkey/token"
	"io"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		identMap := make(map[string]int)
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		fmt.Println("Type\t\tLiteral")
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%v\t\t%v\n", tok.Type, tok.Literal)
			// 标识符计数
			if tok.Type == token.IDENT {
				if _, ok := identMap[tok.Literal]; !ok {
					identMap[tok.Literal] = 1
				} else {
					identMap[tok.Literal]++
				}
			}
		}
		fmt.Println(identMap)
	}
}
