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
		}
	}
}

// lab1
