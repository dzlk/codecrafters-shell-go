package main

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type Executor struct {
	out io.Writer
	err io.Writer
}

func NewExecutor(w io.Writer) *Executor {
	return &Executor{
		out: w,
		err: w,
	}
}

func (e *Executor) Exec(cmd Cmd) error {
	switch {
	case cmd.IsExit():
		return nil
	case cmd.Name == EchoCmd:
		return e.echo(cmd.Args)
	case cmd.Name == TypeCmd:
		return e.exec_type(cmd)
	case cmd.IsExternal():
		return e.exec(cmd)
	}
	return e.printf("%s: command not found\n", cmd.Args)
}

func (e *Executor) printf(format string, args ...any) error {
	sf := fmt.Sprintf(format, args...)
	buf := bytes.NewBufferString(sf)
	_, err := buf.WriteTo(e.out)
	return err
}

func (e *Executor) echo(args string) error {
	return e.printf("%s\n", args)
}

func (e *Executor) exec_type(typeCmd Cmd) error {
	cmd := typeCmd.SubCmd
	if cmd.IsExternal() {
		return e.printf("%s is %s\n", cmd.Name, cmd.Path)
	}

	if cmd.IsUnknown() {
		return e.printf("%s: not found\n", cmd.Args)
	}

	return e.printf("%s is a shell builtin\n", cmd.Name)
}

func (e *Executor) exec(cmd Cmd) error {
	args := strings.Fields(cmd.Args)
	c := exec.Command(cmd.Name, args...)
	c.Stdout = e.out
	c.Stderr = e.err

	return c.Run()
}
