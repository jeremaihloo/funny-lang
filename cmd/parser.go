/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jeremaihloo/funny/langs"
	"github.com/spf13/cobra"
)

// parserCmd represents the parser command
var parserCmd = &cobra.Command{
	Use:   "parser",
	Short: "Parser dumps json parse a funny script file or funny script text into AST.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			filename := args[0]
			if _, err := os.Stat(filename); err != nil {
				fmt.Printf("file not found %s\n", filename)
				return
			}
			data, _ := ioutil.ReadFile(filename)
			parser := langs.NewParser(data)
			parser.Consume("")
			var items langs.Block
			for {
				item := parser.ReadStatement()
				if item == nil {
					break
				}
				items = append(items, item)
				// fmt.Printf("%s %s\n", langs.Typing(item), item.String())
			}
			data, err := json.MarshalIndent(items, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(data))
		}
	},
}

func init() {
	rootCmd.AddCommand(parserCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// parserCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// parserCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
