package repl

import (
	"bufio"
	"fmt"
	"github.com/fanyeke/monkey/evaluator"
	"github.com/fanyeke/monkey/lexer"
	"github.com/fanyeke/monkey/object"
	"github.com/fanyeke/monkey/parser"
	"io"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	// 输入和变量储存环境
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		// 每次取一行, 每一行都建立一颗 ast 树, 进行解释
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		// 建树
		program := p.ParseProgram()
		fmt.Println(program.String())
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		// 解析求值
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func printParserErrors(out io.Writer, errors []string) {
	//io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
