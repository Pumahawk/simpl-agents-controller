package cmd

import "fmt"

type Command interface {
	Name() string
	Run(args []string) error
}

type Cmd struct {
	CName string
	CRun  func(args []string) error
}

func (c *Cmd) Name() string {
	return c.CName
}
func (c *Cmd) Run(args []string) error {
	return c.CRun(args)
}

type Group struct {
	GName string
	Cmds  []Command
}

func (g *Group) Name() string {
	return g.GName
}

func (g *Group) Run(args []string) error {
	if len(args) > 0 {
		name := args[0]
		args := args[1:]
		for i := range g.Cmds {
			cmd := g.Cmds[i]
			if cmd.Name() == name {
				return cmd.Run(args)
			}
		}
		return fmt.Errorf("command %q not found", name)
	}
	return fmt.Errorf("missing command name")
}
