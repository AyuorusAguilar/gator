package commands

import (
	"fmt"
	"github.com/AyuorusAguilar/gator/internal/config"
	"github.com/AyuorusAguilar/gator/internal/state"
)
type Command struct {
	Name string
	Arguments []string
}
type Commander struct {
	CommandMap map[string]func(*state.State, Command) error
}
func (c *Commander) Run(s *state.State, cmd Command) error {	
	if call, exists := c.CommandMap[cmd.Name]; exists {
		err := call(s, cmd)
		return err
	}
	return fmt.Errorf("Command '%s' does not exist\n", cmd.Name)
}
func (c *Commander) Register(name string, f func(*state.State, Command) error) error {
	if _, exists := c.CommandMap[name]; exists {
		return fmt.Errorf("Handler '%s' already exists\n", name)
	}
	c.CommandMap[name] = f
	return nil
}
func NewCommander() Commander {
	return Commander{CommandMap: map[string]func(*state.State, Command) error {}}
}

func HandlerLogin(s *state.State, cmd Command) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("Error:\n\tNo username provided\n")
	}
	err := config.SetUser(s.Cfg, cmd.Arguments[0])
	if err != nil {
		return err
	}
	fmt.Printf("User has been set to %s\n", s.Cfg.CurrentUserName)
	return nil
 }