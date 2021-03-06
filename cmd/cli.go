package cmd

import (
	"os"

	"github.com/tupini07/twitter-tools/app_config"
	"github.com/tupini07/twitter-tools/flow"
	"github.com/tupini07/twitter-tools/print_utils"
	"github.com/tupini07/twitter-tools/twitter_api"
	"github.com/urfave/cli/v2"
)

func setupLogger(c *cli.Context) {
	logOutput := c.String("log-output")
	print_utils.SetupLogger(logOutput)
}

func RunCli() {
	app := &cli.App{
		Name:    "twitter-tools",
		Version: "1.0.0",
		Usage:   "Collection of tools to manage a Twitter account",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-output",
				Aliases: []string{"l"},
				Usage:   "If provided, the path of the file where output will be logged",
				Value:   "",
			},
		},
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
					setupLogger(c)
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
					setupLogger(c)
					amount := c.Int("amount")
					maxTotalFollowing := c.Int("max-total-following")
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
						Aliases: []string{"mf"},
						Usage:   "Maximum number of users you want to be following",
						Value:   4500,
					},
					&cli.IntFlag{
						Name:    "max-sources-to-pick",
						Aliases: []string{"mp"},
						Usage:   "Maximum number of sources to pick at random from the provided 'others'",
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
					setupLogger(c)
					amount := c.Int("amount")
					maxTotalFollowing := c.Int("max-total-following")
					maxSourcesToPick := c.Int("max-sources-to-pick")
					others := c.StringSlice("other")

					twitter_api.FollowFollowersOfOthers(amount,
						maxTotalFollowing,
						maxSourcesToPick,
						others...)

					return nil
				},
			},
			{
				Name:    "do-flow",
				Aliases: []string{"df"},
				Usage:   "Performs the flow actions defined in config.yml",
				Action: func(c *cli.Context) error {
					setupLogger(c)
					cnf := app_config.GetConfig()
					flow.DoFlow(cnf.Flow)

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		print_utils.Fatal(err)
	}
}
