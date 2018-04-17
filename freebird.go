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
	"github.com/dghubble/oauth1"
)

var (
	accessToken    = flag.String("token", "", "API token")
	accessSecret   = flag.String("tokenSecret", "", "API token secret")
	consumerKey    = flag.String("consumerKey", "", "API consumer key")
	consumerSecret = flag.String("consumerSecret", "", "API consumer secret")
	snapshot       = flag.Bool("snapshot", false, "Print all friend's IDs to stdout")
	unfollow       = flag.Bool("unfollow", false, "Unfollow your friends")
	username       = flag.String("username", "", "Your user name")
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

func SnapShot(fids []int64) {
	for _, id := range fids {
		fmt.Println(id)
	}
}

func Unfollow(client *twitter.Client, fids []int64) error {
	fmt.Printf("Are you sure you want to unfollow everyone? (y/n): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Text() != "y" {
		return nil
	}

	for _, id := range fids {
		fp := &twitter.FriendshipDestroyParams{
			UserID: id,
		}

		user, resp, err := client.Friendships.Destroy(fp)
		fmt.Println(resp)
		if err != nil {
			return err
		}

		log.Println("Unfollowed:", user.ScreenName)

	}

	return nil
}

func main() {
	flag.Parse()
	if *accessToken == "" {
		ErrAndExit("no twitter api token specified")
	}
	if *accessSecret == "" {
		ErrAndExit("no twitter api token secret specified")
	}
	if *consumerKey == "" {
		ErrAndExit("no twitter api consumer key specified")
	}
	if *consumerSecret == "" {
		ErrAndExit("no twitter api consumer secret specified")
	}
	if *username == "" {
		ErrAndExit("no twitter username specified")
	}

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	configToken := oauth1.NewToken(*accessToken, *accessSecret)
	httpClient := config.Client(oauth1.NoContext, configToken)
	client := twitter.NewClient(httpClient)

	fp := &twitter.FriendIDParams{
		ScreenName: *username,
		Cursor:     -1,
		Count:      5000,
	}

	for fp.Cursor != 0 {
		fids, resp, err := client.Friends.IDs(fp)
		if err != nil {
			ErrAndExit(err.Error())
		}
		ResetSleep(resp)

		if *snapshot {
			SnapShot(fids.IDs)
		}

		if *unfollow && !*snapshot {
			err := Unfollow(client, fids.IDs)
			if err != nil {
				ErrAndExit(err.Error())
			}
		}

		fp.Cursor = fids.NextCursor
	}
}
