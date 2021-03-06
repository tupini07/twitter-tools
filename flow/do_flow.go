package flow

import (
	"math/rand"
	"time"

	"github.com/golang-collections/collections/set"
	"github.com/tupini07/twitter-tools/app_config"
	"github.com/tupini07/twitter-tools/print_utils"
	"github.com/tupini07/twitter-tools/twitter_api"
)

func runFlowStep(flow *app_config.Flow, step *app_config.FlowStep) {
	if inner := step.Random; inner != nil {
		// if random then choose one at random amongst the provided options
		selected := inner.Options[rand.Intn(len(inner.Options))]
		runFlowStep(flow, &selected)
	}

	if inner := step.FollowAllFollowers; inner != nil {
		twitter_api.FollowAllFollowers(inner.MaxToFollow, flow.MaxTotalFollowing)
	}

	if inner := step.FollowFollowersOfOthers; inner != nil {
		// convert list of "others" into a set
		// TODO consider moving this to the actual twitter_api.FollowFollowersOfOthers
		//      function, or when the configuration itself is loaded
		setOthers := set.New()
		for _, other := range inner.Others {
			setOthers.Insert(other)
		}

		othersDedup := []string{}
		setOthers.Do(func(other interface{}) {
			othersDedup = append(othersDedup, other.(string))
		})

		twitter_api.FollowFollowersOfOthers(inner.MaxToFollow,
			flow.MaxTotalFollowing,
			inner.MaxSourcesToPick,
			othersDedup...)
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
			print_utils.Fatal("Unknown wait_until_day.relative option:", inner.Relative)
		}
	}
}

func DoFlow(flow *app_config.Flow) {

	for {
		for _, step := range flow.Steps {
			runFlowStep(flow, &step)
		}

		if !flow.Repeat {
			break
		}
	}
}
