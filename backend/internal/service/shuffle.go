package service

import "time"

// dailySeed is a shuffle seed (see repository.shuffledBy) that stays the same
// all day for the same id, so random boxes don't reshuffle on every reload,
// and changes from one day, or one id, to the next.
func dailySeed(id uint, now time.Time) int64 {
	days := dayOf(now).Unix() / (24 * 60 * 60)
	return (days*7919 + int64(id)) % 1_000_003
}
