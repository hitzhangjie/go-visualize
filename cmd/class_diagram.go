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
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/hitzhangjie/go-visualize/plantuml"
	"github.com/kazukousen/gouml"
	"github.com/spf13/cobra"
)

// classDiagramCmd represents the class command
var classDiagramCmd = &cobra.Command{
	Use:   "class [directory]",
	Short: "输出类图",
	Long:  `输出类图.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// analyze directory
		dirs, _ := cmd.Flags().GetStringArray("dir")
		ignore, _ := cmd.Flags().GetStringArray("ignore")
		verbose, _ := cmd.Flags().GetBool("verbose")
		puml, _ := cmd.Flags().GetString("puml")

		buf := &bytes.Buffer{}
		buf.WriteString("@startuml\n")

		logger := log.NewNopLogger()
		if err := generate(logger, buf, ignore, dirs, verbose); err != nil {
			return err
		}
		buf.WriteString("@enduml\n")

		if err := writeFile(puml, buf.Bytes()); err != nil {
			return err
		}
		fmt.Printf("generate file: %s\n", puml)

		if err := plantuml.RenderPlantUML(puml); err != nil {
			return err
		}
		png := strings.TrimSuffix(puml, filepath.Ext(puml)) + ".png"
		fmt.Printf("generate file: %s\n", png)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(classDiagramCmd)

	classDiagramCmd.Flags().StringArray("dir", []string{"."}, "分析指定的目录列表")
	classDiagramCmd.Flags().StringArray("file", nil, "分析指定的文件列表")
	classDiagramCmd.Flags().String("puml", "file.class.puml", "输出到指定文件")
	classDiagramCmd.Flags().StringArray("ignore", nil, "忽略指定的文件列表")
	classDiagramCmd.Flags().BoolP("verbose", "v", false, "显示详细信息")
}

func generate(logger log.Logger, buf *bytes.Buffer, ignores []string, targets []string, verbose bool) error {
	gen := gouml.NewGenerator(logger, gouml.PlantUMLParser(logger), verbose)
	if len(ignores) > 0 {
		if err := gen.UpdateIgnore(ignores); err != nil {
			return err
		}
	}
	if len(targets) == 0 {
		targets = []string{"./"}
	}
	if err := gen.Read(targets); err != nil {
		return err
	}

	gen.WriteTo(buf)
	return nil
}

func writeFile(output string, buf []byte) error {
	if len(output) != 0 {
		file, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		file.Write(buf)
		defer file.Close()

		return nil
	}

	io.Copy(os.Stdout, bytes.NewBuffer(buf))
	return nil
}
