package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	api_key    = flag.String("key", "", "API consumer key")
	api_secret = flag.String("secret", "", "API consumer secret")
	username   = flag.String("username", "", "Twitter user name")
)

func main() {
	flag.Parse()

	config := &clientcredentials.Config{
		ClientID:     *api_key,
		ClientSecret: *api_secret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	httpClient := config.Client(oauth2.NoContext)
	client := twitter.NewClient(httpClient)

	user, _, err := client.Users.Show(&twitter.UserShowParams{
		ScreenName: *username,
	})
	if err != nil {
		log.Println(err)
	}

	params := &twitter.FriendListParams{
		UserID:              user.ID,
		ScreenName:          *username,
		Cursor:              0,
		Count:               200,
		SkipStatus:          &[]bool{true}[0],
		IncludeUserEntities: &[]bool{false}[0],
	}

	var users []twitter.User
	var cursor int64 = 1
	for cursor != 0 {
		friends, _, err := client.Friends.List(params)
		if err != nil {
			log.Println(err)
		}

		for _, user := range friends.Users {
			users = append(users, user)
		}

		params.Cursor = friends.NextCursor
		cursor = params.Cursor
	}

	for i, user := range users {
		fmt.Printf("%d\t%-20s%d\n", i, user.ScreenName, user.ID)
	}
}
