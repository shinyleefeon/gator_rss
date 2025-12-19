package main



import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shinyleefeon/gator_rss/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("Username is required")
	}

	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("User %s does not exist", cmd.args[0])
		os.Exit(1)
	}

	err = s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("Username set to:", s.config.Current_user_name)
	return nil
}

func registerUser(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("Username is required")
	}
	username := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err == nil {
		return fmt.Errorf("User %s already exists", username)
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}

	s.config.SetUser(username)
	fmt.Println("User registered and set as current user:", username)
	fmt.Printf("User details: %+v\n", user)

	return nil
}

func deleteUsers(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		fmt.Println("Error deleting users:", err)
		os.Exit(1)
	}
	fmt.Println("All users deleted from the database.")
	return nil
}

func getAllUsers(s *state, cmd command) error {
	users, err := s.db.GetAllUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Registered users:")
	for _, user := range users {
		if user == s.config.Current_user_name {
			fmt.Println("* ", user, "(current)")
		} else {
			fmt.Println("* ", user)
		}
	}
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err :=  http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "GatorRSS/1.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch feed: %s", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}
	
	return &feed, nil
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get next feed to fetch: %v", err)
	}
	s.db.MarkFeedFetched(context.Background(), nextFeed.Name)
	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %v", err)
	}
	for _, item := range feed.Channel.Item {
		fmt.Println(item.Title)
	}
	return nil
}

func aggregateFeeds(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("Time interval between requests is required")
	}
	time_between_reqs := cmd.args[0]
	fmt.Println("Collecting feeds every %s", time_between_reqs)
	duration, err := time.ParseDuration(time_between_reqs)
	if err != nil {
		return err
	}
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func addFeed(s *state, cmd command, user database.User) error {
	
	if len(cmd.args) < 2 {
		return errors.New("Feed name and URL are required")
	}
	feedName := cmd.args[0]
	feedURL := cmd.args[1]
	fmt.Printf("Adding feed: %s (%s)\n", feedName, feedURL, user.ID)

	feedParams := database.CreateFeedParams{
		Name:   feedName,
		Url:    feedURL,
		UserID: user.ID,
	}
	
	_, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}

	err = followFeed(s, command{args: []string{feedURL}}, user)
	if err != nil {
		return err
	}

	return nil
}

func listFeeds(s *state, cmd command) error {
	getAllFeeds, err := s.db.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Registered feeds:")
	for _, feed := range getAllFeeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		username := user.Name
		fmt.Printf("* %s (%s) [User ID: %s]\n", feed.Name, feed.Url, username)
	}
	return nil
}

func followFeed(s *state, cmd command, user database.User) error {
	
	if len(cmd.args) < 1 {
		return errors.New("Feed url is required to follow")
	}
	feedUrl := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("Feed with URL %s does not exist", feedUrl)
	}
	follow := database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.Name,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), follow)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			fmt.Printf("You are already following feed: %s\n", feed.Name)
			return nil
		}
		return err
	}
	fmt.Printf("Successfully followed feed: %s as user %s\n", feedUrl, s.config.Current_user_name)
	return nil
}

func listFollowing(s *state, cmd command, user database.User) error {
	
	followedFeeds,  err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	fmt.Printf("Feeds followed by %s:\n", user.Name)
	for _, feedFollow := range followedFeeds {
		fmt.Printf("* %s (%s)\n", feedFollow.FeedName, feedFollow.FeedID)
	}
	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(s *state, cmd command) error {
	return func(s *state, cmd command) error {
		username := s.config.Current_user_name
		if username == "" {
			return errors.New("No user logged in. Please login first.")
		}
		user, err := s.db.GetUser(context.Background(), username)
		if err != nil {
			return fmt.Errorf("User %s does not exist", username)
		}
		return handler(s, cmd, user)
	}
}

func unfollowFeed(s *state, cmd command, user database.User) error {
	
	if len(cmd.args) < 1 {
		return errors.New("Feed url is required to unfollow")
	}
	feedUrl := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("Feed with URL %s does not exist", feedUrl)
	}

	unfollowParams := database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feed.Name,
	}

	err = s.db.UnfollowFeed(context.Background(), unfollowParams)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully unfollowed feed: %s as user %s\n", feedUrl, s.config.Current_user_name)
	return nil
}