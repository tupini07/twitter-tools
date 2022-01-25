package twitter_api

import (
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tupini07/twitter-tools/database"

	"github.com/dghubble/go-twitter/twitter"
)

// Will get all followers of the currently authenticated user and pipes them
// through followerChannel and close followerChannel once all followers have
// been returned have been processed
func GetFollowersOfUsersStream(screenName string, followerChannel chan twitter.User) {
	mLog := log.WithField("screenName", screenName)
	mLog.Debug("Starting to get followers of user stream")

	const countParam = 200

	cli := getApiClient()

	processed := 0
	var followerChunk *twitter.Followers

	for {
		followerChunk = makeTimeoutHandledRequest(time.Minute, func() (interface{}, *http.Response, error) {
			var nextCursor int64
			if followerChunk != nil {
				nextCursor = followerChunk.NextCursor
			}

			mLog.WithField("next_cursor", nextCursor).Debug("Getting followers from API")

			return cli.Followers.List(&twitter.FollowerListParams{
				ScreenName: screenName,
				Count:      countParam,
				Cursor:     nextCursor,
				SkipStatus: twitter.Bool(true),
			})
		}).(*twitter.Followers)

		mLog.WithField("followers_in_chunk", len(followerChunk.Users)).Debug("Got new chunk of followers from API")

		for _, follower := range followerChunk.Users {
			followerChannel <- follower
			processed += 1
		}

		if len(followerChunk.Users) < 200 {
			break
		}
	}

	mLog.WithField("processed", processed).Debug("Finished getting followers from API")

	close(followerChannel)
}

// Returns a list of all users that are following the currently authenticated
// user
func GetAllFollowersOfUser(screenName string, maxUsers int) (result []twitter.User) {
	mLog := log.WithField("screenName", screenName)
	mLog.Debug("Starting to get all followers of user")

	channel := make(chan twitter.User)

	go GetFollowersOfUsersStream(screenName, channel)

	mLog.Debug("Draining followers from channel")
	for follower := range channel {
		result = append(result, follower)

		if len(result) >= maxUsers {
			break
		}
	}

	mLog.WithField("amount_followers", len(result)).Debug("Finished draining followers from channel")
	return result
}

// Will get all the users being followed by the currently authenticated user and
// pipes them through followingChannel and close followingChannel once all users
// have been returned OR maxUsers have been processed
func GetUsersBeingFollowedStream(screenName string, followingChannel chan twitter.User) {
	mLog := log.WithField("screenName", screenName)
	mLog.Debug("Starting to get friends of user stream")

	const countParam = 200
	cli := getApiClient()

	processed := 0
	var followingChunk *twitter.Friends

	for {
		followingChunk = makeTimeoutHandledRequest(time.Minute, func() (interface{}, *http.Response, error) {
			var nextCursor int64
			if followingChunk != nil {
				nextCursor = followingChunk.NextCursor
			}

			mLog.WithField("next_cursor", nextCursor).Debug("Getting friends from API")

			return cli.Friends.List(&twitter.FriendListParams{
				ScreenName: screenName,
				Count:      countParam,
				Cursor:     nextCursor,
			})
		}).(*twitter.Friends)

		mLog.WithField("friends_in_chunk", len(followingChunk.Users)).Debug("Got new chunk of friends from API")

		for _, follower := range followingChunk.Users {
			followingChannel <- follower
			processed += 1
		}

		if len(followingChunk.Users) < 200 {
			break
		}
	}

	mLog.WithField("processed", processed).Debug("Finished getting friends from API")

	close(followingChannel)
}

// Returns a list of all users that are being followed by the current
// authenticated user
func GetUsersBeingFollowed(screenName string, maxUsers int) (result []twitter.User) {
	mLog := log.WithField("screenName", screenName)
	mLog.Debug("Starting to get all friends of user")

	channel := make(chan twitter.User)

	go GetUsersBeingFollowedStream(screenName, channel)

	mLog.Debug("Draining friends from channel")
	for follower := range channel {
		result = append(result, follower)

		if len(result) >= maxUsers {
			break
		}
	}

	mLog.WithField("amount_followers", len(result)).Debug("Finished draining friends from channel")
	return result
}

func GetAllUserIdsBeingFollowed(screenName string) (result []int64) {
	mLog := log.WithField("screenName", screenName)
	mLog.Debug("Starting to get all friend IDs being followed")

	cli := getApiClient()
	const countParam = 5000

	var followingChunk *twitter.FriendIDs

	for {
		followingChunk = makeTimeoutHandledRequest(time.Minute, func() (interface{}, *http.Response, error) {
			var nextCursor int64
			if followingChunk != nil {
				nextCursor = followingChunk.NextCursor
			}

			mLog.WithField("next_cursor", nextCursor).Debug("Getting friend IDs from API")

			return cli.Friends.IDs(&twitter.FriendIDParams{
				ScreenName: screenName,
				Count:      countParam,
				Cursor:     nextCursor,
			})
		}).(*twitter.FriendIDs)

		mLog.WithField("followers_in_chunk", len(followingChunk.IDs)).Debug("Got new chunk of friend IDs from API")

		result = append(result, followingChunk.IDs...)

		if len(followingChunk.IDs) < countParam {
			break
		}
	}

	mLog.WithField("processed", len(result)).Debug("Finished getting friend IDs from API")

	return result
}

