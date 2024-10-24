package main

import (
	"carrionlang/evaluator"
	"carrionlang/lexer"
	"carrionlang/object"
	"carrionlang/parser"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: carrionlang [file.crl]")
		return
	}
	filePath := os.Args[1]
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file %v", err)
	}
	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			fmt.Println(msg)
		}
		return
	}
	env := object.NewEnvironment()
	evaluator.DefineBuiltins(env)
	evaluator.Eval(program, env)
}
