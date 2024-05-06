package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fanyeke/monkey/evaluator"
	"github.com/fanyeke/monkey/lexer"
	"github.com/fanyeke/monkey/parser"
	"github.com/fanyeke/monkey/repl"
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
	//env := object.NewEnvironment()
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
		p := parser.New(l)
		program := p.ParseProgram()
		out.WriteString(program.String())
		//evaluated := evaluator.Eval(program, env)
		evaluator.TempVarCount = 0

		code := evaluator.GenerateIntermediateCode(program.Statements[0])

		out.WriteString(fmt.Sprintf("\nIntermediate Code:\n%s\n\n", code))
		//if evaluated != nil {
		//	out.WriteString("\n" + evaluated.Inspect() + "\n")
		//}
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
	//env := object.NewEnvironment()

	if err != nil {
		panic(err)
	}
	input := string(file)
	input = strings.ToLower(input)
	out.WriteString("Test Case 1:\n")
	out.WriteString(input + "\n\n")
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	out.WriteString(program.String())
	//evaluated := evaluator.Eval(program, env)
	evaluator.TempVarCount = 0
	code := evaluator.GenerateIntermediateCode(program.Statements[0])
	out.WriteString(fmt.Sprintf("\nIntermediate Code:\n%s\n\n", code))
	err = os.WriteFile("txt_out.txt", out.Bytes(), 0775)
	if err != nil {
		panic(err)
	}
}
