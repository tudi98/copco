package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

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
}
