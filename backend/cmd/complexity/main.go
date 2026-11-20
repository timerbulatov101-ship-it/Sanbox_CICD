package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type FunctionResult struct {
	Name       string
	File       string
	Line       int
	Complexity int
}

func main() {
	root := flag.String("path", ".", "directory to analyze")
	maxComplexity := flag.Int("max", 10, "maximum allowed cyclomatic complexity")
	flag.Parse()

	results, err := Analyze(*root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "analysis error:", err)
		os.Exit(2)
	}

	violations := 0
	fmt.Printf("%-8s %-8s %-40s %s\n", "COMPLEX", "LINE", "FUNCTION", "FILE")
	for _, result := range results {
		fmt.Printf("%-8d %-8d %-40s %s\n", result.Complexity, result.Line, result.Name, result.File)
		if result.Complexity > *maxComplexity {
			violations++
		}
	}

	fmt.Printf("\nFunctions: %d; violations (complexity > %d): %d\n", len(results), *maxComplexity, violations)
	if violations > 0 {
		os.Exit(1)
	}
}

func Analyze(root string) ([]FunctionResult, error) {
	fset := token.NewFileSet()
	results := make([]FunctionResult, 0)

	err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			if path != root && (info.Name() == "vendor" || strings.HasPrefix(info.Name(), ".")) {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}

		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}

		ast.Inspect(file, func(node ast.Node) bool {
			fn, ok := node.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				return true
			}
			name := fn.Name.Name
			if fn.Recv != nil {
				name = receiverName(fn.Recv) + "." + name
			}
			results = append(results, FunctionResult{
				Name: name, File: path, Line: fset.Position(fn.Pos()).Line,
				Complexity: calculateComplexity(fn.Body),
			})
			return false
		})
		return nil
	})
	return results, err
}

func receiverName(fieldList *ast.FieldList) string {
	if fieldList == nil || len(fieldList.List) == 0 {
		return ""
	}
	switch t := fieldList.List[0].Type.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name
		}
	}
	return "receiver"
}

func calculateComplexity(body *ast.BlockStmt) int {
	complexity := 1
	ast.Inspect(body, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt:
			complexity++
		case *ast.CaseClause:
			if len(n.List) > 0 {
				complexity++
			}
		case *ast.CommClause:
			if n.Comm != nil {
				complexity++
			}
		case *ast.BinaryExpr:
			if n.Op == token.LAND || n.Op == token.LOR {
				complexity++
			}
		}
		return true
	})
	return complexity
}
