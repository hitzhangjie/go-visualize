package goast

import (
	"fmt"
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderIfStmtWithConsole(t *testing.T) {
	file := "../../testdata/ast/ifelse.go"

	fset, f, err := ParseFile(file)
	assert.Nil(t, err)

	fnList, err := Functions(fset, f, nil)
	assert.Nil(t, err)
	assert.Equal(t, len(fnList), 1)

	fn := fnList[0]
	for _, stmt := range fn.Body.List {
		if ifstmt, ok := stmt.(*ast.IfStmt); ok && ifstmt != nil {
			err := RenderIfStmtWithConsole(fset, ifstmt, 0)
			assert.Nil(t, err)
		}
	}
}

func TestRenderSwitchStmtWithConsole(t *testing.T) {
	file := "../../testdata/ast/switch.go"

	fset, f, err := ParseFile(file)
	assert.Nil(t, err)

	fnList, err := Functions(fset, f, nil)
	assert.Nil(t, err)
	assert.Equal(t, len(fnList), 1)

	fn := fnList[0]
	for _, stmt := range fn.Body.List {
		if switchStmt, ok := stmt.(*ast.SwitchStmt); ok && switchStmt != nil {
			err := RenderSwitchStmtWithConsole(fset, switchStmt)
			assert.Nil(t, err)
		}
	}
}

func TestRenderForStmtWithConsole(t *testing.T) {
	file := "../../testdata/ast/for.go"

	fset, f, err := ParseFile(file)
	assert.Nil(t, err)

	fnList, err := Functions(fset, f, nil)
	assert.Nil(t, err)
	assert.Equal(t, len(fnList), 1)

	fn := fnList[0]
	for _, stmt := range fn.Body.List {
		if forStmt, ok := stmt.(*ast.ForStmt); ok && forStmt != nil {
			err := RenderForStmtWithConsole(fset, forStmt)
			assert.Nil(t, err)
		}
	}
}

func TestRenderWithConsole(t *testing.T) {

	files := []string{
		"../../testdata/ast/ifelse.go",
		"../../testdata/ast/switch.go",
		"../../testdata/ast/for.go",
	}
	for _, file := range files {

		fset, f, err := ParseFile(file)
		assert.Nil(t, err)

		fnList, err := Functions(fset, f, nil)
		assert.Nil(t, err)
		assert.Equal(t, len(fnList), 1)

		fn := fnList[0]
		for _, stmt := range fn.Body.List {
			err := RenderStmtWithConsole(fset, stmt)
			if err != nil {
				assert.Equal(t, ErrIgnoreStmt, err)
				continue
			}
		}
		fmt.Println()
	}
}
