package goast

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/token"
)

// RenderSwitchStmt 渲染switchstmt，在console中显示
func RenderSwitchStmt(fset *token.FileSet, stmt *ast.SwitchStmt) error {
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

// RenderSwitchStmt 渲染switchstmt，在console中显示
func RenderSwitchStmtWithPlantUML(fset *token.FileSet, stmt *ast.SwitchStmt, buf *bytes.Buffer) error {
	if stmt == nil {
		return errors.New("nil *ast.SwitchStmt")
	}
	participant, err := FunctionNameContainsStmt(fset, stmt)
	if err != nil {
		return err
	}
	if _, ok := renderedParticipants[participant]; !ok {
		fmt.Fprintf(buf, "participant \"%s\"\n", participant)
		renderedParticipants[participant] = true
	}

	// if condition & body
	tag := "tag"
	if ident, ok := stmt.Tag.(*ast.Ident); ok {
		tag = ident.Name
	}

	switchInited := false

	if stmt.Body != nil {
		for _, l := range stmt.Body.List {
			clause, ok := l.(*ast.CaseClause)
			if !ok {
				panic("Assert be *ast.CaseClause")
			}

			// case condition & case body
			n := len(clause.List)

			if n == 0 {
				// n==0为default
				fmt.Fprintf(buf, "else others\n")
			} else {
				// n!=0，有两种情况，alt或者else
				v := clause.List[0]
				s, err := PosToString(fset, v.Pos(), v.End())
				if err != nil {
					return err
				}
				if !switchInited {
					fmt.Fprintf(buf, "alt %s matches cond %s\n", tag, s)
					switchInited = true
				} else {
					fmt.Fprintf(buf, "else %s matches cond %s\n", tag, s)
				}
			}

			for _, l := range clause.Body {
				//s, err := PosToString(fset, l.Pos(), l.End())
				//if err != nil {
				//	return err
				//}
				//fmt.Fprintf(buf, "\t\"%s\"->\"%s\"\n", participant)
				//fmt.Fprintf(buf, "\tnote right: %s\n", joinNewLine(s))

				dat, err := RenderStmtWithPlantUML(fset, l)
				if err != nil {
					return err
				}
				buf.Write(dat)
			}
		}
	}

	if switchInited {
		fmt.Fprintf(buf, "end\n")
	}

	return nil
}
