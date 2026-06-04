package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path"
	"strings"
)

type CmdType string

const (
	ExitCmd CmdType = "exit"
	TypeCmd CmdType = "type"

	EchoCmd CmdType = "echo"

	UnknownCmd CmdType = "unknown"
)

type Cmd struct {
	Type CmdType
	Args string
}

func ParseCmd(s string) Cmd {
	s = strings.Trim(s, " ")
	i := strings.Index(s, " ")

	name := s
	args := ""
	if i > -1 {
		name = s[:i]
		args = s[i+1:]
	}

	var t CmdType
	switch name {
	case string(ExitCmd):
		t = ExitCmd
	case string(EchoCmd):
		t = EchoCmd
	case string(TypeCmd):
		t = TypeCmd
	default:
		t = UnknownCmd
		args = name
	}

	return Cmd{t, args}
}

type Executor struct {
	binDirs []string
	out     io.Writer
}

func NewExecutor(w io.Writer, binDirs []string) *Executor {
	return &Executor{
		out:     w,
		binDirs: binDirs,
	}
}

func (e *Executor) Exec(cmd Cmd) error {
	switch cmd.Type {
	case ExitCmd:
		return nil
	case EchoCmd:
		return e.echo(cmd.Args)
	case TypeCmd:
		return e.exec_type(cmd.Args)
	case UnknownCmd:
		return e.printf("%s: command not found", cmd.Args)
	}

	return errors.New("unknown type of cmd")
}

func (e *Executor) printf(format string, args ...any) error {
	sf := fmt.Sprintf(format, args...)
	buf := bytes.NewBufferString(sf)
	_, err := buf.WriteTo(e.out)
	return err
}

func (e *Executor) echo(args string) error {
	return e.printf("%s", args)
}

func (e *Executor) exec_type(args string) error {
	cmd := ParseCmd(args)

	// check builtins
	if cmd.Type != UnknownCmd {
		return e.printf("%s is a shell builtin", cmd.Type)
	}

	p := cmd.Args

	// check executables
	for _, d := range e.binDirs {
		path, err := exec.LookPath(path.Join(d, p))
		if err == nil {
			return e.printf("%s is %s", p, path)
		}
	}

	// not found
	return e.printf("%s: not found", p)
}
