package main

import (
	"errors"
	"fmt"

	"context"
	"time"

	"strconv"

	"github.com/OmarJarbou/Gator/internal/config"
	"github.com/OmarJarbou/Gator/internal/database"
	"github.com/google/uuid"
)

type commands struct {
	cmdsMap map[string]func(*state, command) error
}

func (cmds *commands) commandMapping(cmd string, args []string) command {
	var cmnd command
	_, ok := cmds.cmdsMap[cmd]
	if ok {
		return command{
			Name:      cmd,
			Arguments: args,
		}
	}
	return cmnd
}

func (cmds *commands) register(name string, f func(*state, command) error) {
	cmds.cmdsMap[name] = f
}

func (cmds *commands) run(state *state, cmd command) error {
	_, ok := cmds.cmdsMap[cmd.Name]
	if !ok {
		return errors.New("COMMAND NOT FOUND")
	}
	err := cmds.cmdsMap[cmd.Name](state, cmd)
	if err != nil {
		return err
	}
	return nil
}

func handleLogin(state *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		return errors.New("THE LOGIN HANDLER EXPECTS A SINGLE ARGUMENT, THE USERNAME")
	}
	_, err := state.DBQueries.GetUser(context.Background(), cmd.Arguments[0])
	if err != nil {
		return errors.New("USER NOT FOUND: " + err.Error())
	}
	err2 := config.SetUser(cmd.Arguments[0], state.Config)
	if err2 != nil {
		return err2
	}
	fmt.Println("User set to", cmd.Arguments[0])
	return nil
}

func handleRegister(state *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		return errors.New("THE REGISTER HANDLER EXPECTS A SINGLE ARGUMENT, THE USERNAME")
	}
	_, err := state.DBQueries.GetUser(context.Background(), cmd.Arguments[0])
	if err == nil {
		return errors.New("USER WITH THIS NAME '" + cmd.Arguments[0] + "' ALREADY EXISTS")
	}

	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Arguments[0],
	}

	usr, err2 := state.DBQueries.CreateUser(context.Background(), user)
	if err2 != nil {
		return errors.New("FAILED TO CREATE USER: " + err2.Error())
	}

	state.Config.CurrentUserName = cmd.Arguments[0]
	fmt.Println("User", usr.Name, "created successfully!")
	fmt.Println("Created at:", usr.CreatedAt)

	err3 := config.SetUser(cmd.Arguments[0], state.Config)
	if err3 != nil {
		return err3
	}
	fmt.Println("User set to", cmd.Arguments[0])

	return nil
}

func handleReset(state *state, cmd command) error {
	if len(cmd.Arguments) != 0 {
		return errors.New("THE RESET HANDLER DOES NOT EXPECT ANY ARGUMENTS")
	}
	err := state.DBQueries.ClearUsers(context.Background())
	if err != nil {
		return errors.New("FAILED TO RESET DB (CLEAR USERS ISSUE): " + err.Error())
	}
	err2 := state.DBQueries.ClearFeeds(context.Background())
	if err2 != nil {
		return errors.New("FAILED TO RESET DB (CLEAR FEEDS ISSUE): " + err2.Error())
	}
	fmt.Println("DB reset successfully!")
	return nil
}

func handleListUsers(state *state, cmd command) error {
	if len(cmd.Arguments) != 0 {
		return errors.New("THE LIST USERS HANDLER DOES NOT EXPECT ANY ARGUMENTS")
	}

	users, err := state.DBQueries.GetUsers(context.Background())
	if err != nil {
		return errors.New("FAILED TO GET USERS: " + err.Error())
	}
	for _, user := range users {
		if state.Config.CurrentUserName == user.Name {
			fmt.Println("*", user.Name, "(current)")
		} else {
			fmt.Println("*", user.Name)
		}
	}
	return nil
}

func handleAggregate(state *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		return errors.New("THE AGGREGATE HANDLER EXPECTS A SINGLE ARGUMENT, THE TIME BETWEEN REQUESTS")
	}

	timeBetweenRequests := cmd.Arguments[0]
	timeBetweenRequestsDuration, err := time.ParseDuration(timeBetweenRequests)
	if err != nil {
		return errors.New("FAILED TO PARSE TIME BETWEEN REQUESTS: " + err.Error())
	}
	fmt.Println("Collecting feeds every:", timeBetweenRequestsDuration)

	ticker := time.NewTicker(timeBetweenRequestsDuration)
	for ; ; <-ticker.C {
		err2 := scrapeFeeds(context.Background(), state)
		if err2 != nil {
			return err2
		}
		fmt.Println("--------------------------------")
	}
}

