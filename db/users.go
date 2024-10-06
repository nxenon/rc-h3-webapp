package db

import (
	"database/sql"
	"fmt"
	"github.com/nxenon/rc-h3-webapp/models"
)

func insertTestUserInMySqlDb(username, passwordHash string, balance int) error {
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

	return nil
}

func GetUserObjectByUsername(username string) (models.UserObject, error) {
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

func GetUserObjectByUserId(userId int) (models.UserObject, error) {
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
