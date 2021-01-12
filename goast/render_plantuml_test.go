package goast

import (
	"bytes"
	"fmt"
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderIfStmtWithPlantuml(t *testing.T) {
	file := "../../testdata/ast/ifelse.go"

	fset, f, err := ParseFile(file)
	assert.Nil(t, err)

	fnList, err := Functions(fset, f, nil)
	assert.Nil(t, err)
	assert.Equal(t, len(fnList), 1)

	fn := fnList[0]
	for _, stmt := range fn.Body.List {
		if ifstmt, ok := stmt.(*ast.IfStmt); ok && ifstmt != nil {
			buf := bytes.Buffer{}
			err := RenderIfStmt(fset, ifstmt, 0, &buf)
			assert.Nil(t, err)
			fmt.Printf("\n%s\n", string(buf.Bytes()))
		}
	}
}

func TestRenderSwitchStmtWithPlantUML(t *testing.T) {
	file := "../../testdata/ast/switch.go"

	fset, f, err := ParseFile(file)
	assert.Nil(t, err)

	fnList, err := Functions(fset, f, nil)
	assert.Nil(t, err)
	assert.Equal(t, len(fnList), 1)

	fn := fnList[0]
	for _, stmt := range fn.Body.List {
		if switchStmt, ok := stmt.(*ast.SwitchStmt); ok && switchStmt != nil {
			buf := bytes.Buffer{}
			err := RenderSwitchStmt(fset, switchStmt, &buf)
			assert.Nil(t, err)
			fmt.Printf("\n%s\n", string(buf.Bytes()))
		}
	}
}

func TestRenderForStmtWithPlantUML(t *testing.T) {
	file := "../../testdata/ast/for.go"

	fset, f, err := ParseFile(file)
	assert.Nil(t, err)

	fnList, err := Functions(fset, f, nil)
	assert.Nil(t, err)
	assert.Equal(t, len(fnList), 1)

	fn := fnList[0]
	for _, stmt := range fn.Body.List {
		if forStmt, ok := stmt.(*ast.ForStmt); ok && forStmt != nil {
			buf := bytes.Buffer{}
			err := RenderForStmt(fset, forStmt, &buf)
			assert.Nil(t, err)
			fmt.Printf("\n%s\n", buf.String())
		}
	}
}

func TestRenderWithPlantUML(t *testing.T) {

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
			buf, err := RenderStmt(fset, stmt)
			if err != nil {
				assert.Equal(t, ErrIgnoreStmt, err)
				continue
			}
			t.Logf("render in puml:\n%s\n", string(buf))
		}
		fmt.Println()
	}
}
