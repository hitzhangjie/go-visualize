package goast

import (
	"bytes"
	"errors"
	"go/ast"
	"go/token"
	"strings"
)

var (
	ErrIgnoreStmt = errors.New("render not supported for this kind of ast.Stmt")
)

var (
	renderedParticipants = map[string]bool{}
)

// RenderStmt render ast.Stmt with plantuml
func RenderStmt(fset *token.FileSet, stmt ast.Stmt) ([]byte, error) {
	var (
		buf = &bytes.Buffer{}
		err error
	)

	switch v := stmt.(type) {
	case *ast.IfStmt:
		err = RenderIfStmt(fset, v, 0, buf)
	case *ast.SwitchStmt:
		err = RenderSwitchStmt(fset, v, buf)
	case *ast.ForStmt:
		err = RenderForStmt(fset, v, buf)
	default:
		err = RenderXSelStmt(fset, v, 0, buf)
		//fmt.Printf("stmt kind not supported: %T, pos: %s\n", v, fset.Position(v.Pos()).String())
	}

	if err != nil {
		return nil, err
	}
	dat := buf.Bytes()

	return dat, nil
}

func joinNewLine(s string) string {
	return strings.ReplaceAll(s, "\n", " ")
}
