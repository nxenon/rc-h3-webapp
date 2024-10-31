package db

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

const (
	maxLoginAttempts = 3
	rateLimitTTL     = 15 * time.Minute // Set TTL to 24 hours, or adjust as needed
)

func IsUserRateLimitedByUsername(username string) (bool, error) {
	ctx := context.Background()
	userKey := "user_login_attempts:" + username

	// Check the current attempt count
	attempts, err := redisDb.Get(ctx, userKey).Result()
	if err == redis.Nil {
		// User does not exist in Redis, set attempt count to 1
		err = redisDb.Set(ctx, userKey, 1, rateLimitTTL).Err()
		if err != nil {
			return false, fmt.Errorf("error setting initial attempt count: %v", err)
		}
		return false, nil // Not rate limited on the first attempt
	} else if err != nil {
		return false, fmt.Errorf("error retrieving attempt count: %v", err)
	}

	// Convert attempts to an integer
	attemptCount, err := strconv.Atoi(attempts)
	if err != nil {
		return false, fmt.Errorf("error converting attempt count to integer: %v", err)
	}

	// Check if attempts exceed the limit
	if attemptCount > maxLoginAttempts {
		return true, nil // User is rate limited
	}

	return false, nil // User is not rate limited
}

func IncreaseUserRateLimitByUsername(username string) error {
	ctx := context.Background()
	userKey := "user_login_attempts:" + username

	// Increment the attempt count
	err := redisDb.Incr(ctx, userKey).Err()
	if err != nil {
		return fmt.Errorf("error incrementing attempt count: %v", err)
	}

	// Set expiration to ensure the rate limit window, if not already set
	ttl, err := redisDb.TTL(ctx, userKey).Result()
	if err != nil {
		return fmt.Errorf("error getting TTL: %v", err)
	}
	if ttl < 0 {
		err = redisDb.Expire(ctx, userKey, rateLimitTTL).Err()
		if err != nil {
			return fmt.Errorf("error setting TTL: %v", err)
		}
	}

	return nil
}

func ResetAllUserLoginAttemptsToZero() error {
	ctx := context.Background()
	pattern := "user_login_attempts:*"
	iter := redisDb.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		err := redisDb.Del(ctx, iter.Val()).Err()
		if err != nil {
			return fmt.Errorf("error deleting key %s: %v", iter.Val(), err)
		}
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("error during key scan: %v", err)
	}

	return nil
}

func ResetUserRateLimitByUsername(username string) error {
	ctx := context.Background()
	userKey := "user_login_attempts:" + username

	// Set the login attempt count to zero
	err := redisDb.Set(ctx, userKey, 0, 0).Err() // Set to 0 without altering TTL
	if err != nil {
		return fmt.Errorf("error resetting rate limit for user %s: %v", username, err)
	}

	return nil
}
