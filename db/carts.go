package db

import (
	"database/sql"
	"fmt"
	"github.com/nxenon/rc-h3-webapp/models"
	"log"
)

func GetCartIdByUserId(userId int) (int, error) {
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

func GetCartById(cartId int) (models.Cart, error) {
	var cart models.Cart

	// Get the cart details
	err := mysqldb.QueryRow("SELECT CART_ID, USER_ID, CART_OVERAL_PRICE FROM CARTS WHERE CART_ID = ?", cartId).
		Scan(&cart.CartId, &cart.UserId, &cart.CartOverallPrice)

	if err == sql.ErrNoRows {
		return cart, nil // No cart found, return an empty cart
	} else if err != nil {
		log.Printf("Error retrieving cart: %v", err)
		return cart, err // Return error if something went wrong
	}

	// Retrieve products in the cart
	rows, err := mysqldb.Query("SELECT p.PRODUCT_ID, p.PRODUCT_NAME, p.PRODUCT_PRICE, p.IMAGE_URL, cp.PRODUCT_IN_CART_ID FROM CART_PRODUCTS cp JOIN PRODUCTS p ON cp.PRODUCT_ID = p.PRODUCT_ID WHERE cp.CART_ID = ?", cart.CartId)
	if err != nil {
		log.Printf("Error retrieving products: %v", err)
		return cart, err // Return error if something went wrong
	}
	defer rows.Close()

	// Populate the products slice
	for rows.Next() {
		var product models.ProductObject
		err := rows.Scan(&product.ProductId, &product.ProductName, &product.ProductPrice, &product.ProductImageUrl, &product.ProductInCartId)
		if err != nil {
			log.Printf("Error scanning product: %v", err)
			return cart, err // Return error if something went wrong
		}
		cart.Products = append(cart.Products, product)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over products: %v", err)
		return cart, err // Return error if something went wrong
	}

	return cart, nil // Return the cart object and nil error
}

func UpdateCartById(cartId int, newOverallPrice int) error {
	// Update query to update the cart's overall price
	updateQuery := `UPDATE CARTS SET CART_OVERAL_PRICE = ? WHERE CART_ID = ?`

	// Execute the update query
	_, err := mysqldb.Exec(updateQuery, newOverallPrice, cartId)
	if err != nil {
		return fmt.Errorf("failed to update cart: %v", err)
	}

	return nil
}

func MakeCartProductsTableEmpty() {
	query := "DELETE FROM CART_PRODUCTS"
	_, err := mysqldb.Exec(query)
	if err != nil {
		panic(err)
	}
	log.Println("Successfully cleared all products in carts")
}

func MakeCartsTableEmpty() {
	query := "DELETE FROM CARTS"

	// Assuming `mysqldb` is your database connection object
	_, err := mysqldb.Exec(query)
	if err != nil {
		panic(err)
	}
	log.Println("Successfully cleared all carts")
}
