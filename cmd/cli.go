package cmd

import (
	"log"
	"os"

	"github.com/tupini07/twitter-tools/app_config"
	"github.com/tupini07/twitter-tools/flow"
	"github.com/tupini07/twitter-tools/twitter_api"
	"github.com/urfave/cli/v2"
)

func RunCli() {
	app := &cli.App{
		Name:    "twitter-tools",
		Version: "0.0.3",
		Usage:   "Collection of tools to manage a Twitter account",
		Commands: []*cli.Command{
			{
				Name:    "unfollow-bad-friends",
				Aliases: []string{"ubf"},
				Usage:   "Unfollows bad friends starting from the oldest friendship",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "amount",
						Aliases: []string{"a"},
						Usage:   "Amount of users to unfollow",
						Value:   400,
					},
				},
				Action: func(c *cli.Context) error {
					amount := c.Int("amount")
					twitter_api.UnfollowBadFriends(amount)
					return nil
				},
			},
			{
				Name:    "follow-all-followers",
				Aliases: []string{"faf"},
				Usage:   "Ensure followers of the current user are being followed",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "amount",
						Aliases: []string{"a"},
						Usage:   "Amount of users to follow",
						Value:   200,
					},
					&cli.IntFlag{
						Name:    "max-total-following",
						Aliases: []string{"m"},
						Usage:   "Maximum number of users you want to be following",
						Value:   4500,
					},
				},
				Action: func(c *cli.Context) error {
					amount := c.Int("amount")
					maxTotalFollowing := c.Int("max-total-followers")
					twitter_api.FollowAllFollowers(amount, maxTotalFollowing)
					return nil
				},
			},
			{
				Name:    "follow-followers-of-other",
				Aliases: []string{"ffoo"},
				Usage:   "Follow followers of other(s) users",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "amount",
						Aliases: []string{"a"},
						Usage:   "Amount of users to follow",
						Value:   200,
					},
					&cli.IntFlag{
						Name:    "max-total-following",
						Aliases: []string{"m"},
						Usage:   "Maximum number of users you want to be following",
						Value:   4500,
					},
					&cli.StringSliceFlag{
						Name:     "other",
						Aliases:  []string{"o"},
						Usage:    "Twitter handler of user to take followers from",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					amount := c.Int("amount")
					maxTotalFollowing := c.Int("max-total-followers")
					others := c.StringSlice("other")

					twitter_api.FollowFollowersOfOthers(amount, maxTotalFollowing, others...)
					return nil
				},
			},
			{
				Name:    "do-flow",
				Aliases: []string{"df"},
				Usage:   "Performs the flow actions defined in config.yml",
				Action: func(c *cli.Context) error {
					cnf := app_config.GetConfig()
					flow.DoFlow(cnf.Flow)

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
