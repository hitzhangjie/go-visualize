package goast

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	gPackages = map[string]*ast.Package{}
)

// ParseFile 解析源文件返回对应fileset、astfile，如果遇到错误返回error
func ParseFile(file string) (*token.FileSet, *ast.File, error) {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}
	return fset, astFile, nil
}

// ParseDir 解析目录对应package，如果recursive为true则递归解析子目录对应package，
// 返回map结构key为package名，value为对应的*ast.Package语法树结构.
func ParseDir(dir string, recursive bool) (*token.FileSet, map[string]*ast.Package, error) {

	dirs := []string{dir}

	if recursive {
		v, err := traverseSubDirs(dir)
		if err != nil {
			return nil, nil, err
		}
		dirs = append(dirs, v...)
	}

	fset := token.NewFileSet()
	allPkgs := map[string]*ast.Package{}

	for _, dir := range dirs {
		pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
		if err != nil {
			return nil, nil, err
		}
		for k, v := range pkgs {
			allPkgs[k] = v
		}
	}

	gPackages = allPkgs

	return fset, allPkgs, nil
}

func traverseSubDirs(dir string) ([]string, error) {
	dirs := []string{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})
	return dirs, err
}

// Filter 如果希望过滤出该函数，则返回true
type Filter func(fn *ast.FuncDecl) bool

// Functions 解析file中对应的functions，如果filter不为nil，则filter(node)返回true时才返回
func Functions(fset *token.FileSet, file *ast.File, filter Filter) ([]*ast.FuncDecl, error) {
	fnList := []*ast.FuncDecl{}
	ast.Inspect(file, func(node ast.Node) bool {
		if fn, ok := node.(*ast.FuncDecl); ok {
			if filter != nil {
				if filter(fn) {
					fnList = append(fnList, fn)
				}
				return true
			}
			fnList = append(fnList, fn)
		}
		return true
	})
	return fnList, nil
}

// BUG: recvType是否只检查了接受者类型，没有检查package名?
func FindMethod(pkgs map[string]*ast.Package, recvType, funcName string) *ast.FuncDecl {

	// TODO how to seperate *Student and Student
	recvType = strings.TrimPrefix(recvType, "*")

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			filter := func(fn *ast.FuncDecl) bool {
				if _, ok := isTargetMethod(fn, recvType, funcName); ok {
					return true
				}
				if _, ok := isTargetPackageExportedFunction(fn, pkg.Name, recvType, funcName); ok {
					return true
				}
				return false
			}
			fnList, _ := Functions(nil, file, filter) // never return error
			if len(fnList) != 0 {
				return fnList[0]
			}
		}
	}

	return nil
}

// ReceiverTypeOrPackageName the function to analyze, like main.main or main.(*helloworldServiceImpl).Hello
func FindFunction(pkgs map[string]*ast.Package, funcName string) (*ast.FuncDecl, error) {

	// TODO how to seperate *Student and Student
	vals := strings.Split(funcName, ".")
	switch len(vals) {
	case 2:
		return FindPackageExportedFunction(pkgs, vals[0], vals[1])
	case 3:

		return FindServiceExportedFunction(pkgs, vals[0], vals[1], vals[2])
	default:
		return nil, errors.New("invalid function: must be main.main or main.*helloworldServiceImpl.Hello or main.(*helloworldServiceImpl).Hello")
	}
}

// TODO how to seperate *Student and Student
func FindPackageExportedFunction(pkgs map[string]*ast.Package, pkgName, funcName string) (*ast.FuncDecl, error) {

	for name, pkg := range pkgs {
		if name != pkgName {
			continue
		}
		for _, file := range pkg.Files {
			filter := func(fn *ast.FuncDecl) bool {
				if fn.Name.Name != funcName {
					return false
				}
				return true
			}

			fnList, err := Functions(nil, file, filter) // never return error
			if err != nil {
				return nil, err
			}

			switch len(fnList) {
			case 0:
				continue
			case 1:
				return fnList[0], nil
			default:
				return nil, errors.New("code has errors: found more than 1 instances")
			}

			if len(fnList) != 0 {
				return fnList[0], nil
			}
		}
	}

	return nil, errors.New("not found")
}

func FindServiceExportedFunction(pkgs map[string]*ast.Package, pkgName, recvType, funcName string) (*ast.FuncDecl, error) {

	recvType = strings.TrimPrefix(recvType, "(")
	recvType = strings.TrimSuffix(recvType, ")")

	for name, pkg := range pkgs {
		if pkgName != name {
			continue
		}
		for _, file := range pkg.Files {
			filter := func(fn *ast.FuncDecl) bool {
				if fn.Name.Name != funcName {
					return false
				}
				if v, err := MethodReceiverTypeName(fn); err != nil || v != recvType {
					fmt.Println("method recvType:", v)
					return false
				}
				return true
			}

			fnList, err := Functions(nil, file, filter) // never return error
			if err != nil {
				return nil, err
			}

			switch len(fnList) {
			case 0:
				continue
			case 1:
				return fnList[0], nil
			default:
				return nil, errors.New("code has errors: found more than 1 instances")
			}

			if len(fnList) != 0 {
				return fnList[0], nil
			}
		}
	}

	return nil, errors.New("not found")
}

