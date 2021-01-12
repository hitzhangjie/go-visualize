package goast

import (
	"bytes"
	"go/ast"
	"go/token"
)

// RenderFunctionWithPlantUML 渲染一个函数，在plantuml中显示
func RenderFunctionWithPlantUML(funcDecl *ast.FuncDecl, fset *token.FileSet, pkgs map[string]*ast.Package) (*bytes.Buffer, error) {

	buf := bytes.Buffer{}

	for _, stmt := range funcDecl.Body.List {
		dat, err := RenderStmtWithPlantUML(fset, stmt)
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

// Deprecated don't use this method, it's ugly
//func RenderFunction(funcDecl *ast.FuncDecl, fset *token.FileSet,
//	pkgs map[string]*ast.Package, depth int) (string, []Step) {
//
//	depth++
//	steps := []Step{}
//
//	// 首先打印当前service method的全名 $receiverType.$methodName
//	funcFullName, err := FunctionFullName(fset, funcDecl)
//	if err != nil {
//		panic(err)
//	}
//
//	// 然后继续分析函数体中的语句列表，重点分析以下形式的语句：
//	// - x.Sel，包含了对象之间的通信、包之间的通信情景
//	// - 分支控制：ifstmt、switchstmt、forstmt
//
//	for _, stmt := range funcDecl.Body.List {
//		v, err := processStmt(fset, pkgs, funcFullName, stmt, depth)
//		if err != nil {
//			panic(err)
//		}
//		steps = append(steps, v...)
//	}
//
//	return funcFullName, steps
//}

//func processStmt(fset *token.FileSet, pkgs map[string]*ast.Package,
//	funcFullName string, stmt ast.Stmt, depth int) ([]Step, error) {
//
//	steps := []Step{}
//
//	var (
//		pos       = fset.Position(stmt.Pos())
//		args      = "..."
//		statement = ""
//		typ       = "unknown"
//		selName   = ""
//		callExpr  *ast.CallExpr
//		err       error
//	)
//
//	switch v := stmt.(type) {
//	case *ast.ExprStmt: // 包导出函数
//		call, ok := v.X.(*ast.CallExpr)
//		if !ok {
//			return nil, errors.New("not *ast.CallExpr")
//		}
//		callExpr = call
//	case *ast.AssignStmt: // 赋值运算符右边操作数检测
//		call, ok := v.Rhs[0].(*ast.CallExpr)
//		if !ok {
//			// 这种情况下，没有要处理的语句，直接跳过吧
//			return nil, nil
//		}
//		callExpr = call
//	case *ast.IfStmt:
//		// TODO
//		x := Step{
//			Stmt:               v,
//			CallHierarchyDepth: depth,
//			Position:           pos,
//			Statement:          "",
//			Caller:             funcFullName,
//		}
//		return []Step{x}, nil
//	case *ast.SwitchStmt:
//		// TODO
//		x := Step{
//			Stmt:               v,
//			CallHierarchyDepth: depth,
//			Position:           pos,
//			Statement:          "",
//			Caller:             funcFullName,
//		}
//		return []Step{x}, nil
//	case *ast.ForStmt:
//		// TODO
//		x := Step{
//			Stmt:               v,
//			CallHierarchyDepth: depth,
//			Position:           pos,
//			Statement:          "",
//			Caller:             funcFullName,
//		}
//		return []Step{x}, nil
//	default:
//		return nil, nil
//	}
//
//	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
//	if !ok {
//		return nil, errors.New("not *ast.SelectorExpr")
//	}
//
//	// x.Name represents receiver variable, or package name
//	x, ok := selectorExpr.X.(*ast.Ident)
//	if !ok {
//		// TODO 这里有些类型忽略掉了，比如表达式类型*go/ast.Expr (like httpReq.Header...)
//		return nil, errors.New("not *ast.Ident")
//	}
//	selName = selectorExpr.Sel.Name
//
//	if x.Obj != nil {
//		// x.Obj != nil, it's methodName
//		spew.Dump(selectorExpr)
//
//		// BUG: panic: interface conversion: interface {} is *ast.ValueSpec, not *ast.AssignStmt
//		if obj := selectorExpr.X.(*ast.Ident).Obj; obj != nil {
//
//			if decl := obj.Decl; decl != nil {
//				switch v := decl.(type) {
//				case *ast.AssignStmt: // `a := 0`
//					rhs := v.Rhs
//					typ = rhs[0].(*ast.UnaryExpr).X.(*ast.CompositeLit).Type.(*ast.Ident).Name
//					if op := rhs[0].(*ast.UnaryExpr).Op.String(); len(op) != 0 {
//						if op == "&" {
//							op = "*"
//						}
//						typ = op + typ
//					}
//				case *ast.DeclStmt: // `var a = 1` or `var a int` or `var a int = 1`
//					decl, ok := v.Decl.(*ast.GenDecl)
//					if !ok {
//						panic("not *ast.GenDecl")
//					}
//				NextSpec:
//					for _, spec := range decl.Specs {
//						valspec, ok := spec.(*ast.ValueSpec)
//						if !ok {
//							panic("not *ast.ValueSpec")
//						}
//						for _, name := range valspec.Names {
//							if name.Name == obj.Name {
//								switch star := valspec.Type.(type) {
//								case *ast.Ident:
//									typ = valspec.Type.(*ast.Ident).Name
//								case *ast.StarExpr:
//									typ, err = PosToString(fset, star.Pos(), star.End())
//									if err != nil {
//										panic(err)
//									}
//								default:
//									panic("not supported *ast.DeclStmt")
//								}
//								break NextSpec
//							}
//						}
//					}
//
//				default:
//					panic("what's the type")
//				}
//
//			} else {
//				panic("nil *obj.Decl")
//			}
//		}
//	} else {
//		// x.Obj == nil, it's package exported function
//		typ = selectorExpr.X.(*ast.Ident).Name
//	}
//
//	if len(typ) != 0 {
//		statement = fmt.Sprintf("%s %s.%s(%s)\n", pos, x.Name, selName, args)
//	} else {
//		statement = fmt.Sprintf("%s %s.%s(%s)\n", pos, x.Name, selName, args)
//	}
//
//	// TODO recursively expand the body at function
//	findNode := FindMethod(pkgs, typ, selName)
//	if len(statement) != 0 {
//		comment := ""
//		if findNode != nil {
//			//fmt.Println("comment:", findNode.Doc.Text())
//			tmp := strings.TrimSpace(findNode.Doc.Text())
//			comment = strings.TrimPrefix(tmp, findNode.Name.Name)
//		}
//
//		steps = append(steps, Step{
//			Position:           pos,
//			Statement:          statement,
//			Comment:            comment,
//			Caller:             funcFullName,
//			CallHierarchyDepth: depth,
//			X:                  x.Name,
//			Typ:                typ,
//			Seletor:            selName,
//			Args:               []string{args},
//		})
//	}
//
//	if findNode != nil {
//		// recursive expand function body
//		_, nestedSteps := inspectServiceMethod(findNode, fset, pkgs, depth)
//		steps = append(steps, nestedSteps...)
//	} else {
//		fmt.Printf("not found funcNode, %s.%s\n", typ, selName)
//	}
//	return steps, nil
//}

// inspectServiceMethod analyze function code flow, like rpc call hierarchy, control flow, etc
//
// what should we visualize?
// - OOP communication, this depicts the dependencies btw different components
//
// 	 case1 : communication btw components by calling obj's method
// 	 ss := &student{}
// 	 ss.Name()
//
// 	 case2: communication btw components by calling pkg's exported function
// 	 pkg.Statement()
//
// - TODO control flow, if, for, switch, this depicts some important logic
// - TODO concurrency like go func(), wg.Wait()
//
//func inspectServiceMethod(method *ast.FuncDecl, fset *token.FileSet,
//	pkgs map[string]*ast.Package, depth int) (string, []Step) {
//
//	depth++
//	steps := []Step{}
//
//	// TODO methodFullName没有全部统一成一种形式，如pkg.processStmt，或者receivertype.xxx的形式？
//	methodFullName := method.Name.Name
//
//	// 首先打印当前service method的全名 $receiverType.$methodName
//	service, _ := MethodReceiverTypeName(method)
//	if len(service) != 0 {
//		methodFullName = fmt.Sprintf("%s.%s", service, method.Name.Name)
//	}
//
//	// 然后继续分析service method中的语句列表，对符合x.Sel的ast语句结构进行分析
//	for _, stmt := range method.Body.List {
//
//		//spew.Dump(stmt)
//
//		// TODO arguments
//		var (
//			pos       = fset.Position(stmt.Pos())
//			args      = "..."
//			statement = ""
//			typ       = "unknown"
//			selName   = ""
//			callExpr  *ast.CallExpr
//		)
//
//		switch v := stmt.(type) {
//		case *ast.ExprStmt: // 包导出函数
//			call, ok := v.X.(*ast.CallExpr)
//			if !ok {
//				continue
//			}
//			callExpr = call
//		case *ast.AssignStmt: // 赋值运算符右边操作数检测
//			call, ok := v.Rhs[0].(*ast.CallExpr)
//			if !ok {
//				continue
//			}
//			callExpr = call
//		default:
//			continue
//		}
//
//		selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
//		if !ok {
//			continue
//		}
//
//		// x.Name represents receiver variable, or package name
//		x, ok := selectorExpr.X.(*ast.Ident)
//		if !ok {
//			// TODO 这里有些类型忽略掉了，比如表达式类型*go/ast.Expr (like httpReq.Header...)
//			continue
//		}
//		selName = selectorExpr.Sel.Name
//
//		if x.Obj != nil {
//			// x.Obj != nil, it's methodName
//			spew.Dump(selectorExpr)
//
//			// BUG: TODO panic: interface conversion: interface {} is *ast.ValueSpec, not *ast.AssignStmt
//			// 这里处理变量初始化  a = b? 可以忽略？我们只想处理selector.doSomething()这种形式的
//			//rhs := selectorExpr.X.(*ast.Ident).Obj.Decl.(*ast.AssignStmt).Rhs
//			if obj := selectorExpr.X.(*ast.Ident).Obj; obj != nil {
//				if decl := obj.Decl; decl != nil {
//					assign, ok := decl.(*ast.AssignStmt)
//					if ok {
//						rhs := assign.Rhs
//						typ = rhs[0].(*ast.UnaryExpr).X.(*ast.CompositeLit).Type.(*ast.Ident).Name
//						if op := rhs[0].(*ast.UnaryExpr).Op.String(); len(op) != 0 {
//							if op == "&" {
//								op = "*"
//							}
//							typ = op + typ
//						}
//					}
//
//				}
//			}
//		} else {
//			// x.Obj == nil, it's package exported function
//			typ = selectorExpr.X.(*ast.Ident).Name
//		}
//
//		if len(typ) != 0 {
//			statement = fmt.Sprintf("%s%s%s (%s)%s.%s(%s)",
//				log.COLOR_GREEN, pos, log.COLOR_RESET, typ, x.Name, selName, args)
//		} else {
//			statement = fmt.Sprintf("%s%s%s %s.%s(%s)\n",
//				log.COLOR_GREEN, pos, log.COLOR_RESET, x.Name, selName, args)
//		}
//
//		// TODO recursively expand the body at function
//		findNode := FindMethod(pkgs, typ, selName)
//		if len(statement) != 0 {
//			comment := ""
//			if findNode != nil {
//				//fmt.Println("comment:", findNode.Doc.Text())
//				tmp := strings.TrimSpace(findNode.Doc.Text())
//				comment = strings.TrimPrefix(tmp, findNode.Name.Name)
//			}
//
//			//TODO 想办法分析下标准库中的！如fmt.Printf
//			//fmt.Println("node not found")
//
//			steps = append(steps, Step{
//				Position:  pos,
//				Statement: statement,
//				Comment:   comment,
//				// TODO methodFullName没有全部统一成一种形式，如pkg.processStmt，或者receivertype.xxx的形式？
//				Caller:             methodFullName,
//				CallHierarchyDepth: depth,
//				X:                  x.Name,
//				Typ:                typ,
//				Seletor:            selName,
//				Args:               []string{args},
//			})
//		}
//
//		if findNode != nil {
//			//fmt.Printf("found funcNode, %s.%s, %+v\n", typ, selName, findNode)
//			// recursive expand function body
//			_, nestedSteps := inspectServiceMethod(findNode, fset, pkgs, depth)
//			steps = append(steps, nestedSteps...)
//		} else {
//			// TODO 先临时不分析标准库代码
//			if isStandardPackage(typ) {
//				continue
//			}
//			fmt.Printf("not found funcNode, %s.%s\n", typ, selName)
//		}
//	}
//
//	return methodFullName, steps
//}

//func isStandardPackage(pkg string) bool {
//	stds := []string{"fmt", "json"}
//	for _, std := range stds {
//		if std == pkg {
//			return true
//		}
//	}
//	return false
//}