func handleAddFeed(state *state, cmd command, currentUser database.User) error {
	if len(cmd.Arguments) != 2 {
		return errors.New("THE ADD FEED HANDLER EXPECTS TWO ARGUMENT, THE NAME AND THE URL OF THE FEED")
	}

	feed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Arguments[0],
		Url:       cmd.Arguments[1],
	}

	createdFeed, err2 := state.DBQueries.CreateFeed(context.Background(), feed)
	if err2 != nil {
		return errors.New("FAILED TO CREATE FEED: " + err2.Error())
	}

	fmt.Println("Feed", createdFeed.Name, "added successfully!")
	fmt.Println("Created at:", createdFeed.CreatedAt)
	fmt.Println("URL:", createdFeed.Url)

	feedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    createdFeed.ID,
	}

	createdFeedFollow, err3 := state.DBQueries.CreateFeedFollow(context.Background(), feedFollow)
	if err3 != nil {
		return errors.New("FAILED TO CREATE FEED FOLLOW: " + err3.Error())
	}

	fmt.Println("Feed", createdFeedFollow.FeedName, "followed successfully by user", createdFeedFollow.UserName+"!")

	return nil
}

func handleListFeeds(state *state, cmd command) error {
	if len(cmd.Arguments) != 0 {
		return errors.New("THE LIST FEEDS HANDLER DOES NOT EXPECT ANY ARGUMENTS")
	}

	feeds, err := state.DBQueries.GetFeeds(context.Background())
	if err != nil {
		return errors.New("FAILED TO GET FEEDS: " + err.Error())
	}

	for _, feed := range feeds {
		fmt.Println("Feed #" + feed.ID.String())
		fmt.Println("Name:", feed.Name)
		fmt.Println("URL:", feed.Url)
		fmt.Println("Created at:", feed.CreatedAt)
		fmt.Println("Updated at:", feed.UpdatedAt)
	}
	return nil
}

func handleFollowFeed(state *state, cmd command, currentUser database.User) error {
	if len(cmd.Arguments) != 1 {
		return errors.New("THE FOLLOW FEED HANDLER EXPECTS A SINGLE ARGUMENT, THE FEED URL")
	}

	feedURL := cmd.Arguments[0]
	feed, err := state.DBQueries.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return errors.New("FAILED TO GET FEED: " + err.Error())
	}

	feedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	}

	createdFeedFollow, err3 := state.DBQueries.CreateFeedFollow(context.Background(), feedFollow)
	if err3 != nil {
		return errors.New("FAILED TO CREATE FEED FOLLOW: " + err3.Error())
	}

	fmt.Println("Feed", createdFeedFollow.FeedName, "followed successfully by user", createdFeedFollow.UserName+"!")

	return nil
}

func handleListFollowing(state *state, cmd command, currentUser database.User) error {
	if len(cmd.Arguments) != 0 {
		return errors.New("THE LIST FOLLOWING HANDLER DOES NOT EXPECT ANY ARGUMENTS")
	}

	followedFeeds, err := state.DBQueries.GetFeedFollowsForUser(context.Background(), currentUser.Name)
	if err != nil {
		return errors.New("FAILED TO GET FOLLOWED FEEDS: " + err.Error())
	}

	fmt.Println("Feeds followed by", currentUser.Name+":")
	for _, feed := range followedFeeds {
		fmt.Println(feed.FeedName)
	}
	return nil
}

func handleUnfollowFeed(state *state, cmd command, currentUser database.User) error {
	if len(cmd.Arguments) != 1 {
		return errors.New("THE UNFOLLOW FEED HANDLER EXPECTS A SINGLE ARGUMENT, THE FEED URL")
	}

	feedURL := cmd.Arguments[0]
	deleteFeedFollowParams := database.DeleteFeedFollowParams{
		Url:    feedURL,
		UserID: currentUser.ID,
	}
	err := state.DBQueries.DeleteFeedFollow(context.Background(), deleteFeedFollowParams)
	if err != nil {
		return errors.New("FAILED TO DELETE FEED FOLLOW: " + err.Error())
	}

	fmt.Println("Feed", feedURL, "unfollowed successfully by user", currentUser.Name+"!")

	return nil
}

func handleBrowsePosts(state *state, cmd command, currentUser database.User) error {
	if len(cmd.Arguments) != 0 && len(cmd.Arguments) != 1 {
		return errors.New("THE BROWSE POSTS HANDLER DOES NOT EXPECT ANY ARGUMENTS OR A SINGLE ARGUMENT, THE LIMIT")
	}

	var limit int64
	var err error
	if len(cmd.Arguments) == 0 {
		limit = 2
	} else {
		limit, err = strconv.ParseInt(cmd.Arguments[0], 10, 64)
		if err != nil {
			return errors.New("FAILED TO PARSE LIMIT: " + err.Error())
		}
	}

	params := database.GetPostsForUserParams{
		UserID: currentUser.ID,
		Limit:  int32(limit),
	}

	posts, err2 := state.DBQueries.GetPostsForUser(context.Background(), params)
	if err2 != nil {
		return errors.New("FAILED TO GET POSTS: " + err2.Error())
	}

	for _, post := range posts {
		fmt.Println("Post #" + post.ID.String())
		fmt.Println("Feed ID:", post.FeedID)
		fmt.Println("Title:", post.Title)
		fmt.Println("URL:", post.Url)
		fmt.Println("Description:", post.Description)
		fmt.Println("Published at:", post.PublishedAt)
		fmt.Println("Created at:", post.CreatedAt)
		fmt.Println("Updated at:", post.UpdatedAt)
		fmt.Println("--------------------------------")
	}
	return nil
}
