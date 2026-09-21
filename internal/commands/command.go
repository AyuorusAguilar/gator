package commands

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/AyuorusAguilar/gator/internal/config"
	"github.com/AyuorusAguilar/gator/internal/database"
	"github.com/AyuorusAguilar/gator/internal/rss"
	"github.com/AyuorusAguilar/gator/internal/state"
	"github.com/google/uuid"
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

	name := cmd.Arguments[0]
	if _, err := s.Db.GetUser(context.Background(), name); err != nil {
		return fmt.Errorf("Error:\n\tUser with name %s doesn't exist, register it first!\n", name)
	}

	err := config.SetUser(s.Cfg, name)
	if err != nil {
		return err
	}

	fmt.Printf("User has been set to %s\n", s.Cfg.CurrentUserName)
	return nil
 }
func HandlerRegister(s *state.State, cmd Command) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("Error:\n\tNo username provided\n")
	}
	name := cmd.Arguments[0]
	if _, err := s.Db.GetUser(context.Background(), name); err == nil {
		return fmt.Errorf("Error:\n\tName already in use, pick a diferent one\n")
	}
	if _, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: sql.NullTime{Time: time.Now()},
		UpdatedAt: sql.NullTime{Time: time.Now()},
		Name: name,
	}); err != nil {
		return fmt.Errorf("Error:\n\tAn error ocurred. Our team is currently working very hard investigating the causes\n") 
	}

	err := config.SetUser(s.Cfg, cmd.Arguments[0])
	if err != nil {
		return err
	}

	fmt.Printf("Registred user %s\n", name)

	return nil
 }
func HandlerReset(s *state.State, cmd Command) error {
	if err := s.Db.DropIsmu(context.Background()); err != nil {
		return fmt.Errorf("Error:\n\t??? A truly unexpected event...\n")
	}

	return nil
 }

 func HandlerUsers(s *state.State, cmd Command) error {
	if data, err := s.Db.GetUsers(context.Background()); err != nil {
		return fmt.Errorf("Error:\n\t??? A truly unexpected event...\n")
	} else {
		fmt.Printf("Listing %d users:\n", len(data))
		for _, user := range data {
			str := "\t*  %s\n"
			if user == s.Cfg.CurrentUserName {
				str ="\t*  %s (current)\n"
			}
			fmt.Printf(str, user)
		}
		return nil
	}
	
 }
 
 func HandlerAggregator(s *state.State, cmd Command) error {
	content, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("Error:\n\t??? A truly unexpected event...\n%v\n", err)
	}
	fmt.Printf("Content:\n%v", content)
	return nil
 }
