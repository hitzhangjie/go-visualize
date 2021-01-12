package goast

import (
	"bytes"
	"go/ast"
	"go/token"
)

// RenderFunction 渲染一个函数，在plantuml中显示
func RenderFunction(funcDecl *ast.FuncDecl, fset *token.FileSet, pkgs map[string]*ast.Package) (*bytes.Buffer, error) {

	buf := bytes.Buffer{}

	for _, stmt := range funcDecl.Body.List {
		dat, err := RenderStmt(fset, stmt)
		if err == ErrIgnoreStmt {
			continue
		}
		if err != nil {
			return nil, err
		}
		buf.Write(dat)
	}

	return &buf, nil
}
