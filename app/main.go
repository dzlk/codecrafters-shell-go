package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	e := NewExecutor(os.Stdout, getExecutablePaths())

	for {
		fmt.Print("$ ")

		cmd, err := readCmd()
		if err != nil {
			handleError(fmt.Errorf("parse cmd failed: %w", err))
			break
		}

		if cmd.Type == ExitCmd {
			break
		}

		err = e.Exec(cmd)
		if err != nil {
			handleError(fmt.Errorf("exec cmd failed: %w", err))
			break
		}

		fmt.Println()
	}
}

func handleError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func readCmd() (Cmd, error) {
	s, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return Cmd{}, err
	}

	cmd := ParseCmd(s[:len(s)-1])
	return cmd, nil
}

func getExecutablePaths() []string {
	pathEnv, ok := os.LookupEnv("PATH")
	binDirs := []string{}
	if ok {
		dirs := filepath.SplitList(pathEnv)

		for _, d := range dirs {
			info, err := os.Stat(d)
			if !os.IsNotExist(err) && info.IsDir() {
				binDirs = append(binDirs, d)
			}
		}
	}

	return binDirs
}
