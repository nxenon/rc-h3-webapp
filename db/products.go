package db

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/nxenon/rc-h3-webapp/models"
	"log"
)

// Check if a product exists and insert if it doesn't
func addProductIfNotExists(productName string, productPrice int, imageURL string) error {
	if dbType == "mysql" {
		return addProductInMySqlDbIfNotExists(productName, productPrice, imageURL)
	} else if dbType == "redis" {
		return addProductInRedisDbIfNotExists(productName, productPrice, imageURL)
	} else {
		panic(fmt.Sprintf("invalid DB type: %s", dbType))
	}
}

func addProductInMySqlDbIfNotExists(productName string, productPrice int, imageURL string) error {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM PRODUCTS WHERE PRODUCT_NAME = ?)`

	// Check if the product exists
	err := mysqldb.QueryRow(query, productName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if product exists: %v", err)
	}

	if exists {
		// If product exists, update the price
		updateQuery := `UPDATE PRODUCTS SET PRODUCT_PRICE = ?, IMAGE_URL = ? WHERE PRODUCT_NAME = ?`
		_, err = mysqldb.Exec(updateQuery, productPrice, imageURL, productName)
		if err != nil {
			return fmt.Errorf("failed to update product price: %v", err)
		}
		fmt.Println("Product price updated:", productName)
	} else {
		// If product does not exist, insert a new product
		insertQuery := `INSERT INTO PRODUCTS (PRODUCT_NAME, PRODUCT_PRICE, IMAGE_URL) VALUES (?, ?, ?)`
		_, err = mysqldb.Exec(insertQuery, productName, productPrice, imageURL)
		if err != nil {
			return fmt.Errorf("failed to insert product: %v", err)
		}
		fmt.Println("Product inserted:", productName)
	}

	return nil
}

func addProductInRedisDbIfNotExists(productName string, productPrice int, imageURL string) error {
	// todo
	return nil
}

func GetProduct(productId int) (models.ProductObject, error) {
	if dbType == "mysql" {
		return GetProductInMySqlDb(productId)
	} else if dbType == "redis" {
		return GetProductInRedisDb(productId)
	} else {
		panic(fmt.Sprintf("Invalid Db Type: %s", dbType))
	}
}

func GetProductInMySqlDb(productId int) (models.ProductObject, error) {
	// Query to get user info by username
	query := `SELECT PRODUCT_ID, PRODUCT_NAME, PRODUCT_PRICE, IMAGE_URL FROM PRODUCTS WHERE PRODUCT_ID = ?`

	// Struct to store the user info
	var product models.ProductObject

	// Execute the query and scan the result into the user struct
	err := mysqldb.QueryRow(query, productId).Scan(&product.ProductId, &product.ProductName, &product.ProductPrice, &product.ProductImageUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no user found, return a custom error
			return models.ProductObject{}, fmt.Errorf("product ID '%s' not found", productId)
		}
		return models.ProductObject{}, fmt.Errorf("failed to get product: %v", err)
	}

	return product, nil
}

func GetProductInRedisDb(productId int) (models.ProductObject, error) {
	// Query to get user info by username
	// todo
	return models.ProductObject{}, nil
}

func GetAllProducts() ([]models.ProductObject, error) {
	if dbType == "mysql" {
		return GetAllProductsInMySqlDb()
	} else if dbType == "redis" {
		return GetAllProductsInRedisDb()
	} else {
		panic(fmt.Sprintf("Invalid Db Type: %s", dbType))
	}
}

func GetAllProductsInMySqlDb() ([]models.ProductObject, error) {
	// Query to get all products
	query := `SELECT PRODUCT_ID, PRODUCT_NAME, PRODUCT_PRICE, IMAGE_URL FROM PRODUCTS`

	// Prepare a slice to hold the product objects
	var products []models.ProductObject

	// Execute the query
	rows, err := mysqldb.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close() // Ensure rows are closed after we finish with them

	// Loop through the results and scan into product objects
	for rows.Next() {
		var product models.ProductObject
		err := rows.Scan(&product.ProductId, &product.ProductName, &product.ProductPrice, &product.ProductImageUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		products = append(products, product) // Add the product to the slice
	}

	// Check for errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error encountered while iterating over rows: %v", err)
	}

	// Return the slice of products
	return products, nil
}

func GetAllProductsInRedisDb() ([]models.ProductObject, error) {
	// todo
	return []models.ProductObject{}, nil
}

func AddProductToUserCart(userId, productId int) error {

	if dbType == "mysql" {
		return AddProductToUserCartInRedisDb(userId, productId)
	} else {
		panic(fmt.Sprintf("Invalid Db Type: %s", dbType))
	}

}

func AddProductToUserCartInMySqlDb(userId, productId int) error {

	// Check if user has an existing cart
	var cart models.Cart
	var err error
	err = mysqldb.QueryRow("SELECT CART_ID, CART_OVERAL_PRICE FROM CARTS WHERE USER_ID = ?", userId).Scan(&cart.CartId, &cart.CartOverallPrice)
	productObject, err2 := GetProduct(productId)
	if err2 != nil {
		return err
	}

	if err == sql.ErrNoRows {
		// If no cart exists, create a new cart for the user
		result, err := mysqldb.Exec("INSERT INTO CARTS (USER_ID, CART_OVERAL_PRICE) VALUES (?, ?)", userId, 0)
		if err != nil {
			return err
		}
		// Get the last inserted cart ID
		cartId, err := result.LastInsertId()

		if err != nil {
			return err
		}
		cart.CartId = int(cartId)
		cart.UserId = userId
		cart.CartOverallPrice = 0
	} else if err != nil {
		return err
	}

	// Update the cart's overall price
	newPrice := cart.CartOverallPrice + productObject.ProductPrice
	_, err = mysqldb.Exec("UPDATE CARTS SET CART_OVERAL_PRICE = ? WHERE CART_ID = ?", newPrice, cart.CartId)
	if err != nil {
		return err
	}

	// Add product to the CART_PRODUCTS table
	_, err = mysqldb.Exec("INSERT INTO CART_PRODUCTS (PRODUCT_ID, CART_ID) VALUES (?, ?)", productId, cart.CartId)
	if err != nil {
		return err
	}

	return nil
}

func AddProductToUserCartInRedisDb(userId, productId int) error {

	userCart, err := GetCartByUserIdInRedisDb(userId)
	if err != nil {
		return err
	}

	p, err2 := GetProductInMySqlDb(productId)
	if err2 != nil {
		return err2
	}

	p.ProductInCartUUID = uuid.New().String()

	userCart.Products = append(userCart.Products, p)
	userCart.CartOverallPrice = userCart.CartOverallPrice + p.ProductPrice

	return UpdateCartByUserId(userId, userCart)

}

//
//func RemoveProductFromCartByPRODUCT_IN_CART_ID(productInCartId int, cartId int) error {
//	if dbType == "mysql" {
//		return RemoveProductFromCartByPRODUCT_IN_CART_IDInMySqlDb(productInCartId, cartId)
//	} else if dbType == "redis" {
//		return RemoveProductFromCartByPRODUCT_IN_CART_IDInRedisDb(productInCartId, cartId)
//	} else {
//		panic(fmt.Sprintf("Invalid Db Type: %s", dbType))
//	}
//}

func RemoveProductFromCartByUserId(productInCartUUID string, userId int) error {
	if dbType == "mysql" {
		return RemoveProductFromCartByUserIdInRedisDb(productInCartUUID, userId)
	} else {
		panic(fmt.Sprintf("Invalid Db Type: %s", dbType))
	}
}

func RemoveProductFromCartByUserIdInRedisDb(productInCartUUID string, userId int) error {

	userCart, err := GetCartByUserId(userId)
	if err != nil {
		return err
	}

	var tempProducts []models.ProductObject
	var productExists bool
	productExists = false
	for _, p := range userCart.Products {
		if p.ProductInCartUUID != productInCartUUID {
			tempProducts = append(tempProducts, p)
		} else {
			productExists = true
			userCart.CartOverallPrice -= p.ProductPrice
		}
	}

	userCart.Products = tempProducts
	err2 := UpdateCartByUserIdInRedisDb(userId, userCart)
	if err2 == nil {
		if !productExists {
			return fmt.Errorf("product does not exist")
		}
	}
	return err2

}

func RemoveProductFromCartByPRODUCT_IN_CART_IDInMySqlDb(productInCartId int, cartId int) error {
	query := "DELETE FROM CART_PRODUCTS WHERE PRODUCT_IN_CART_ID = ? AND CART_ID = ?"
	// Execute the query
	result, err := mysqldb.Exec(query, productInCartId, cartId)
	if err != nil {
		log.Printf("Error removing product from cart: %v", err)
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error checking affected rows: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func RemoveProductFromCartByPRODUCT_IN_CART_IDInRedisDb(productInCartId int, cartId int) error {
	// todo
	return nil
}

func GetProductFromPRODUCT_IN_CART_ID(productInCartId int) (models.ProductObject, error) {
	if dbType == "mysql" {
		return GetProductFromPRODUCT_IN_CART_IDInMySqlDb(productInCartId)
	} else if dbType == "redis" {
		return GetProductFromPRODUCT_IN_CART_IDInRedisDb(productInCartId)
	} else {
		panic(fmt.Sprintf("Invalid Db Type: %s", dbType))
	}
}

func GetProductFromPRODUCT_IN_CART_IDInMySqlDb(productInCartId int) (models.ProductObject, error) {
	var product models.ProductObject
	err := mysqldb.QueryRow(`
		SELECT 
			p.PRODUCT_ID, 
			p.PRODUCT_NAME, 
			p.PRODUCT_PRICE, 
			p.IMAGE_URL 
		FROM 
			CART_PRODUCTS cp 
		JOIN 
			PRODUCTS p 
		ON 
			cp.PRODUCT_ID = p.PRODUCT_ID 
		WHERE 
			cp.PRODUCT_IN_CART_ID = ?`, productInCartId).Scan(&product.ProductId, &product.ProductName, &product.ProductPrice, &product.ProductImageUrl)

	if err != nil {
		return product, err
	}

	// Set the `ProductInCartId` field to the one used for querying
	product.ProductInCartId = productInCartId

	return product, nil
}

func GetProductFromPRODUCT_IN_CART_IDInRedisDb(productInCartId int) (models.ProductObject, error) {
	// todo
	return models.ProductObject{}, nil
}
