package utils

import (
	"fmt"
	"time"
	"zipride/database"
)

// stroing refresh token in redis
func StoreRefresh(tokenid string, userid uint, expiry time.Duration) error {
	key := fmt.Sprintf("refresh_token:%s", tokenid)
	return database.RDB.Set(database.Ctx, key, userid, expiry).Err()
}

// delete refresh token
func DeleteRefresh(tokenid string) error {
	key := fmt.Sprintf("refresh_token:%s", tokenid)
	return database.RDB.Del(database.Ctx, key).Err()
}

// checking refresh token valid
func CheckToken(tokenid string) bool {
	key := fmt.Sprintf("refresh_token:%s", tokenid)

	_, err := database.RDB.Get(database.Ctx, key).Result()

	return err == nil
}
