package goast

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"os"
)

type Step struct {
	Package            ast.Package  // 包名
	Func               ast.FuncDecl // 函数
	Stmt               ast.Stmt     // 函数语句
	Comment            string       // 语句注释
	CallHierarchyDepth int          // 调用深度

	Position  token.Position //Deprecated，直接使用stmt.Pos/End就可以
	Statement string         //Deprecated，这个也可以去掉，后续再格式化处理
	Caller    string         //Deprecated，这个可以后续再解析 caller, $package.$function

	Typ     string   //Deprecated，这个可以后续再解析 receiver type
	X       string   //Deprecated，这个可以后续再解析 receiver variable, or package name
	Seletor string   //Deprecated，这个可以后续再解析 member function or package exported function
	Args    []string //Deprecated，这个可以后续再解析 function args
}

func RenderXSelStmt(fset *token.FileSet, stmt ast.Stmt, depth int, buf *bytes.Buffer) error {

	pos := fset.Position(stmt.Pos()).String()

	// 提取callexpr语句
	var callExprStmt *ast.CallExpr
	switch v := stmt.(type) {
	case *ast.ExprStmt: // 包导出函数
		expr, ok := v.X.(*ast.CallExpr)
		if !ok {
			return errors.New("not *ast.CallExpr")
		}
		callExprStmt = expr
	case *ast.AssignStmt: // 赋值运算符右边操作数检测
		expr, ok := v.Rhs[0].(*ast.CallExpr)
		if !ok {
			return nil
		}
		callExprStmt = expr
	case *ast.ReturnStmt: // ignore
		fmt.Println("found an return, how to render this? branch?")
		return nil
	default:
		// just ignore
		fmt.Fprintf(os.Stderr, "doesn't have *ast.CallExpr node, pos: %s\n", pos)
		return nil
	}

	// 当前语句中的对象类型或者包名，添加到参与者
	funcDecl, err := FunctionContainsStmt(fset, stmt)
	if err != nil {
		return err
	}

	// 获取 newParticipant，通常为语句所在的函数的类型名，或者语句所在的函数的包名
	pkgOfNewPartiicpant, err := PackageNameContainsStmt(fset, stmt)
	if err != nil {
		return err
	}
	recvOfNewPartipant, _ := MethodReceiverTypeName(funcDecl)

	var xType string
	if len(recvOfNewPartipant) == 0 {
		xType = pkgOfNewPartiicpant
	} else {
		//xType = pkgOfNewPartiicpant + ".(" + recvOfNewPartipant + ")"
		xType = pkgOfNewPartiicpant + "." + recvOfNewPartipant
	}

	newParticipant := xType + "." + funcDecl.Name.Name

	// 获取x.selector语句对应的x（变量名或者包名）
	// x.Name represents receiver variable name (not typename), or package name
	var (
		typ string
		//x       = ""
		selName string
	)

	switch expr := callExprStmt.Fun.(type) {

	case *ast.Ident:
		//NOTE: 包内部定义的非导出函数，直接忽略
		if c := expr.Name[0]; c < 'A' {
			return nil
		}
		typ = pkgOfNewPartiicpant

	case *ast.SelectorExpr:
		selName = expr.Sel.Name

		typ, err = ReceiverTypeOrPackageName(fset, expr)
		if err != nil {
			return err
		}
	}

	if _, ok := renderedParticipants[newParticipant]; !ok {
		fmt.Fprintf(buf, "participant \"%s\"\n", newParticipant)
		renderedParticipants[newParticipant] = true
	}

	caller, err := FunctionFullName(fset, funcDecl)
	if err != nil {
		return err
	}

	participantCaller := caller

	source, _ := PosToString(fset, stmt.Pos(), stmt.End())
	fmt.Fprintf(buf, "\"%s\"->\"%s\" : %s\n", participantCaller, newParticipant, joinNewLine(source))
	fmt.Fprintf(buf, "activate \"%s\"\n", newParticipant)
	if newParticipant != participantCaller {
		fmt.Fprintf(buf, "\"%s\" -> \"%s\"\n", newParticipant, participantCaller)
	}
	fmt.Fprintf(buf, "deactivate \"%s\"\n", newParticipant)

	// TODO recursively expand the body at function
	pkgs := gPackages
	if len(typ) == 0 || len(selName) == 0 {
		return nil
	}

	findNode := FindMethod(pkgs, typ, selName)
	if findNode != nil {
		// recursive expand function body
		fmt.Println(typ, selName)
		dat, err := RenderFunctionWithPlantUML(findNode, fset, pkgs)
		if err != nil {
			return err
		}
		buf.Write(dat.Bytes())
	}

	fmt.Printf("not found funcNode, %s.%s\n", typ, selName)
	return nil
}

func Assert(condition bool, errmsg string) {
	if !condition {
		panic("MUST: " + errmsg)
	}
}
