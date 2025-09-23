package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserFollowers struct {
	FollowerID *uuid.UUID `db:"follower_id"`
	FolloweeID *uuid.UUID `db:"followee_id"`
	CreatedAt  *time.Time `db:"created_at"`
}
