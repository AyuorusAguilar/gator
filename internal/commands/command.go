package commands

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/AyuorusAguilar/gator/internal/config"
	"github.com/AyuorusAguilar/gator/internal/database"
	"github.com/AyuorusAguilar/gator/internal/rss"
	"github.com/AyuorusAguilar/gator/internal/state"
	"github.com/google/uuid"
)

type Command struct {
	Name      string
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
	return Commander{CommandMap: map[string]func(*state.State, Command) error{}}
}

func HandlerLogin(s *state.State, cmd Command) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("Error:\n\tNo username provided\n")
	}

	name := cmd.Arguments[0]
	if res, err := s.Db.GetUser(context.Background(), name); err != nil {
		return fmt.Errorf("Error:\n\tUser with name %s doesn't exist, register it first!\n", name)
	} else {
		err := config.SetUser(s.Cfg, res.Name, res.ID)
		if err != nil {
			return err
		}
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
	if res, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}); err != nil {
		return fmt.Errorf("Error:\n\tAn error ocurred. Our team is currently working very hard investigating the causes\n")
	} else {
		err := config.SetUser(s.Cfg, res.Name, res.ID)
		if err != nil {
			return err
		}
		fmt.Printf("Registred user %s\n", name)
		return nil
	}
}
func HandlerReset(s *state.State, cmd Command) error {
	if err := s.Db.DropIsmu(context.Background()); err != nil {
		return fmt.Errorf("Error:\n\t??? A truly unexpected event...\n")
	}
	config.ResetConf(s.Cfg)
	return nil
}
func HandlerUsers(s *state.State, cmd Command) error {
	if data, err := s.Db.GetUsers(context.Background()); err != nil {
		return fmt.Errorf("Error:\n\t??? A truly unexpected event...\n")
	} else {
		fmt.Printf("Listing %d users:\n", len(data))
		for _, user := range data {
			str := "\t*  %s\n"
			if user.Name == s.Cfg.CurrentUserName {
				str = "\t*  %s (current)\n"
			}
			fmt.Printf(str, user.Name)
		}
		return nil
	}
}
func HandlerAggregator(s *state.State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("Error:\n\tNo interval argument provided. Please Write a valid time! Example 1h30m10s or 5s\n")
	}
	timetofetch := cmd.Arguments[0]
	ActualInterval, err := time.ParseDuration(timetofetch)
	if err != nil {
		return fmt.Errorf("Error:\n\tPlease Write a valid time! Example 1h30m10s or 5s\n")
	}

	fmt.Printf("\n Fetching channels each %s\n", timetofetch)
	tic := time.NewTicker(ActualInterval)
	for ; ; <-tic.C {
		toFetch, err := s.Db.GetOldestFetchedByUserId(context.Background(), s.Cfg.CurrentUserId)
		if err != nil {
			return fmt.Errorf("Error:\n\tAn error ocurred fetching registered feeds, please make sure there is at least one registered feed!s\n")
		}

		content, err := rss.FetchFeed(context.Background(), toFetch.Url)
		if err != nil {
			fmt.Printf("Error:\n\tAn error ocurred fetching feed. We'll retry later, but maybe you'd want to check the feed info using the 'feeds' command and make sure it is all right!.\n Error info:\n\t%v\n", err)
			continue
		}
		err = s.Db.MarkFetched(context.Background(), database.MarkFetchedParams{
			ID:            toFetch.ID,
			LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
		})
		if err != nil {
			return fmt.Errorf("Error:\n\tAn error ocurred updating the timestamp of the feed, closing the loop to avoid repeted requests\n", err)
		}

		
		for _, item := range content.Channel.Item {
			var deit time.Time
			var err error
			var pubdeit sql.NullTime
			deit, err = time.Parse("02 Jan 2006 15:04:05 -0700", item.PubDate)
			if err != nil {
				pubdeit = sql.NullTime{
					Valid: false,
				}
			} else {
				pubdeit = sql.NullTime{
					Time: deit,
					Valid: true,
				}
			}
			_, err = s.Db.CreatePost(context.Background(), database.CreatePostParams{
				ID: uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Title: item.Title,
				Description: item.Description,
				Url: item.Link,
				PublishedAt: pubdeit,
				FeedID: toFetch.ID,
			})
			if err != nil && !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				fmt.Printf("Failed to save post on database! How werid! %v", item)
			}
		}
		
		/* fmt.Printf("Got %d titles from %s:\n", len(content.Channel.Item), toFetch.Name)
		for _, item := range content.Channel.Item {
			fmt.Printf("\t%s\n", item.Title)
		} */
	}
}
func HandlerBrowse(s *state.State, cmd Command, user database.User) error {
	var err error
	limit := 2
	if len(cmd.Arguments) > 0 {
		limit, err = strconv.Atoi(cmd.Arguments[0])
		if err != nil {
			fmt.Printf("No valid limit provided, defaulting to 2\n")
		}
		limit = 2
		err = nil
	}

	posts, err := s.Db.GetPostsByUserId(context.Background(), database.GetPostsByUserIdParams{
		UserID: s.Cfg.CurrentUserId,
		Limit: int32(limit),
	})
	if err != nil {
		return fmt.Errorf("Error:\n\tAn error ocurred while getting the posts! Make sure you are subscribed to an added feed and you've run the 'agg' command at least once to fetch posts!\ns")
	}

	for _, post := range posts {
		fmt.Printf("%s\n\n%s\n---\n", post.Title, post.Description)
	}

	return nil
}
func HandlerAddFeed(s *state.State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("Error:\n\tNot enough arguments! Provide a name and a url for the new feed\n")
	}
	name := cmd.Arguments[0]
	url := cmd.Arguments[1]

	var err error
	feed_id := uuid.New()
	_, err = s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        feed_id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    s.Cfg.CurrentUserId,
	})
	if err != nil {
		return fmt.Errorf("Error:\n\tAn error ocurred. Our team is currently working very hard investigating the causes\n")
	}

	if res, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    s.Cfg.CurrentUserId,
		FeedID:    feed_id,
	}); err != nil {
		return fmt.Errorf("Error:\n\tAn error ocurred registering the follow: %v\n", err)
	} else {
		fmt.Printf("%s added and started to follow \"%s\" at \"%s\"!\n", res.UserName, res.FeedName, res.FeedUrl)
		return nil
	}
}
func GetUserID(s *state.State, name string) (uuid.UUID, error) {
	id, err := s.Db.GetId(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}
func HandlerGetFeeds(s *state.State, cmd Command) error {
	var data []database.Feed
	var users []database.GetUsersRow
	usersDict := make(map[uuid.UUID]string)
	var err error

	if data, err = s.Db.GetFeeds(context.Background()); err != nil {
		return fmt.Errorf("Error:\n\t??? A truly unexpected event...\n")
	}
	if len(data) < 1 {
		fmt.Printf("\tCurrently, there are no feeds registered! Please register a feed using 'addfeed' command")
		return nil
	}

	if users, err = s.Db.GetUsers(context.Background()); err != nil {
		return fmt.Errorf("Error:\n\t??? A truly unexpected event...\n")
	}

	for _, user := range users {
		usersDict[user.ID] = user.Name
	}

	fmt.Printf("Listing %d feeds:\n", len(data))
	for _, feed := range data {
		fmt.Printf(">\tName:%s\n\tUrl:%s\n\tAdded By:%s\n\n", feed.Name, feed.Url, usersDict[feed.UserID])
	}
	return nil
}
func HandlerFollow(s *state.State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("No url to follow provided!")
	}
	var feed database.Feed
	var err error
	feed, err = s.Db.GetFeedByUrl(context.Background(), cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("Error:\n\tThat feed hasn't been registered yet! Please add it with 'add <url>'\n")
	}

	if res, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    s.Cfg.CurrentUserId,
		FeedID:    feed.ID,
	}); err != nil {
		return fmt.Errorf("Error:\n\tAn error ocurred registering the follow: %v\n", err)
	} else {
		fmt.Printf("%s now following \"%s\" at \"%s\"!\n", res.UserName, res.FeedName, res.FeedUrl)
		return nil
	}
}
func HandlerFollowing(s *state.State, cmd Command, user database.User) error {
	res, err := s.Db.GetFeedFollowsForUser(context.Background(), s.Cfg.CurrentUserId)
	if err != nil {
		return fmt.Errorf("Error:\n\tAn error ocurred getting the followed feeds: %v\n", err)
	}
	if len(res) < 1 {
		fmt.Println("You're not following any feed yet! Add them with 'add <url>' or follow an existing one with 'follow <url>'")
		return nil
	}
	fmt.Printf("User %s is currently following %d feeds:\n", s.Cfg.CurrentUserName, len(res))
	for _, feed := range res {
		fmt.Printf(">\tName:%s\n\tUrl:%s\n\n", feed.Name, feed.Url)
	}
	return nil
}
func HandlerUnfollow(s *state.State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("Error:\n\tNo url provided\n")
	}

	feed, err := s.Db.GetFeedByUrl(context.Background(), cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("Error:\n\tUrl doesn't match any registered feed\n")
	}
	err = s.Db.DeleteFollow(context.Background(), database.DeleteFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Error:\n\tAn error ocurred deleting the registry: %v\n", err)
	}
	fmt.Println("Succesfuly unsubscribed!")
	return nil
}

func MiddlewareLoggedIn(handler func(s *state.State, cmd Command, user database.User) error) func(*state.State, Command) error {
	return func(s *state.State, cmd Command) error {

		// Check if the user is atually logged in
		user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("User doesn't exist! Please register!")
		}
		// return the function
		return handler(s, cmd, user)
	}

}
