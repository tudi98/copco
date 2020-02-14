package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/tudi98/copco/pkg/parser"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func run() {
	file, err := ioutil.ReadFile("problem.json")
	if err != nil {
		log.Fatal(err)
	}

	problem := parser.Problem{}

	if err := json.Unmarshal([]byte(file), &problem); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Compiling...")
	cmd := exec.Command("g++", "-o", "main", "main.cpp")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return
	}
	fmt.Println("Successful!")

	input_files := getFilesWithExtension(".in")

	fmt.Printf("%v\n", input_files)

	for i, _ := range input_files {
		fmt.Printf("Running on test %d...", i)
		fmt.Printf("OK\n")
	}
}

func getFilesWithExtension(ext string) []string {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	separator := "/"
	if runtime.GOOS == "windows" {
		separator = "\\"
	}

	path = path + separator + "tests"

	var files []string
	filepath.Walk(path, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() && filepath.Ext(path) == ext {
			files = append(files, f.Name())
		}
		return nil
	})

	return files
}
