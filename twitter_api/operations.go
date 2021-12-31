package twitter_api

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tupini07/twitter-tools/database"

	"github.com/dghubble/go-twitter/twitter"
)

// Will follow maxNumber followers of the current authenticated user. If
// maxNumber is negative then all followers of the current user will be
// processed.
func FollowAllFollowers(maxNumber int) {
	printTitle("Starting to follow all followers")

	authedUser := GetAuthedUserInformation()

	idsBeingFollowed := GetAllUserIdsBeingFollowed(authedUser.ScreenName)
	idsThatFollowMe := GetFollowersOfUser(authedUser.ScreenName)

	idsToFollow := make([]int64, 0)

	for _, followerId := range idsThatFollowMe {
		areWeFollowingUser := false

		for _, followedId := range idsBeingFollowed {
			// we know we don't need to follow if we're already following us
			if followedId == followerId {
				areWeFollowingUser = true
				break
			}
		}

		if !areWeFollowingUser {
			idsToFollow = append(idsToFollow, followerId)
		}
	}

	printAction(fmt.Sprintf("Starting to follow '%s' unfollowed users", green(len(idsToFollow))))

	processed := 0
	for _, friendId := range idsToFollow {
		processed += 1

		str := fmt.Sprintf("Following user: %s", green(friendId))
		printStepAction(processed, maxNumber, str)
		FollowUserId(friendId)

		now := time.Now()
		database.CreateFriend(&database.Friend{
			UserId:     friendId,
			FollowedOn: &now,
		})

		if processed >= maxNumber {
			break
		}

	}
}

// gets the followers of every user in screenNames and follows up to maxNumber
// of them
func FollowFollowersOfOthers(maxNumber int, screenNames ...string) {
	authedUser := GetAuthedUserInformation()
	processed := 0

	if maxNumber <= 0 {
		log.Fatalf("Error! maxNumbers needs to be greater than 0 but %d was provided\n", maxNumber)
	}

	log.WithField("screen_names", screenNames).Debug("Starting following followers of others")

	for _, sourceName := range screenNames {
		printTitle(fmt.Sprintf("Starting to follow all followers of %s", green(sourceName)))

		sourceLog := log.WithFields(log.Fields{
			"source_name": sourceName,
		})

		channel := make(chan twitter.User)
		go GetFollowersOfUsersStream(sourceName, channel)

		for user := range channel {
			sourceLog.WithField("screen_name", user.ScreenName).Debug("Processing user")

			haveWeUnfollowedThisUserInThePast := false
			if dbEntry := database.GetFriendByUserId(user.ID); dbEntry != nil {
				haveWeUnfollowedThisUserInThePast = dbEntry.UnfollowedOn != nil
			}

			// only send requests to users that we're not yet following and that don't
			// have a private account
			if user.ID != authedUser.ID && !haveWeUnfollowedThisUserInThePast && !user.Following && !user.FollowRequestSent && !user.Protected {
				processed += 1

				str := fmt.Sprintf("Following user: %s", green(user.ScreenName))
				printStepAction(processed, maxNumber, str)

				FollowUserScreenName(user.ScreenName)

				now := time.Now()
				database.CreateFriend(&database.Friend{
					UserId:     user.ID,
					ScreenName: &user.ScreenName,
					FollowedOn: &now,
				})

				if processed >= maxNumber {
					sourceLog.Debug("Done following all followers of source")
					log.WithField("processed", processed).Debug("Done following all followers")
					return
				}
			}
		}

	}

}

// Unfollows maxNumber amount of users being followed by the current user, from
// oldest to newest.
func UnfollowBadFriends(maxNumber int) {
	printTitle("Unfollowing bad friends")

	authedUser := GetAuthedUserInformation()

	idsBeingFollowed := GetAllUserIdsBeingFollowed(authedUser.ScreenName)
	idsThatFollowMe := GetFollowersOfUser(authedUser.ScreenName)

	idsToUnfollow := make([]int64, 0)

	for _, followedId := range idsBeingFollowed {
		// don't mark for unfollow if we've already unfollowed this user in the past
		dbUser := database.GetFriendByUserId(followedId)
		if dbUser != nil {
			// don't unfollow user if we followed them less than a month ago
			isRecentlyFollowed := dbUser.FollowedOn != nil && dbUser.FollowedOn.After(time.Now().Add(-30*24*time.Hour))
			haveWeAlreadyUnfollowedInThePast := dbUser.UnfollowedOn != nil

			if haveWeAlreadyUnfollowedInThePast || isRecentlyFollowed {
				log.WithFields(log.Fields{
					"haveWeAlreadyUnfollowedInThePast": haveWeAlreadyUnfollowedInThePast,
					"isRecentlyFollowed":               isRecentlyFollowed,
					"userId":                           dbUser.UserId,
					"screenName":                       dbUser.ScreenName,
				}).Debug("Skipping unfollow for user")
				continue
			}
		}

		isUserFollowingMe := false

		for _, followerId := range idsThatFollowMe {
			// we know we don't need to unfollow if id follows us
			if followedId == followerId {
				isUserFollowingMe = true
				break
			}
		}

		if !isUserFollowingMe {
			idsToUnfollow = append(idsToUnfollow, followedId)
		}
	}

	printAction(fmt.Sprintf("Starting to unfollow '%s' bad friends", green(len(idsToUnfollow))))

	processed := 0
	for i := len(idsToUnfollow) - 1; i >= 0; i-- {
		badFriendId := idsToUnfollow[i]

		processed += 1

		var str string
		if dbUser := database.GetFriendByUserId(badFriendId); dbUser != nil {
			userIdentifier := string(dbUser.UserId)
			if dbUser.ScreenName != nil {
				userIdentifier = *dbUser.ScreenName
			}

			str = fmt.Sprintf("Unfollowing bad friend '%s' who we had originally followed on '%s'",
				green(userIdentifier),
				green(dbUser.FollowedOn.String()))
		} else {
			str = fmt.Sprintf("Unfollowing bad friend: %s", green(badFriendId))
		}

		printStepAction(processed, maxNumber, str)

		UnfollowUserId(badFriendId)

		now := time.Now()
		database.CreateFriend(&database.Friend{
			UserId:       badFriendId,
			UnfollowedOn: &now,
		})

		if processed >= maxNumber {
			break
		}
	}
}