func isTargetMethod(fn *ast.FuncDecl, recvType, funcName string) (*ast.FuncDecl, bool) {

	if len(recvType) != 0 {
		// fn is function, rather than methods
		if fn.Recv == nil || len(fn.Recv.List) == 0 || fn.Recv.List[0] == nil || fn.Recv.List[0].Type == nil {
			return nil, false
		}

		// gorpc template make sure receiver type of methods of generated implemention of service interface
		// always conforms to form `(s *${service}) RPCMethod(ctx, req, rsp) error`.
		typ, ok := fn.Recv.List[0].Type.(*ast.StarExpr)
		if !ok {
			return nil, false
		}

		// filter out the methods whose receiver type has the same type as registered services
		ident, ok := typ.X.(*ast.Ident)
		if !ok || recvType != ident.Name {
			return nil, false
		}
	}

	// filter out the methods whose name not matches
	if funcName != fn.Name.Name {
		return nil, false
	}

	return fn, true
}

func isTargetPackageExportedFunction(fn *ast.FuncDecl, pkg, x, funcName string) (*ast.FuncDecl, bool) {
	if pkg == x && fn.Name.Name == funcName {
		return fn, true
	}
	return nil, false
}

// MethodReceiverTypeName 返回方法接收器类型
func MethodReceiverTypeName(fn *ast.FuncDecl) (string, error) {

	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return "", errors.New("invalid receiver type")
	}

	typ, ok := fn.Recv.List[0].Type.(*ast.StarExpr)
	if !ok {
		return "", errors.New("not *ast.StarExpr")
	}

	ident, ok := typ.X.(*ast.Ident)
	if !ok {
		panic("not *ast.Ident")
	}

	// differentiate value from pointer
	if typ.Star.IsValid() {
		return "*" + ident.Name, nil
	}
	return ident.Name, nil
}

// PackageNameContainsFunc 返回方法所在包名
func PackageNameContainsFunc(fset *token.FileSet, fn *ast.FuncDecl) (string, error) {

	fname := fset.Position(fn.Pos()).Filename
	_, file, err := ParseFile(fname)
	if err != nil {
		return "", err
	}
	return file.Name.Name, nil
}

// Print print AST of `filename`
func Print(filename string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parser.ParseFile error: %v", err)
	}
	return ast.Print(fset, f)
}

// PosToString 将pos转换成对应字符串
//
// token.NoPos==0, otherwise, token.Pos>0
func PosToString(fset *token.FileSet, begin, end token.Pos) (string, error) {
	if !begin.IsValid() || !end.IsValid() || begin >= end {
		return "", errors.New("invalid token.Pos")
	}

	// see https://abhinavg.net/posts/understanding-token-pos/
	base := fset.File(begin).Base()

	position := fset.Position(begin)
	dat, err := ioutil.ReadFile(position.Filename)
	if err != nil {
		return "", err
	}

	n := len(dat)
	if n < int(end)-base-1 {
		return "", errors.New("invalid token.Pos, exceed filesize")
	}

	idx1 := int(begin) - base - 1
	//idx2 := int(end) - base - 1
	idx2 := int(end) - base

	s := string(dat[idx1:idx2])
	return strings.TrimSpace(s), nil
}

// FunctionFullName 返回完整函数名，如果是方法，则返回$receiver.$func，如果是包导出函数，则返回$package.$func
func FunctionFullName(fset *token.FileSet, fn *ast.FuncDecl) (string, error) {

	pkg, err := PackageNameContainsFunc(fset, fn)
	if err != nil {
		return "", errors.New("BUG: Assert find the package")
	}

	// type exported method
	recvType, _ := MethodReceiverTypeName(fn)
	if len(recvType) != 0 {
		//fullName := fmt.Sprintf("%s.(%s).%s", pkg, recvType, fn.Name.Name)
		fullName := fmt.Sprintf("%s.%s.%s", pkg, recvType, fn.Name.Name)
		return fullName, nil
	}

	// package exported function
	fullName := fmt.Sprintf("%s.%s", pkg, fn.Name.Name)
	return fullName, nil
}

// FunctionContainsStmt 返回包含语句stmt的函数全名$xxx.$func
func FunctionContainsStmt(fset *token.FileSet, stmt ast.Stmt) (*ast.FuncDecl, error) {

	// 找到原文件名
	base := fset.File(stmt.Pos()).Base()
	idx1 := int(stmt.Pos()) - base
	idx2 := int(stmt.End()) - base

	fname := fset.Position(stmt.Pos()).Filename
	fset, file, err := ParseFile(fname)
	if err != nil {
		return nil, err
	}

	fnList, err := Functions(fset, file, func(fn *ast.FuncDecl) bool {
		if int(fn.Pos()) <= idx1 && int(fn.End()) >= idx2 {
			return true
		}
		return false
	})
	if err != nil {
		return nil, err
	}

	Assert(len(fnList) == 1, "should only 1 function existed")
	Assert(fnList[0] != nil, "function shouldn't be nil")

	return fnList[0], nil
}

