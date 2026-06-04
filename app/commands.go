package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
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
	out io.Writer

	Type CmdType
	Args string
}

func NewCmd(out io.Writer, s string) Cmd {
	t, args := parseCmd(s)

	return Cmd{
		out:  out,
		Type: t,
		Args: args,
	}
}

func parseCmd(s string) (CmdType, string) {
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

	return t, args
}

func (cmd Cmd) Exec() error {
	switch cmd.Type {
	case ExitCmd:
		return nil
	case EchoCmd:
		return cmd.echo()
	case TypeCmd:
		return cmd.exec_type()
	case UnknownCmd:
		return cmd.printf("%s: command not found", cmd.Args)
	}

	return errors.New("unknown type of cmd")
}

func (cmd Cmd) printf(format string, args ...any) error {
	sf := fmt.Sprintf(format, args...)
	buf := bytes.NewBufferString(sf)
	_, err := buf.WriteTo(cmd.out)
	return err
}

func (cmd Cmd) echo() error {
	return cmd.printf("%s", cmd.Args)
}

func (cmd Cmd) exec_type() error {
	sc, args := parseCmd(cmd.Args)

	if sc == UnknownCmd {
		return cmd.printf("%s: not found", args)
	}

	return cmd.printf("%s: is a shell builtin", sc)
}
