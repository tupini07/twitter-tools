package twitter_api

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tupini07/twitter-tools/data_utils"
	"github.com/tupini07/twitter-tools/database"

	"github.com/golang-collections/collections/set"
)

// Will follow maxNumber followers of the current authenticated user. If
// maxNumber is negative then all followers of the current user will be
// processed.
func FollowAllFollowers(maxNumber, maxTotalFollowing int) {
	printTitle("Starting to follow all followers")

	authedUser := GetAuthedUserInformation()

	currentNumberOfFollowing := authedUser.FriendsCount
	if currentNumberOfFollowing >= maxTotalFollowing {
		printAction("Not following any users since we're already above 'maxTotalFollowing'")
	}

	idsBeingFollowed := GetAllUserIdsBeingFollowed(authedUser.ScreenName)
	idsThatFollowMe := GetFollowersIDsOfUser(authedUser.ScreenName)

	idsToFollow := make([]int64, 0)

	for _, followerId := range idsThatFollowMe {

		// if we're already following friend on the DB then skip them
		if dbFriend := database.GetFriendByUserId(followerId); dbFriend != nil {
			if dbFriend.FollowedOn != nil {
				continue
			}
		}

		// otherwise, skip them if we're following them on twitter
		for _, followedId := range idsBeingFollowed {
			// we know we don't need to follow if we're already following us
			if followedId == followerId {
				continue
			}
		}

		idsToFollow = append(idsToFollow, followerId)
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

		if processed >= maxNumber || currentNumberOfFollowing+processed >= maxTotalFollowing {
			break
		}

	}
}

// gets the followers of every user in screenNames and follows up to maxNumber
// of them
func FollowFollowersOfOthers(maxNumber, maxTotalFollowing, maxSourcesToPick int, screenNames ...string) {
	printTitle("Following followers of others")

	authedUser := GetAuthedUserInformation()

	currentNumberOfFollowing := authedUser.FriendsCount
	if currentNumberOfFollowing >= maxTotalFollowing {
		printAction("Not following any users since we're already above 'maxTotalFollowing'")
	}

	processed := 0

	if maxNumber <= 0 {
		log.Fatalf("Error! maxNumbers needs to be greater than 0 but %d was provided\n", maxNumber)
	}

	log.WithField("screen_names", screenNames).Debug("Starting following followers of others")

	allFollowerIdsOfOthersSet := set.New()

	data_utils.ShuffleArrayInplace(screenNames)

	pickedSources := 0
	for _, sourceName := range screenNames {
		log.WithField("source_name", sourceName).Debug("Getting follower ids of other")
		for _, id := range GetFollowersIDsOfUser(sourceName) {
			allFollowerIdsOfOthersSet.Insert(id)
		}
		pickedSources += 1

		// if maxSources is 0 or -1 then pick all sources
		if maxSourcesToPick > 0 && pickedSources >= maxSourcesToPick {
			break
		}
	}

	myFollowersIds := GetFollowersIDsOfUser(authedUser.ScreenName)
	myFollowersSet := set.New()
	for _, mFID := range myFollowersIds {
		myFollowersSet.Insert(mFID)
	}

	// so that we don't try to follow ourselves
	myFollowersSet.Insert(authedUser.ID)

	potentialIdsToFollow := make([]int64, 0)

	// consider followers which are not my followers
	allFollowerIdsOfOthersSet.Difference(myFollowersSet).Do(func(i interface{}) {
		potentialIdsToFollow = append(potentialIdsToFollow, i.(int64))
	})

	data_utils.ShuffleArrayInplace(potentialIdsToFollow)

	for _, userId := range potentialIdsToFollow {
		log.WithField("user_id", userId).Debug("Processing user")

		haveWeUnfollowedThisUserInThePast := false
		if dbEntry := database.GetFriendByUserId(userId); dbEntry != nil {
			haveWeUnfollowedThisUserInThePast = dbEntry.UnfollowedOn != nil
		}

		if !haveWeUnfollowedThisUserInThePast {
			processed += 1

			str := fmt.Sprintf("Following user: %s", green(userId))
			printStepAction(processed, maxNumber, str)

			FollowUserId(userId)

			now := time.Now()
			database.CreateFriend(&database.Friend{
				UserId:     userId,
				FollowedOn: &now,
			})

			if processed >= maxNumber || currentNumberOfFollowing+processed >= maxTotalFollowing {
				log.WithField("processed", processed).Debug("Done following all followers")
				return
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
	idsThatFollowMe := GetFollowersIDsOfUser(authedUser.ScreenName)

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

	actualAmountToUnfollow := len(idsToUnfollow)
	printAction(fmt.Sprintf("Starting to unfollow '%s' bad friends", green(actualAmountToUnfollow)))

	processed := 0
	for i := len(idsToUnfollow) - 1; i >= 0; i-- {
		badFriendId := idsToUnfollow[i]

		processed += 1

		var str string
		if dbUser := database.GetFriendByUserId(badFriendId); dbUser != nil {
			userIdentifier := fmt.Sprint(dbUser.UserId)
			if dbUser.ScreenName != nil {
				userIdentifier = *dbUser.ScreenName
			}

			str = fmt.Sprintf("Unfollowing bad friend '%s' who we had originally followed on '%s'",
				green(userIdentifier),
				green(dbUser.FollowedOn.String()))
		} else {
			str = fmt.Sprintf("Unfollowing bad friend: %s", green(badFriendId))
		}

		printStepAction(processed, actualAmountToUnfollow, str)

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
