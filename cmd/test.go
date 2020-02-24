package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/shirou/gopsutil/process"
	"github.com/spf13/cobra"
	"github.com/tudi98/copco/pkg/models"
	"github.com/tudi98/copco/pkg/utils"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run source code on all test cases",
	Run: func(cmd *cobra.Command, args []string) {
		test()
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}

func test() {
	file, err := ioutil.ReadFile("problem.json")
	if err != nil {
		log.Fatal("Error while opening problem.json")
	}
	problem := models.Problem{}
	if err := json.Unmarshal(file, &problem); err != nil {
		log.Fatal("Error while parsing problem.json")
	}

	compileCmd := exec.Command("g++", "-O2", "-o", "main", "main.cpp")
	compileCmd.Stderr = os.Stderr
	compileCmd.Stdout = os.Stdout
	if err := compileCmd.Run(); err != nil {
		color.Red("Compilation error!")
		return
	}
	color.Green("Compiled successfully!")

	inputFiles := utils.GetFilesWithExtension(".in")

	separator := "/"
	if runtime.GOOS == "windows" {
		separator = "\\"
	}

	fmt.Println("Running on tests: ")
	for _, filePath := range inputFiles {
		testName := strings.Split(filePath, separator)[2]
		runCmd := exec.Command("./main")

		fmt.Printf(" - %s : ", testName)

		stdin, _ := runCmd.StdinPipe()
		inFile, err := os.Open(filePath)
		if err != nil {
			log.Fatalf("Error while opening %s", filePath)
		}

		_, err = io.Copy(stdin, inFile)
		if err != nil {
			log.Fatal(err)
		}

		_ = inFile.Close()
		_ = stdin.Close()

		outFilePath := filePath[:len(filePath)-2] + "out"
		outFile, err := os.OpenFile(outFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatalf("Error when opening/creating %s", outFilePath)
		}
		runCmd.Stdout = outFile

		if err := runCmd.Start(); err != nil {
			color.Red("Runtime Error -> %s", err.Error())
			continue
		}

		pid := int32(runCmd.Process.Pid)

		ch := make(chan error)
		go func() {
			ch <- runCmd.Wait()
		}()

		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(time.Duration(problem.TimeLimit) * time.Millisecond)
			timeout <- true
		}()

		finished := false
		errorFree := true

		for !finished {
			select {
			case err := <-ch:
				if err != nil {
					color.Red("Runtime Error -> %s", err.Error())
					errorFree = false
				}
				finished = true
			case _ = <-timeout:
				_ = runCmd.Process.Signal(syscall.SIGINT)
				color.Red("Time Limit Exceeded")
				finished = true
				errorFree = false
			default:
				newProcess, err := process.NewProcess(pid)
				if err == nil {
					memoryInfo, err := newProcess.MemoryInfo()
					if err == nil && memoryInfo.RSS > uint64(problem.MemoryLimit) {
						_ = runCmd.Process.Signal(syscall.SIGINT)
						color.Red("Memory Limit Exceeded")
						finished = true
						errorFree = false
					}
				}
			}
		}

		if !errorFree {
			continue
		}

		_ = outFile.Close()

		okFilePath := filePath[:len(filePath)-2] + "ok"

		okFileBytes, err := ioutil.ReadFile(okFilePath)
		if err != nil {
			log.Fatalf("Error while opening %s", okFilePath)
		}

		outFileBytes, err := ioutil.ReadFile(outFilePath)
		if err != nil {
			log.Fatalf("Error while opening %s", outFilePath)
		}

		ok := utils.TrimLines(okFileBytes)
		out := utils.TrimLines(outFileBytes)

		if bytes.Compare(ok, out) == 0 {
			color.Green("Passed")
		} else {
			color.Red("Wrong Answer")
		}
	}
}
