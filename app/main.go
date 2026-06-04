package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	pathEnv, _ := os.LookupEnv("PATH")
	p := NewCmdParser(pathEnv)
	e := NewExecutor(os.Stdout)

	for {
		fmt.Print("$ ")

		cmd, err := readCmd(p)
		if err != nil {
			handleError(fmt.Errorf("parse cmd failed: %w", err))
			break
		}

		if cmd.IsExit() {
			break
		}

		err = e.Exec(cmd)
		if err != nil {
			handleError(fmt.Errorf("exec cmd failed: %w", err))
			break
		}
	}
}

func handleError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func readCmd(p *CmdParser) (Cmd, error) {
	s, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return Cmd{}, err
	}

	cmd := p.Parse(s[:len(s)-1])
	return cmd, nil
}
