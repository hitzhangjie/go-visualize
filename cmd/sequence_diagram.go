/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package code

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"git.code.oa.com/gitcode/gitcode/cmd/annotations"
	"git.code.oa.com/gitcode/gitcode/http"
	"git.code.oa.com/gitcode/gitcode/util/goast"
	"github.com/spf13/cobra"
)

func init() {
	visualizeCmd.AddCommand(visualizeSequenceCmd)

	visualizeSequenceCmd.Flags().String("func", "main.main", "specify the function to analyze")

	visualizeSequenceCmd.Flags().StringP("render", "r", "plantuml", "specify the render mode, console or plantuml")
	visualizeSequenceCmd.Flags().String("puml", "file.sequence.puml", "输出到指定文件")
}

// visualizeSequenceCmd represents the update command
var visualizeSequenceCmd = &cobra.Command{
	Use:   "sequence [directory]",
	Short: "visualize code in sequence diagram",
	Long:  `visualize code in sequence diagram`,
	Annotations: map[string]string{
		annotations.CmdGroupAnnotation: annotations.CmdGroupCode,
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		function, _ := cmd.Flags().GetString("func")
		if len(function) == 0 {
			return errors.New("invalid --func")
		}
		puml, _ := cmd.Flags().GetString("puml")
		//render, _ := cmd.Flags().GetString("render")

		// analyze this directory
		dir := "."
		if len(args) != 0 {
			dir = args[0]
		}

		var (
			fset *token.FileSet
			pkgs map[string]*ast.Package
			err  error
		)

		// analyze this function
		fmt.Println("analyze function:", function)
		fset, pkgs, err = goast.ParseDir(dir, true)
		if err != nil {
			return err
		}

		funcDecl, err := goast.FindFunction(pkgs, function)
		if err != nil {
			return err
		}

		buf, err := goast.RenderFunctionWithPlantUML(funcDecl, fset, pkgs)
		if err != nil {
			return err
		}

		dat := bytes.Buffer{}
		dat.WriteString("@startuml\n")
		dat.Write(buf.Bytes())
		dat.WriteString("@enduml\n")

		if err := writeFile(puml, dat.Bytes()); err != nil {
			return err
		}
		fmt.Printf("generate file: %s\n", puml)

		if err := renderPlantUML(puml); err != nil {
			return err
		}
		png := strings.TrimSuffix(puml, filepath.Ext(puml)) + ".png"
		fmt.Printf("generate file: %s\n", png)

		return nil
	},
}

func renderPlantUML(pumlFile string) error {

	// download plantuml.jar if not found
	home, _ := os.UserHomeDir()
	jar := filepath.Join(home, "plantuml.jar")
	_, err := os.Lstat(jar)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		err := http.DownloadFile(jar, "https://nchc.dl.sourceforge.net/project/plantuml/plantuml.jar")
		if err != nil {
			return fmt.Errorf("download plantuml.jar error: %v", err)
		}
	}

	cmd := exec.Command("java", "-jar", jar, "-progress", pumlFile)
	if msg, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("run plantuml error: %v,\nmsg:%s", err, msg)
	}
	return nil
}
