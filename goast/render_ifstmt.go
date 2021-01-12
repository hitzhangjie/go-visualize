package goast

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

// RenderIfStmt 渲染ifstmt语句，在plantuml中显示
//
// TODO: if条件，使用条件表达式，代替源码位置
func RenderIfStmt(fset *token.FileSet, stmt *ast.IfStmt, depth int, buf *bytes.Buffer) error {

	var (
		participant string = "ifstmt"
		err         error
	)
	if depth == 0 {
		participant, err = FunctionNameContainsStmt(fset, stmt)
		if err != nil {
			return err
		}
		if _, ok := renderedParticipants[participant]; !ok {
			fmt.Fprintf(buf, "participant \"%s\"\n", participant)
			renderedParticipants[participant] = true
		}
	}

	indent := strings.Repeat("\t", depth)

	if stmt == nil {
		return errors.New("nil *ast.IfStmt")
	}

	// if init
	if stmt.Init != nil {
		s, err := PosToString(fset, stmt.Init.Pos(), stmt.Init.End())
		if err != nil {
			return err
		}
		fmt.Fprintf(buf, "%s\t\"%s\"->\"%s\"\n", indent, participant, participant)
		fmt.Fprintf(buf, "note right: %s\n", joinNewLine(s))
	}

	// if condition & body
	s, err := PosToString(fset, stmt.Cond.Pos(), stmt.Cond.End())
	if err != nil {
		return err
	}
	fmt.Fprintf(buf, "%salt %s\n", indent, s)
	defer fmt.Fprintf(buf, "%send\n", indent)

	if stmt.Body != nil {
		for _, l := range stmt.Body.List {
			fmt.Fprintf(buf, "%s\t\"%s\"->\"%s\"\n", indent, participant, participant)
			s, err := PosToString(fset, l.Pos(), l.End())
			if err != nil {
				return err
			}
			fmt.Fprintf(buf, "%s\tnote right: stmt %s\n", indent, joinNewLine(s))
		}
	}

	// else & body
	if stmt.Else == nil {
		return nil
	}

	// else if ==> IfStmt
	if nestedIfStmt, ok := stmt.Else.(*ast.IfStmt); ok && nestedIfStmt != nil {
		fmt.Fprintf(buf, "%selse others\n", indent)
		return RenderIfStmt(fset, nestedIfStmt, depth+1, buf)
	}

	// else ==> BlockStmt
	if blk, ok := stmt.Else.(*ast.BlockStmt); ok && blk != nil {
		fmt.Fprintf(buf, "%selse others\n", indent)
		for _, l := range blk.List {

			// ifstmt中的body部分语句，需要回到一般化处理
			dat, err := RenderStmt(fset, l)
			if err == ErrIgnoreStmt {
				continue
			}
			if err != nil {
				return err
			}

			buf.Write(dat)
		}
		return nil
	}

	return errors.New("invalid Else")
}
