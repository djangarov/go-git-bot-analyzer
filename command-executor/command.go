package command

import (
	"bytes"
	"os/exec"
	"strings"
)

const WHITESPACE = " "

type Argument struct {
	Argument string
	Param    string
}

type Command struct {
	App      string
	Argument []Argument
}

const ShellToUse = "bash"

func (c *Command) parse() string {
	var command strings.Builder
	command.WriteString(c.App)
	command.WriteString(WHITESPACE)

	for _, args := range c.Argument {
		command.WriteString(args.Argument)
		command.WriteString(WHITESPACE)
		command.WriteString(args.Param)
		command.WriteString(WHITESPACE)
	}
	return command.String()
}

func Execute(command Command) (error, string, string) {

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command.parse())
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}
