package goast

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

// RenderStmtWithConsole render ast.Stmt
//
// Deprecated
func RenderStmtWithConsole(fset *token.FileSet, stmt ast.Stmt) error {
	switch v := stmt.(type) {
	case *ast.IfStmt:
		return RenderIfStmtWithConsole(fset, v, 0)
	case *ast.SwitchStmt:
		return RenderSwitchStmtWithConsole(fset, v)
	case *ast.ForStmt:
		return RenderForStmtWithConsole(fset, v)
	default:
		//fmt.Printf("stmt kind not supported: %T, pos: %s\n", v, fset.Position(v.Pos()).String())
		return ErrIgnoreStmt
	}
}

// RenderIfStmtWithConsole 渲染ifstmt语句，在console中显示
//
// Deprecated
func RenderIfStmtWithConsole(fset *token.FileSet, stmt *ast.IfStmt, depth int) error {
	if stmt == nil {
		return errors.New("nil *ast.IfStmt")
	}

	// if condition & body
	s, err := PosToString(fset, stmt.Cond.Pos(), stmt.Cond.End())
	if err != nil {
		return err
	}
	printWithIndent(depth, "if condition: %s", s)

	if stmt.Body != nil {
		for _, l := range stmt.Body.List {
			s, err := PosToString(fset, l.Pos(), l.End())
			if err != nil {
				return err
			}
			printWithIndent(depth, "\tstmt: %s", s)
		}
	}

	// else & body
	if stmt.Else == nil {
		return nil
	}

	printWithIndent(depth, "Else:")

	if blk, ok := stmt.Else.(*ast.BlockStmt); ok && blk != nil {
		for _, l := range blk.List {
			s, err := PosToString(fset, l.Pos(), l.End())
			if err != nil {
				return err
			}
			printWithIndent(depth, "\tstmt: %s", s)
		}
		return nil
	}
	if nestedIfStmt, ok := stmt.Else.(*ast.IfStmt); ok && nestedIfStmt != nil {
		return RenderIfStmtWithConsole(fset, nestedIfStmt, depth+1)
	}
	return errors.New("invalid Else")
}

func printWithIndent(indent int, format string, args ...interface{}) {
	prefix := strings.Repeat("\t", indent)
	fmt.Printf(prefix+format+"\n", args...)
}

// RenderSwitchStmtWithConsole 渲染switchstmt，在console中显示
//
// Deprecated
func RenderSwitchStmtWithConsole(fset *token.FileSet, stmt *ast.SwitchStmt) error {
	if stmt == nil {
		return errors.New("nil *ast.SwitchStmt")
	}

	// tag
	s, err := PosToString(fset, stmt.Tag.Pos(), stmt.Tag.End())
	if err != nil {
		return err
	}
	printWithIndent(0, "switch %v", s)

	if stmt.Body != nil {
		for _, l := range stmt.Body.List {
			clause, ok := l.(*ast.CaseClause)
			if !ok {
				panic("Assert be *ast.CaseClause")
			}

			// case condition & case body
			n := len(clause.List)
			if n == 0 {
				printWithIndent(0, "\tdefault")
			} else {
				v := clause.List[n-1]
				s, err := PosToString(fset, v.Pos(), v.End())
				if err != nil {
					return err
				}
				printWithIndent(0, "\tcase: %s", s)
			}

			for _, l := range clause.Body {
				s, err := PosToString(fset, l.Pos(), l.End())
				if err != nil {
					return err
				}
				printWithIndent(0, "\t\tstmt: %s", s)
			}
		}
	}

	return nil
}

// RenderForStmtWithConsole 渲染forstmt，在console中显示
//
// Deprecated
func RenderForStmtWithConsole(fset *token.FileSet, stmt *ast.ForStmt) error {
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
