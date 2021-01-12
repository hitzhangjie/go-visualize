package goast

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/token"
)

// RenderForStmt 渲染forstmt，在console中显示
func RenderForStmt(fset *token.FileSet, stmt *ast.ForStmt) error {
	if stmt == nil {
		return errors.New("nil *ast.ForStmt")
	}

	fmt.Printf("for\n")

	// init
	s, err := PosToString(fset, stmt.Init.Pos(), stmt.Init.End())
	if err != nil {
		return err
	}
	fmt.Printf("\tinit: %s\n", s)

	// cond
	s, err = PosToString(fset, stmt.Cond.Pos(), stmt.Cond.End())
	if err != nil {
		return err
	}
	fmt.Printf("\tcond: %s\n", s)

	s, err = PosToString(fset, stmt.Post.Pos(), stmt.Post.End())
	if err != nil {
		return err
	}
	fmt.Printf("\tpost: %s\n", s)

	//body
	fmt.Printf("\tbody:\n")
	for _, s := range stmt.Body.List {
		s, err := PosToString(fset, s.Pos(), s.End())
		if err != nil {
			return err
		}
		fmt.Printf("\t\tstmt: %s\n", s)
	}
	return nil
}

// RenderForStmtWithPlantUML 渲染forstmt，在plantuml中显示
func RenderForStmtWithPlantUML(fset *token.FileSet, stmt *ast.ForStmt, buf *bytes.Buffer) error {
	if stmt == nil {
		return errors.New("nil *ast.ForStmt")
	}

	participant, err := FunctionNameContainsStmt(fset, stmt)
	if err != nil {
		return err
	}
	if _, ok := renderedParticipants[participant]; !ok {
		fmt.Fprintf(buf, "participant \"%s\"\n", participant)
		renderedParticipants[participant] = true
	}

	fmt.Fprintf(buf, "group ForStmt\n")
	defer fmt.Fprintf(buf, "end\n")

	// cond init
	fmt.Fprintf(buf, "\t\"%s\"->\"%s\"\n", participant, participant)

	if stmt.Init != nil {
		s, err := PosToString(fset, stmt.Init.Pos(), stmt.Init.End())
		if err != nil {
			return err
		}
		fmt.Fprintf(buf, "\tnote right: init %s\n", joinNewLine(s))
	}

	// cond determination
	if stmt.Cond != nil {
		s, err := PosToString(fset, stmt.Cond.Pos(), stmt.Cond.End())
		if err != nil {
			return err
		}
		fmt.Fprintf(buf, "\tloop %s\n", s)
	}

	// body statements
	for _, s := range stmt.Body.List {
		//fmt.Fprintf(buf, "\t\t\"%s\"->\"%s\"\n", participant, participant)
		//s, err := PosToString(fset, s.Pos(), s.End())
		//if err != nil {
		//	return err
		//}
		//fmt.Fprintf(buf, "\t\tnote right: body stmt: %s\n", joinNewLine(s))
		dat, err := RenderStmtWithPlantUML(fset, s)
		if err != nil {
			return err
		}
		buf.Write(dat)
	}

	// cond post
	if stmt.Post != nil {
		fmt.Fprintf(buf, "\t\t\"%s\"->\"%s\"\n", participant, participant)
		s, err := PosToString(fset, stmt.Post.Pos(), stmt.Post.End())
		if err != nil {
			return err
		}
		fmt.Fprintf(buf, "\t\tnote right: post %s\n", joinNewLine(s))
	}

	fmt.Fprintf(buf, "\tend\n")

	return nil
}
