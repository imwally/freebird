package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	api_key    = flag.String("key", "", "API consumer key")
	api_secret = flag.String("secret", "", "API consumer secret")
	username   = flag.String("username", "", "Twitter user name")
)

type Friend struct {
	ScreenName string
	ID         int64
}

func GetFriends(client *twitter.Client, username string) ([]Friend, error) {
	user, _, err := client.Users.Show(&twitter.UserShowParams{
		ScreenName: username,
	})
	if err != nil {
		return nil, err
	}

	params := &twitter.FriendListParams{
		UserID:              user.ID,
		ScreenName:          username,
		Cursor:              -1,
		Count:               200,
		SkipStatus:          &[]bool{true}[0],
		IncludeUserEntities: &[]bool{false}[0],
	}

	var friends []Friend
	for params.Cursor != 0 {
		f, _, err := client.Friends.List(params)
		if err != nil {
			return nil, err
		}

		for _, user := range f.Users {
			friends = append(friends, Friend{user.ScreenName, user.ID})
		}

		params.Cursor = f.NextCursor
	}

	return friends, nil
}

func main() {
	flag.Parse()
	if *api_key == "" {
		fmt.Fprintf(os.Stderr, "%s: Please specify a Twitter API consumer key.\n", os.Args[0])
		os.Exit(1)
	}
	if *api_secret == "" {
		fmt.Fprintf(os.Stderr, "%s: Please specify a Twitter API consumer secret.\n", os.Args[0])
		os.Exit(1)
	}
	if *username == "" {
		fmt.Fprintf(os.Stderr, "%s: Please specify a Twitter username.\n", os.Args[0])
		os.Exit(1)
	}

	config := &clientcredentials.Config{
		ClientID:     *api_key,
		ClientSecret: *api_secret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	httpClient := config.Client(oauth2.NoContext)
	client := twitter.NewClient(httpClient)

	friends, err := GetFriends(client, *username)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(friends)
}
