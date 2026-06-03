package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	for {
		fmt.Print("$ ")

		cmd, err := readCmd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading command:", err)
			os.Exit(1)
		}

		if cmd == "exit" {
			break
		}

		if strings.HasPrefix(cmd, "echo") {
			fmt.Println(cmd[5:])
			continue
		}

		fmt.Printf("%s: command not found", cmd)
		fmt.Println()
	}
}

func readCmd() (string, error) {
	command, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}

	return command[:len(command)-1], nil
}
