package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/nxenon/rc-h3-webapp/models"
	"github.com/nxenon/rc-h3-webapp/utils"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
)

func GetCartIdByUserId(userId int) (int, error) {
	if dbType == "mysql" {
		return GetCartIdByUserIdInMySqlDb(userId)
	} else if dbType == "redis" {
		return GetCartIdByUserIdInRedisDb(userId)
	} else {
		panic(fmt.Sprintf("invalid DB type: %s", dbType))
	}
}

func GetCartIdByUserIdInMySqlDb(userId int) (int, error) {
	var cartId int

	// Query to get the cart ID for the specified user ID
	err := mysqldb.QueryRow("SELECT CART_ID FROM CARTS WHERE USER_ID = ?", userId).Scan(&cartId)

	if err == sql.ErrNoRows {
		// No cart found for the user, returning 0 and nil error
		return 0, nil
	} else if err != nil {
		log.Printf("Error retrieving cart ID: %v", err)
		return 0, err // Return error if something went wrong
	}

	return cartId, nil // Return the found cart ID and nil error
}

func GetCartIdByUserIdInRedisDb(userId int) (int, error) {
	return 0, nil
}

func GetCartByUserId(userId int) (models.Cart, error) {
	if dbType == "mysql" {
		return GetCartByUserIdInRedisDb(userId)
	} else {
		panic(fmt.Sprintf("invalid DB type: %s", dbType))
	}
}

//func GetCartByIdInMySqlDb(cartId int) (models.Cart, error) {
//	var cart models.Cart
//
//	// Get the cart details
//	err := mysqldb.QueryRow("SELECT CART_ID, USER_ID, CART_OVERAL_PRICE FROM CARTS WHERE CART_ID = ?", cartId).
//		Scan(&cart.CartId, &cart.UserId, &cart.CartOverallPrice)
//
//	if err == sql.ErrNoRows {
//		return cart, nil // No cart found, return an empty cart
//	} else if err != nil {
//		log.Printf("Error retrieving cart: %v", err)
//		return cart, err // Return error if something went wrong
//	}
//
//	// Retrieve products in the cart
//	rows, err := mysqldb.Query("SELECT p.PRODUCT_ID, p.PRODUCT_NAME, p.PRODUCT_PRICE, p.IMAGE_URL, cp.PRODUCT_IN_CART_ID FROM CART_PRODUCTS cp JOIN PRODUCTS p ON cp.PRODUCT_ID = p.PRODUCT_ID WHERE cp.CART_ID = ?", cart.CartId)
//	if err != nil {
//		log.Printf("Error retrieving products: %v", err)
//		return cart, err // Return error if something went wrong
//	}
//	defer rows.Close()
//
//	// Populate the products slice
//	for rows.Next() {
//		var product models.ProductObject
//		err := rows.Scan(&product.ProductId, &product.ProductName, &product.ProductPrice, &product.ProductImageUrl, &product.ProductInCartId)
//		if err != nil {
//			log.Printf("Error scanning product: %v", err)
//			return cart, err // Return error if something went wrong
//		}
//		cart.Products = append(cart.Products, product)
//	}
//
//	// Check for errors from iterating over rows
//	if err = rows.Err(); err != nil {
//		log.Printf("Error iterating over products: %v", err)
//		return cart, err // Return error if something went wrong
//	}
//
//	return cart, nil // Return the cart object and nil error
//}

// todo
func GetCartByUserIdInRedisDb(userId int) (models.Cart, error) {
	ctx := context.Background()
	var cart models.Cart
	cart.UserId = userId

	userCartKey := "user_cart:" + strconv.Itoa(userId)

	cartInDb, err := redisDb.HGetAll(ctx, userCartKey).Result()
	if err == redis.Nil {
		err2 := redisDb.HSet(
			ctx,
			userCartKey,
			utils.StructToMap(cart),
		).Err()

		if err2 != nil {
			panic(err2)
		}
		return cart, nil

	} else if err != nil {
		panic(err)
	}
	cop, _ := strconv.Atoi(cartInDb["CartOverallPrice"])
	cart.CartOverallPrice = cop
	err = json.Unmarshal([]byte(cartInDb["Products"]), &cart.Products)
	if err != nil {
		cart.Products = []models.ProductObject{}
	}

	return cart, nil

}

func UpdateCartByUserId(userId int, cart models.Cart) error {
	if dbType == "mysql" {
		return UpdateCartByUserIdInRedisDb(userId, cart)
	} else {
		panic(fmt.Sprintf("invalid DB type: %s", dbType))
	}
}

func UpdateCartByIdInMysqlDb(cartId int, newOverallPrice int) error {
	// Update query to update the cart's overall price
	updateQuery := `UPDATE CARTS SET CART_OVERAL_PRICE = ? WHERE CART_ID = ?`

	// Execute the update query
	_, err := mysqldb.Exec(updateQuery, newOverallPrice, cartId)
	if err != nil {
		return fmt.Errorf("failed to update cart: %v", err)
	}

	return nil
}

func UpdateCartByUserIdInRedisDb(userId int, cart models.Cart) error {
	ctx := context.Background()
	userCartKey := "user_cart:" + strconv.Itoa(userId)

	productString, err := json.Marshal(cart.Products)
	if err != nil {
		return err
	}

	x := map[string]interface{}{
		"CartId":           0,
		"UserId":           userId,
		"CartOverallPrice": cart.CartOverallPrice,
		"Products":         productString,
	}

	err3 := redisDb.HSet(ctx, userCartKey, x).Err()
	return err3
}

func MakeCartProductsTableEmpty() {
	if dbType == "mysql" {
		MakeCartProductsTableEmptyInMysqlDb()
	} else if dbType == "redis" {
		MakeCartProductsTableEmptyInRedisDb()
	} else {
		panic(fmt.Sprintf("invalid DB type: %s", dbType))
	}
}

func MakeCartProductsTableEmptyInMysqlDb() {
	query := "DELETE FROM CART_PRODUCTS"
	_, err := mysqldb.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully cleared all products in carts")
}

func MakeCartProductsTableEmptyInRedisDb() {
	// todo
}

func MakeCartsTableEmpty() {
	if dbType == "mysql" {
		MakeCartsTableEmptyInMysqlDb()
		MakeCartsTableEmptyInRedisDb()
	} else {
		panic(fmt.Sprintf("invalid DB type: %s", dbType))
	}
}

func MakeCartsTableEmptyInMysqlDb() {
	query := "DELETE FROM CARTS"

	// Assuming `mysqldb` is your database connection object
	_, err := mysqldb.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully cleared all carts")
}

func MakeCartsTableEmptyInRedisDb() {
	ctx := context.Background()
	iter := redisDb.Scan(ctx, 0, "user_cart:*", 0).Iterator()
	for iter.Next(ctx) {
		// Delete each matching key
		if err := redisDb.Del(ctx, iter.Val()).Err(); err != nil {
			fmt.Printf("error deleting key %s: %v\n", iter.Val(), err)
		}
	}

	// Check for any errors during iteration
	if err := iter.Err(); err != nil {
		fmt.Printf("error during scan iteration: %v\n", err)
	}

	fmt.Println("Cleared Redis Products in User Carts.")

}
