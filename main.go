package main

import (
	"fmt"
	"github.com/fanyeke/monkey/input"
	"github.com/fanyeke/monkey/lexer"
	"github.com/fanyeke/monkey/queue"
	"github.com/fanyeke/monkey/repl"
	"github.com/fanyeke/monkey/token"
	"os"
	user2 "os/user"
	"strings"
)

func main() {
	user, err := user2.Current()
	if err != nil {
		panic(err)
	}
	fmt.Println("======== lab1 ========")
	fmt.Println("1: command line input")
	fmt.Println("2 json file input")
	fmt.Println("3 txt file input")
	fmt.Println("Please enter 1 - 3:")
	var n int
	fmt.Scan(&n)
	if err != nil {
		panic(err)
	}
	switch n {
	case 1:
		fmt.Printf("Hello %s!This is the Monkey programming language!\n", user.Username)
		fmt.Printf("Feel free to type in commands\n")
		repl.Start(os.Stdin, os.Stdout)
	case 2:
		readInput()
	case 3:
		readTxt()
	}

}

func readTxt() {
	outFile := ""

	inputFile, err := os.ReadFile("input.txt")
	inp := string(inputFile)
	if err != nil {
		panic(err)
	}
	// 取一个程序
	inp = strings.ToLower(inp)
	l := lexer.New(inp)
	fmt.Println("Type\t\tLiteral")
	// 初始化一个队列和计数map
	q := queue.NewQueue()
	identMap := make(map[string]int)
	// 处理取出的这一个程序
	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {

		//fmt.Printf("%v\t\t%v\n", tok.Type, tok.Literal)
		// 标识符计数
		if tok.Type == token.IDENT {
			if _, ok := identMap[tok.Literal]; !ok {
				q.Enqueue(tok.Literal)
				identMap[tok.Literal] = 1
			} else {
				identMap[tok.Literal]++
			}
		}

	}

	// 每一行输出一次
	for !q.IsEmpty() {
		var dq string
		dq = q.Dequeue().(string)
		fmt.Printf("(%s:\t%d)\n", dq, identMap[dq])
		outFile += fmt.Sprintf("(%s:\t%d)\n", dq, identMap[dq])
	}
	err = os.WriteFile("txt_out.txt", []byte(outFile), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully wrote input to out.txt")
}

func readInput() {
	ri := input.Ri
	outFile := ""
	for index, inp := range ri.Input {
		// 结果写入文件
		outFile += fmt.Sprintf("out%d:\n", index+1)
		// 取一个程序
		inp = strings.ToLower(inp)
		l := lexer.New(inp)
		fmt.Println("Type\t\tLiteral")
		// 初始化一个队列和计数map
		q := queue.NewQueue()
		identMap := make(map[string]int)
		// 处理取出的这一个程序
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {

			fmt.Printf("%v\t\t%v\n", tok.Type, tok.Literal)
			// 标识符计数
			if tok.Type == token.IDENT {
				if _, ok := identMap[tok.Literal]; !ok {
					q.Enqueue(tok.Literal)
					identMap[tok.Literal] = 1
				} else {
					identMap[tok.Literal]++
				}
			}

		}

		// 每一行输出一次
		for !q.IsEmpty() {
			var dq string
			dq = q.Dequeue().(string)
			fmt.Printf("(%s:\t%d)\n", dq, identMap[dq])
			outFile += fmt.Sprintf("(%s:\t%d)\n", dq, identMap[dq])
		}
	}
	err := os.WriteFile("json_out.txt", []byte(outFile), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully wrote input to out.txt")
}
