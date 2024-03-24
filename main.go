package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fanyeke/monkey/lexer"
	"github.com/fanyeke/monkey/repl"
	"github.com/fanyeke/monkey/token"
	"os"
	"strings"
)

type testJSON struct {
	TestCases []string `json:"testCases"`
}

func main() {
	readFromTxt()
	readFromJSON()
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}

func readFromJSON() {
	test := testJSON{}
	file, err := os.ReadFile("input/input.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &test)
	if err != nil {
		panic(err)
	}

	out := bytes.Buffer{}
	for i, testCase := range test.TestCases {
		out.WriteString(fmt.Sprintf("test case %d:\n", i+1))
		out.WriteString(testCase + "\n\n")
		l := lexer.New(testCase)
		out.WriteString(fmt.Sprintf("Type\t\tLiteral\n"))
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			out.WriteString(fmt.Sprintf("%v\t\t%v\n", tok.Type, tok.Literal))
		}
		out.WriteString("-------------\n\n")
	}
	err = os.WriteFile("json_out.txt", out.Bytes(), 0775)
	if err != nil {
		panic(err)
	}
}

func readFromTxt() {
	out := bytes.Buffer{}
	// 读取文本文件
	file, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	input := string(file)
	input = strings.ToLower(input)
	out.WriteString("Test Case 1:\n")
	out.WriteString(input + "\n\n")
	l := lexer.New(input)
	out.WriteString(fmt.Sprintf("Type\t\tLiteral\n"))
	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		out.WriteString(fmt.Sprintf("%v\t\t%v\n", tok.Type, tok.Literal))
	}
	err = os.WriteFile("txt_out.txt", out.Bytes(), 0775)
	if err != nil {
		panic(err)
	}
}
