# 'Cause I'm as free as a bird now

Well this was originally written as a utility to easily unfollow all
the accounts you follow on Twitter. Unfortunately I didn't take the
time to read Twitter's policy on [automation
rules](https://help.twitter.com/en/rules-and-policies/twitter-automation).

> Automated following/unfollowing: You may not follow or unfollow
> Twitter accounts in a bulk, aggressive, or indiscriminate
> manner. Aggressive following is a violation of the Twitter
> Rules. Please also review our following rules and best practices to
> ensure you are in compliance. Note that applications that claim to
> get users more followers are also prohibited under the Twitter
> Rules.

Turns out that after about 100 POST requests the API will revoke your
access token.

## Snapshots

Twitter however is kind enough to allow you to request the IDs of the
accounts you follow in batches of 5000. This is limited to 15 requests
within 15 minutes giving you the ability to fetch 75,000 IDs within
each window.

If for some reason you want a snapshot (list of IDs) of the accounts
you follow then this utility can at least handle that.

```
Usage of freebird:
  -consumerKey string
        API consumer key
  -consumerSecret string
        API consumer secret
  -snapshot
        Print all friend's IDs to stdout
  -token string
        API token
  -tokenSecret string
        API token secret
  -unfollow
        Unfollow your friends
  -username string
        Your user name
```
