package goast

import (
	"fmt"
	"go/ast"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestParseFile(t *testing.T) {
	fset, file, err := ParseFile("../../testdata/trpc/main.go")
	assert.Nil(t, err)

	t.Logf("fset: %v", fset)
	t.Logf("file: %v", file)
}

func TestFunctions(t *testing.T) {
	fset, file, err := ParseFile("../../testdata/trpc/main.go")
	assert.Nil(t, err)

	fnList, err := Functions(fset, file, nil)
	assert.Nil(t, err)
	spew.Printf("functions: %v\n", fnList)

	filter := func(node *ast.FuncDecl) bool {
		if strings.Contains(node.Name.Name, "main") {
			return true
		}
		return false
	}
	fnList, err = Functions(fset, file, filter)
	assert.Nil(t, err)
	spew.Printf("functions: %v\n", fnList)
}

func TestParseDir(t *testing.T) {
	fset, pkgs, err := ParseDir("../../testdata/trpc", false)
	assert.Nil(t, err)
	assert.NotNil(t, fset)
	assert.NotEmpty(t, pkgs)
	fmt.Printf("len(pkgs) == %d\n", len(pkgs))

	fset, pkgs, err = ParseDir("../../testdata/trpc", true)
	assert.Nil(t, err)
	assert.NotNil(t, fset)
	assert.NotEmpty(t, pkgs)
	fmt.Printf("len(pkgs) == %d\n", len(pkgs))

	for k, v := range pkgs {
		spew.Printf("%s => %v\n", k, v)
	}
}

func TestMethodReceiverType(t *testing.T) {
	fset, file, err := ParseFile("../../testdata/trpc/helloworld.go")
	assert.Nil(t, err)

	filter := func(node *ast.FuncDecl) bool {
		if strings.Contains(node.Name.Name, "Hello") {
			return true
		}
		return false
	}
	fnList, err := Functions(fset, file, filter)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(fnList))

	typ, err := MethodReceiverTypeName(fnList[0])
	assert.Nil(t, err)
	assert.Equal(t, "*helloworldServiceImpl", typ)
}

func TestPrint(t *testing.T) {
	err := Print("./goast.go")
	assert.Nil(t, err)
}

func TestPosToString(t *testing.T) {
	fset, file, err := ParseFile("../../testdata/trpc/helloworld.go")
	assert.Nil(t, err)

	filter := func(node *ast.FuncDecl) bool {
		if strings.Contains(node.Name.Name, "Hello") {
			return true
		}
		return false
	}
	fnList, err := Functions(fset, file, filter)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(fnList))

	fn := fnList[0]
	s, err := PosToString(fset, fn.Name.Pos(), fn.Name.End())
	assert.Nil(t, err)
	assert.NotEmpty(t, s)
	t.Logf("pos:%v, string:%s", fn.Pos(), s)
}

func TestFunctionFullName(t *testing.T) {
	fset, file, err := ParseFile("../../testdata/ast/func.go")
	assert.Nil(t, err)

	fnList, err := Functions(fset, file, nil)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(fnList))

	for idx, fn := range fnList {
		fullName, err := FunctionFullName(fset, fn)
		assert.Nil(t, err)
		switch idx {
		case 0:
			assert.Equal(t, "ast.(*Cat).Run", fullName)
		case 1:
			assert.Equal(t, "ast.Hello", fullName)
		}
		t.Logf("function full name: %s", fullName)
	}
}

func TestFunctionContainsStmt(t *testing.T) {
	fset, file, err := ParseFile("../../testdata/ast/func.go")
	assert.Nil(t, err)

	fnList, err := Functions(fset, file, nil)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(fnList))

	expected := map[int]string{
		0: "ast.(*Cat).Run",
		1: "ast.Hello",
	}
	for idx, fn := range fnList {
		for _, stmt := range fn.Body.List {
			fnName, err := FunctionNameContainsStmt(fset, stmt)
			assert.Nil(t, err)
			assert.Equal(t, expected[idx], fnName)
		}
	}
}

func TestPackageContainsStmt(t *testing.T) {
	fset, file, err := ParseFile("../../testdata/ast/func.go")
	assert.Nil(t, err)

	fnList, err := Functions(fset, file, nil)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(fnList))

	for _, fn := range fnList {
		for _, stmt := range fn.Body.List {
			pkgName, err := PackageNameContainsStmt(fset, stmt)
			assert.Nil(t, err)
			assert.Equal(t, "ast", pkgName)
		}
	}
}
