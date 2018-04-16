package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	api_key    = flag.String("key", "", "API consumer key")
	api_secret = flag.String("secret", "", "API consumer secret")
	snapshot   = flag.Bool("snapshot", false, "Print all friend's IDs to stdout")
	unfollow   = flag.Bool("unfollow", false, "Unfollow your friends")
	username   = flag.String("username", "", "Your user name")
)

func ErrAndExit(s string) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], s)
	os.Exit(1)
}

func ResetSleep(resp *http.Response) {
	if resp.Header["X-Rate-Limit-Remaining"][0] != "0" {
		return
	}

	resetTimeStr := resp.Header["X-Rate-Limit-Reset"][0]
	resetTime, err := strconv.ParseInt(resetTimeStr, 10, 64)
	if err != nil {
		log.Println(err)
	}

	log.Printf("%s: sleeping until %s\n", os.Args[0], time.Unix(resetTime, 0))
	until := time.Until(time.Unix(resetTime, 0))
	time.Sleep(until)
}

func main() {
	flag.Parse()
	if *api_key == "" {
		ErrAndExit("no twitter api consumer key specified")
	}
	if *api_secret == "" {
		ErrAndExit("no twitter api consumer secret specified")
	}
	if *username == "" {
		ErrAndExit("no twitter username specified")
	}

	// Generate a twitter client
	config := &clientcredentials.Config{
		ClientID:     *api_key,
		ClientSecret: *api_secret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	httpClient := config.Client(oauth2.NoContext)
	client := twitter.NewClient(httpClient)

	// Need user's ID to generate friend list params
	user, _, err := client.Users.Show(&twitter.UserShowParams{
		ScreenName: *username,
	})
	if err != nil {
		log.Println(err)
	}

	// Params for friend IDs
	fp := &twitter.FriendIDParams{
		UserID:     user.ID,
		ScreenName: *username,
		Cursor:     -1,
		Count:      5000,
	}

	// Round up those friends (really, just their IDs)
	var friendIDs []int64
	for fp.Cursor != 0 {
		fids, resp, err := client.Friends.IDs(fp)
		if err != nil {
			log.Println(err)
		}
		ResetSleep(resp)

		for _, fid := range fids.IDs {
			friendIDs = append(friendIDs, fid)
		}

		fp.Cursor = fids.NextCursor
	}

	// Print friend's IDs to stdout
	if *snapshot {
		for _, id := range friendIDs {
			fmt.Println(id)
		}
		return
	}

	// Fly away, free bird
	if *unfollow {
		fmt.Printf("Are you sure you want to unfollow everyone? (y/n): ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if scanner.Text() != "y" {
			return
		}

		fmt.Println("BOOM!")
	}
}
