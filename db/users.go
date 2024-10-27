package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/nxenon/rc-h3-webapp/models"
	"github.com/redis/go-redis/v9"
	"strconv"
)

func insertTestUserInDb(username, passwordHash string, balance float64) error {
	if dbType == "mysql" {
		return insertTestUserInMySqlDb(username, passwordHash, balance)
	} else if dbType == "redis" {
		return insertTestUserInMyRedisDb(username, passwordHash, balance)
	} else {
		panic(fmt.Sprintf("invalid DB type: %s", dbType))
	}
}

func insertTestUserInMySqlDb(username, passwordHash string, balance float64) error {
	// Check if a user exists and insert if it doesn't

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM USERS WHERE USERNAME = ?)`

	// Check if the user exists
	err := mysqldb.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if user exists: %v", err)
	}

	if exists {
		// If user exists, update the balance
		updateQuery := `UPDATE USERS SET BALANCE = ? WHERE USERNAME = ?`
		_, err = mysqldb.Exec(updateQuery, balance, username)
		if err != nil {
			return fmt.Errorf("failed to update user balance: %v", err)
		}
		fmt.Println("User balance in db restarted:", username)
	} else {
		// If user does not exist, insert a new user
		insertQuery := `INSERT INTO USERS (USERNAME, PASSWORD_HASH, BALANCE) VALUES (?, ?, ?)`
		_, err = mysqldb.Exec(insertQuery, username, passwordHash, balance)
		if err != nil {
			return fmt.Errorf("failed to insert user: %v", err)
		}
		fmt.Println("User inserted in db:", username)
	}

	return insertUserBalanceInRedisDb(username, balance)
}

func insertTestUserInMyRedisDb(username, passwordHash string, balance float64) error {
	// todo
	return nil
}

func GetUserObjectByUsername(username string) (models.UserObject, error) {
	if dbType == "mysql" {
		return GetUserObjectByUsernameInMySqlDb(username)
	} else if dbType == "redis" {
		return GetUserObjectByUsernameInRedisDb(username)
	} else {
		panic(fmt.Sprintf("invalid DB type: %s", dbType))
	}
}

func GetUserObjectByUsernameInMySqlDb(username string) (models.UserObject, error) {
	// Query to get user info by username
	query := `SELECT USER_ID, USERNAME, PASSWORD_HASH, BALANCE FROM USERS WHERE USERNAME = ?`

	// Struct to store the user info
	var user models.UserObject

	// Execute the query and scan the result into the user struct
	err := mysqldb.QueryRow(query, username).Scan(&user.UserId, &user.Username, &user.UserHashedPassword, &user.UserBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no user found, return a custom error
			return models.UserObject{}, fmt.Errorf("user with username '%s' not found", username)
		}
		return models.UserObject{}, fmt.Errorf("failed to get user by username: %v", err)
	}

	return user, nil
}

func GetUserObjectByUsernameInRedisDb(username string) (models.UserObject, error) {
	// todo
	return models.UserObject{}, nil
}

func GetUserObjectByUserId(userId int) (models.UserObject, error) {
	if dbType == "mysql" {
		return GetUserObjectByUserIdInMysqlDb(userId)
	} else if dbType == "redis" {
		return GetUserObjectByUserIdInRedisDb(userId)
	} else {
		panic(fmt.Sprintf("invalid DB type: %s", dbType))
	}
}

func GetUserObjectByUserIdInMysqlDb(userId int) (models.UserObject, error) {
	// Query to get user info by username
	query := `SELECT USER_ID, USERNAME, PASSWORD_HASH, BALANCE FROM USERS WHERE USER_ID = ?`

	// Struct to store the user info
	var user models.UserObject

	// Execute the query and scan the result into the user struct
	err := mysqldb.QueryRow(query, userId).Scan(&user.UserId, &user.Username, &user.UserHashedPassword, &user.UserBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no user found, return a custom error
			return models.UserObject{}, fmt.Errorf("user with username '%s' not found", userId)
		}
		return models.UserObject{}, fmt.Errorf("failed to get user by username: %v", err)
	}

	return user, nil
}

func GetUserObjectByUserIdInRedisDb(userId int) (models.UserObject, error) {
	// todo
	return models.UserObject{}, nil
}

func GetUserBalanceFromRedisDb(userId int) (float64, error) {

	ctx := context.Background()
	userBalanceKey := "user_balance:" + strconv.Itoa(userId)

	// Get the user's balance from Redis
	balanceStr, err := redisDb.Get(ctx, userBalanceKey).Result()
	if err == redis.Nil {
		// The key does not exist
		return 0, fmt.Errorf("user not found")
	} else if err != nil {
		return 0, err // Return the error for further handling
	}

	// Convert the balance from string to integer
	balance, err := strconv.ParseFloat(balanceStr, 64)
	if err != nil {
		return 0, err // Return the error if conversion fails
	}

	return balance, nil // Return the balance if successful
}

func insertUserBalanceInRedisDb(username string, balance float64) error {
	ctx := context.Background()

	userObject, err2 := GetUserObjectByUsername(username)
	if err2 != nil {
		return err2
	}

	userBalanceKey := "user_balance:" + strconv.Itoa(userObject.UserId)

	// Convert balance to string for Redis storage
	balanceStr := strconv.FormatFloat(balance, 'f', -1, 64) // Convert to string

	// Set the user's balance in Redis
	err := redisDb.Set(ctx, userBalanceKey, balanceStr, 0).Err() // 0 means no expiration
	if err != nil {
		return err // Return the error for further handling
	}

	return nil // Return nil if successful
}

func insertUserBalanceByUserIdInRedisDb(userId int, balance float64) error {
	ctx := context.Background()

	userBalanceKey := "user_balance:" + strconv.Itoa(userId)

	// Convert balance to string for Redis storage
	balanceStr := strconv.FormatFloat(balance, 'f', -1, 64) // Convert to string

	// Set the user's balance in Redis
	err := redisDb.Set(ctx, userBalanceKey, balanceStr, 0).Err() // 0 means no expiration
	if err != nil {
		return err // Return the error for further handling
	}

	return nil // Return nil if successful
}

func TransferMoneyToUser(fromUser int, toUsername string, amount float64) error {

	if amount < 0 {
		return fmt.Errorf("amount must be positiv")
	}

	fromUserBalance, err := GetUserBalanceFromRedisDb(fromUser)
	if err != nil {
		return err
	}

	if fromUserBalance < amount {
		return fmt.Errorf("not enough balance")
	}

	toUserObject, err := GetUserObjectByUsername(toUsername)
	if err != nil {
		return err
	}

	toUserBalance, err := GetUserBalanceFromRedisDb(toUserObject.UserId)
	if err != nil {
		return err
	}

	toUserBalance += amount
	fromUserBalance -= amount

	err = insertUserBalanceByUserIdInRedisDb(fromUser, fromUserBalance)
	if err != nil {
		return err
	}

	err = insertUserBalanceByUserIdInRedisDb(toUserObject.UserId, toUserBalance)
	if err != nil {
		return err
	}

	return err

}