// FunctionNameContainsStmt 返回包含语句stmt的函数全名$xxx.$func
func FunctionNameContainsStmt(fset *token.FileSet, stmt ast.Stmt) (string, error) {

	// 找到原文件名
	base := fset.File(stmt.Pos()).Base()
	idx1 := int(stmt.Pos()) - base
	idx2 := int(stmt.End()) - base

	fname := fset.Position(stmt.Pos()).Filename
	fset, file, err := ParseFile(fname)
	if err != nil {
		return "", err
	}

	fnList, err := Functions(fset, file, func(fn *ast.FuncDecl) bool {
		if int(fn.Pos()) <= idx1 && int(fn.End()) >= idx2 {
			return true
		}
		return false
	})
	if err != nil {
		return "", err
	}

	Assert(len(fnList) == 1, "should only 1 function existed")
	Assert(fnList[0] != nil, "function shouldn't be nil")

	return FunctionFullName(fset, fnList[0])
}

// PackageNameContainsStmt 返回包含语句stmt的package名
func PackageNameContainsStmt(fset *token.FileSet, stmt ast.Stmt) (string, error) {
	fname := fset.Position(stmt.Pos()).Filename
	fset, file, err := ParseFile(fname)
	if err != nil {
		return "", err
	}
	return file.Name.Name, nil
}

func ReceiverTypeOrPackageName(fset *token.FileSet, selectorExpr *ast.SelectorExpr) (typ string, err error) {

	pos := fset.Position(selectorExpr.Pos()).String()
	hint, _ := PosToString(fset, selectorExpr.Pos(), selectorExpr.End())

	ident, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		// TODO 这里有些类型忽略掉了，比如表达式类型*go/ast.Expr (like httpReq.Header...)
		return "", errors.New("not *ast.Ident")
	}

	// providing typ is package exported function
	typ = selectorExpr.X.(*ast.Ident).Name

	// it's a package exported function
	if ident.Obj == nil {
		return typ, nil
	}

	// it's a method name
	decl := ident.Obj.Decl
	if decl == nil {
		return "", errors.New("nil *obj.Decl")
	}

	// BUG: panic: interface conversion: interface {} is *ast.ValueSpec, not *ast.AssignStmt
	switch v := decl.(type) {
	case *ast.ValueSpec: // `var db *sql.DB`
		typ, _ = PosToString(fset, v.Pos(), v.End())
	case *ast.AssignStmt: // `a := 0` or `s := trpc.NewServer()`
		rhs := v.Rhs
		switch expr := rhs[0].(type) {
		case *ast.UnaryExpr: // `a := 0`
			typ = expr.X.(*ast.CompositeLit).Type.(*ast.Ident).Name
			if op := rhs[0].(*ast.UnaryExpr).Op.String(); len(op) != 0 {
				if op == "&" {
					op = "*"
				}
				typ = op + typ
			}
		case *ast.CallExpr: // `s := grpc.NewServer()`
			// TODO 这里需要从依赖的gomodules找到trpc.NewServer的定义，然后查出其返回值类型，
			// 如：github.com/golang/grpc/server.*Server
			//
			// 这里临时为了方便，先不做这部分外部依赖的检查，直接用语句字符串表示其类型了
			typ, _ = PosToString(fset, expr.Pos(), expr.End())
		case *ast.CompositeLit: // `wg := sync.WaitGroup{}`
			// TODO 这里需要从依赖的gomodules找到sync.WaitGroup的定义，然后查出其返回值类型，
			// 如：sync.WaitGroup
			typ, _ = PosToString(fset, expr.Pos(), expr.End())
		default:
			return "", fmt.Errorf("rhs[0] should be *ast.UnaryExpr or *ast.CallExpr, pos: %s, hint: %s", pos, hint)
		}
	case *ast.DeclStmt: // `var a = 1` or `var a int` or `var a int = 1`
		decl, ok := v.Decl.(*ast.GenDecl)
		if !ok {
			panic("not *ast.GenDecl")
		}
	NextSpec:
		for _, spec := range decl.Specs {
			valspec, ok := spec.(*ast.ValueSpec)
			if !ok {
				panic("not *ast.ValueSpec")
			}
			for _, name := range valspec.Names {
				if name.Name == ident.Obj.Name {
					switch star := valspec.Type.(type) {
					case *ast.Ident:
						typ = valspec.Type.(*ast.Ident).Name
					case *ast.StarExpr:
						typ, err = PosToString(fset, star.Pos(), star.End())
						if err != nil {
							panic(err)
						}
					default:
						panic("not supported *ast.DeclStmt")
					}
					break NextSpec
				}
			}
		}

	default:
		panic(fmt.Sprintf("what's the type, pos: %s, hint: %s, decl: %T", pos, hint, decl))
	}
	return typ, nil
}
