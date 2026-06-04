package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
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

		err = cmd.Exec()
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

	cmd := NewCmd(os.Stdout, s[:len(s)-1])
	return cmd, nil

}
