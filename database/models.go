package database

import (
	"time"

	"gorm.io/gorm"
)

type Friend struct {
	gorm.Model
	UserId                    int64
	ScreenName                *string
	UserAskedWeDontFollowThem bool
	FollowedOn                *time.Time
	UnfollowedOn              *time.Time
}
