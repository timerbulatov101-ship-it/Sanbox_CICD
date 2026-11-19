package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestCalculateComplexity(t *testing.T) {
	source := `package sample
	func example(a, b bool, n int) int {
		if a && b { return 1 } else if n > 0 {
			for i := 0; i < n; i++ { println(i) }
		}
		switch n {
		case 1: return 2
		case 2, 3: return 3
		default: return 0
		}
	}`
	file, err := parser.ParseFile(token.NewFileSet(), "sample.go", source, 0)
	if err != nil {
		t.Fatal(err)
	}
	fn := file.Decls[0].(*ast.FuncDecl)
	if got, want := calculateComplexity(fn.Body), 7; got != want {
		t.Fatalf("calculateComplexity() = %d, want %d", got, want)
	}
}
