package main

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const (
	ExitCmd string = "exit"
	TypeCmd string = "type"
	EchoCmd string = "echo"

	UnknownCmd string = "unknown"
)

type Cmd struct {
	Name string
	Args string

	Path   string
	SubCmd *Cmd
}

func (cmd *Cmd) IsExit() bool {
	return cmd.Name == ExitCmd
}

func (cmd *Cmd) IsExternal() bool {
	return cmd.Path != ""
}

func (cmd *Cmd) IsUnknown() bool {
	return cmd.Name == UnknownCmd
}

type CmdParser struct {
	binDirs []string
}

func NewCmdParser(pathList string) *CmdParser {
	binDirs := []string{}
	dirs := filepath.SplitList(pathList)

	for _, d := range dirs {
		info, err := os.Stat(d)
		if !os.IsNotExist(err) && info.IsDir() {
			binDirs = append(binDirs, d)
		}
	}

	return &CmdParser{binDirs: binDirs}
}

func (p *CmdParser) Parse(s string) Cmd {
	s = strings.Trim(s, " ")
	name, args, _ := strings.Cut(s, " ")

	switch name {
	case ExitCmd:
		return Cmd{Name: ExitCmd}
	case EchoCmd:
		return Cmd{Name: EchoCmd, Args: args}
	case TypeCmd:
		cmd := p.Parse(args)
		return Cmd{Name: TypeCmd, Args: args, SubCmd: &cmd}
	}

	// check if cmd is external program
	cmdPath := ""
	for _, d := range p.binDirs {
		path, err := exec.LookPath(path.Join(d, name))
		if err == nil {
			cmdPath = path
			break
		}
	}

	if cmdPath != "" {
		return Cmd{Name: name, Args: args, Path: cmdPath}
	}

	return Cmd{Name: UnknownCmd, Args: name}
}
