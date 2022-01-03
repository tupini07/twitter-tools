package flow

import (
	"log"
	"math/rand"
	"time"

	"github.com/tupini07/twitter-tools/app_config"
	"github.com/tupini07/twitter-tools/print_utils"
	"github.com/tupini07/twitter-tools/twitter_api"
)

func runFlowStep(step *app_config.FlowStep) {
	if inner := step.FollowAllFollowers; inner != nil {
		twitter_api.FollowAllFollowers(inner.MaxToFollow)
	}

	if inner := step.FollowFollowersOfOthers; inner != nil {
		others := inner.Others

		// if random is specified then shuffle the order of the sources
		if inner.Random {
			rand.Shuffle(len(inner.Others), func(i, j int) {
				others[i], others[j] = others[j], others[i]
			})
		}

		twitter_api.FollowFollowersOfOthers(inner.MaxToFollow, others...)
	}

	if inner := step.UnfollowBadFriends; inner != nil {
		twitter_api.UnfollowBadFriends(inner.MaxToUnfollow)
	}

	if inner := step.Wait; inner != nil {
		s := inner.Seconds * int64(time.Second)
		m := inner.Minutes * int64(time.Minute)
		h := inner.Hours * int64(time.Hour)

		print_utils.WaitWithBar(time.Duration(s+m+h), "Flow sleep")
	}

	if inner := step.WaitUntilDay; inner != nil {
		switch inner.Relative {
		case "tomorrow":
			print_utils.WaitUntilDay(
				time.Now().AddDate(0, 0, 1))
		default:
			log.Fatal("Unknown wait_until_day.relative option:", inner.Relative)
		}
	}
}

func DoFlow(flow *app_config.Flow) {

	for {
		for _, step := range flow.Steps {
			runFlowStep(&step)
		}

		if !flow.Repeat {
			break
		}
	}
}
