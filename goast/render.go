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

// RenderStmt render ast.Stmt
func RenderStmt(fset *token.FileSet, stmt ast.Stmt) error {
	switch v := stmt.(type) {
	case *ast.IfStmt:
		return RenderIfStmt(fset, v, 0)
	case *ast.SwitchStmt:
		return RenderSwitchStmt(fset, v)
	case *ast.ForStmt:
		return RenderForStmt(fset, v)
	default:
		//fmt.Printf("stmt kind not supported: %T, pos: %s\n", v, fset.Position(v.Pos()).String())
		return ErrIgnoreStmt
	}
}

// RenderStmtWithPlantUML render ast.Stmt with plantuml
func RenderStmtWithPlantUML(fset *token.FileSet, stmt ast.Stmt) ([]byte, error) {
	var (
		buf = &bytes.Buffer{}
		err error
	)

	switch v := stmt.(type) {
	case *ast.IfStmt:
		err = RenderIfStmtWithPlantuml(fset, v, 0, buf)
	case *ast.SwitchStmt:
		err = RenderSwitchStmtWithPlantUML(fset, v, buf)
	case *ast.ForStmt:
		err = RenderForStmtWithPlantUML(fset, v, buf)
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
