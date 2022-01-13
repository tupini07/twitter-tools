package database

// save friend as is. To be used after modifying a friend instance that came
// from the DB
func SaveFriend(friend *Friend) {
	db := getDbInstance()
	db.Save(friend)
}

// save friend but if friend exist then only override the value of the
// properties that are set in the incoming `friend`
func CreateFriend(friend *Friend) {
	db := getDbInstance()

	existingEntry := GetFriendByUserId(friend.UserId)
	if existingEntry == nil {
		db.Create(friend)
	} else {
		if friend.ScreenName != nil {
			existingEntry.ScreenName = friend.ScreenName
		}

		if friend.FollowedOn != nil {
			existingEntry.FollowedOn = friend.FollowedOn
		}

		if friend.UnfollowedOn != nil {
			existingEntry.UnfollowedOn = friend.UnfollowedOn
		}

		// once set this cannot be unset
		existingEntry.UserAskedWeDontFollowThem = existingEntry.UserAskedWeDontFollowThem || friend.UserAskedWeDontFollowThem

		db.Save(existingEntry)
	}
}

func GetFriendByUserId(userId int64) *Friend {
	db := getDbInstance()
	friend := &Friend{}
	queryRes := db.Where(&Friend{UserId: userId}).First(friend)

	if queryRes.Error == nil {
		return friend
	} else {
		return nil
	}
}

func GetFriendByDbId(id int) *Friend {
	db := getDbInstance()
	res := &Friend{}
	db.First(res, id)
	return res
}