func FollowUserScreenName(screenName string) {
	log.WithField("screen_name", screenName).Debug("Following user")

	cli := getApiClient()
	makeTimeoutHandledRequest(time.Minute, func() (interface{}, *http.Response, error) {
		data, resp, err := cli.Friendships.Create(&twitter.FriendshipCreateParams{
			ScreenName: screenName,
		})

		if err != nil {
			if strings.Contains(err.Error(), "160 You've already requested to follow") {
				// 160 You've already requested to follow
				printActionLog("Skipping user since follow request has already been sent")
				return nil, resp, nil
			}

			if strings.Contains(err.Error(), "108 Cannot find specified user") {
				printAction(yellow("Skipping user since Twitter says it can't find them"))
				return nil, resp, nil
			}

			if strings.Contains(err.Error(), "162 You have been blocked from following this account at the request of the user") {
				printAction(yellow("Skipping user since they have asked we don't follow them"))
				return nil, resp, nil
			}
		}

		return data, resp, err
	})
}

func FollowUserId(userId int64) {
	// check if user hasn't asked that we don't dollow them
	if dbEntry := database.GetFriendByUserId(userId); dbEntry != nil {
		if dbEntry.UserAskedWeDontFollowThem {
			printAction(yellow("Skipping user since they have asked we don't follow them"))
			return
		}
	}

	log.WithField("user_id", userId).Debug("Following user")

	cli := getApiClient()
	makeTimeoutHandledRequest(time.Minute, func() (interface{}, *http.Response, error) {
		data, resp, err := cli.Friendships.Create(&twitter.FriendshipCreateParams{
			UserID: userId,
		})

		if err != nil {
			if strings.Contains(err.Error(), "160 You've already requested to follow") {
				// 160 You've already requested to follow
				printActionLog(yellow("Skipping user since follow request has already been sent"))
				return nil, resp, nil
			}

			if strings.Contains(err.Error(), "108 Cannot find specified user") {
				printAction(yellow("Skipping user since Twitter says it can't find them"))
				return nil, resp, nil
			}

			if strings.Contains(err.Error(), "162 You have been blocked from following this account") {
				database.CreateFriend(&database.Friend{
					UserId:                    userId,
					UserAskedWeDontFollowThem: true,
				})
				printAction(yellow("Skipping user since they have asked we don't follow them"))
				return nil, resp, nil
			}
		}

		return data, resp, err
	})
}

func UnfollowUserScreenName(screenName string) {
	log.WithField("screen_name", screenName).Debug("Unfollowing user")

	cli := getApiClient()
	makeTimeoutHandledRequest(time.Minute, func() (interface{}, *http.Response, error) {
		return cli.Friendships.Destroy(&twitter.FriendshipDestroyParams{
			ScreenName: screenName,
		})
	})
}
func UnfollowUserId(userId int64) {
	log.WithField("user_id", userId).Debug("Unfollowing user")

	cli := getApiClient()
	makeTimeoutHandledRequest(time.Minute, func() (interface{}, *http.Response, error) {
		return cli.Friendships.Destroy(&twitter.FriendshipDestroyParams{
			UserID: userId,
		})
	})
}

var authedUserInstace *twitter.User

func GetAuthedUserInformation() *twitter.User {
	log.Debug("Getting Authed user information")

	if authedUserInstace == nil {
		cli := getApiClient()
		authedUserInstace = makeTimeoutHandledRequest(10*time.Second, func() (interface{}, *http.Response, error) {
			return cli.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{
				IncludeEntities: twitter.Bool(false),
				SkipStatus:      twitter.Bool(true),
			})
		}).(*twitter.User)
	}

	return authedUserInstace
}

func GetFollowersIDsOfUser(screenName string) (result []int64) {
	mLog := log.WithField("screenName", screenName)
	mLog.Debug("Starting to get all followers of user")

	cli := getApiClient()
	const countParam = 5000

	var followersChunk *twitter.FollowerIDs

	for {
		followersChunk = makeTimeoutHandledRequest(time.Minute, func() (interface{}, *http.Response, error) {
			var nextCursor int64
			if followersChunk != nil {
				nextCursor = followersChunk.NextCursor
			}

			mLog.WithField("next_cursor", nextCursor).Debug("Getting follower IDs from API")

			return cli.Followers.IDs(&twitter.FollowerIDParams{
				ScreenName: screenName,
				Count:      countParam,
				Cursor:     nextCursor,
			})
		}).(*twitter.FollowerIDs)

		mLog.WithField("followers_in_chunk", len(followersChunk.IDs)).Debug("Got new chunk of follower IDs from API")

		result = append(result, followersChunk.IDs...)

		if len(followersChunk.IDs) < countParam {
			break
		}
	}

	mLog.WithField("processed", len(result)).Debug("Finished getting follower IDs from API")

	return result
}

func GetFriendship(sourceId, targetId int64) twitter.RelationshipSource {
	log.WithFields(log.Fields{
		"sourceId": sourceId,
		"targetId": targetId,
	}).Debug("Getting friendship relation")

	cli := getApiClient()
	relationship := makeTimeoutHandledRequest(3*time.Second, func() (interface{}, *http.Response, error) {
		return cli.Friendships.Show(&twitter.FriendshipShowParams{
			SourceID: sourceId,
			TargetID: targetId,
		})
	}).(*twitter.Relationship)

	return relationship.Source
}
