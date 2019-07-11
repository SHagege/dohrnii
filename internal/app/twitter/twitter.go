package twitter

import (
	"log"
	"flag"
	"time"
	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2"
	"github.com/coreos/pkg/flagutil"
)

// Tweet represents all the proprieties of a tweet
type Tweet struct {
	User		string		`json:"user"`
	Content 	string		`json:"content"`
	Timestamp	time.Time	`json:"timestamp"`
}

// GetTweets get the latest tweets from a specific user
func GetTweets() []Tweet {
	flags := flag.NewFlagSet("app-auth", flag.ExitOnError)
	accessToken := flags.String("app-access-token", "AAAAAAAAAAAAAAAAAAAAAKDv%2FAAAAAAADfoFqiqL7romy%2Bv3I0d0vBsdR0Y%3DNF6444DjREZu76Wi492UUJvfg8wT0F8UgCfOaBhhN2MV72u5GY", "Twitter Application Access Token")
	flagutil.SetFlagsFromEnv(flags, "TWITTER")

	if *accessToken == "" {
		log.Fatal("Application Access Token required")
	}

	config := &oauth2.Config{}
	token := &oauth2.Token{AccessToken: *accessToken}
	// OAuth2 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth2.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	userTimelineParams := &twitter.UserTimelineParams{ScreenName: "NakamotoQuotes", Count: 1}
	tweets, _, _ := client.Timelines.UserTimeline(userTimelineParams)

	tweet := []Tweet{}
	for i := 0; i < len(tweets); i++ {
		timeTweet, _ := tweets[i].CreatedAtTime()
		tweet = append(tweet, Tweet{User: tweets[i].User.Name, Content: tweets[i].Text, Timestamp: timeTweet})
	}
	return tweet
}