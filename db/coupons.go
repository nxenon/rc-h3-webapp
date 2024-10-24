package db

import (
	"database/sql"
	"fmt"
	"github.com/nxenon/rc-h3-webapp/models"
)

func InsertOrUpdateCouponCode(isValid bool, discountPercent int, couponValue string) error {

	if dbType == "mysql" {
		return InsertOrUpdateCouponCodeInMySqlDb(isValid, discountPercent, couponValue)
	} else if dbType == "redis" {
		return InsertOrUpdateCouponCodeInRedisDb(isValid, discountPercent, couponValue)
	} else {
		panic(fmt.Sprintf("invalid DB type: %s", dbType))
	}

}

func InsertOrUpdateCouponCodeInRedisDb(isValid bool, discountPercent int, couponValue string) error {
	// todo
	return nil
}

func InsertOrUpdateCouponCodeInMySqlDb(isValid bool, discountPercent int, couponValue string) error {
	// Check if the coupon already exists
	var couponId int
	query := `SELECT COUPON_ID FROM COUPONS WHERE COUPON_VALUE = ?`
	err := mysqldb.QueryRow(query, couponValue).Scan(&couponId)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check coupon existence: %v", err)
	}

	if err == sql.ErrNoRows {
		// Insert new coupon if it doesn't exist
		insertQuery := `INSERT INTO COUPONS (IS_VALID, DISCOUNT_PERCENT, COUPON_VALUE) VALUES (?, ?, ?)`
		_, err = mysqldb.Exec(insertQuery, isValid, discountPercent, couponValue)
		if err != nil {
			return fmt.Errorf("failed to insert new coupon: %v", err)
		}
		fmt.Println("Inserted new coupon successfully")
	} else {
		// Update existing coupon if it exists
		updateQuery := `UPDATE COUPONS SET IS_VALID = ?, DISCOUNT_PERCENT = ? WHERE COUPON_ID = ?`
		_, err = mysqldb.Exec(updateQuery, isValid, discountPercent, couponId)
		if err != nil {
			return fmt.Errorf("failed to update coupon: %v", err)
		}
		fmt.Println("Updated existing coupon successfully")
	}

	return nil
}

func GetCouponByValue(couponValue string) (models.CouponCode, error) {
	if dbType == "mysql" {
		return GetCouponByValueInMySqlDb(couponValue)
	} else if dbType == "redis" {
		return GetCouponByValueInRedisDb(couponValue)
	} else {
		panic(fmt.Sprintf("Invalid Db Type: %s", dbType))
	}
}

func GetCouponByValueInMySqlDb(couponValue string) (models.CouponCode, error) {
	var coupon models.CouponCode

	query := `SELECT COUPON_ID, IS_VALID, DISCOUNT_PERCENT, COUPON_VALUE FROM COUPONS WHERE COUPON_VALUE = ?`
	err := mysqldb.QueryRow(query, couponValue).Scan(&coupon.CouponId, &coupon.IsValid, &coupon.DiscountPercent, &coupon.CouponValue)

	if err != nil {
		if err == sql.ErrNoRows {
			return coupon, fmt.Errorf("coupon not found")
		}
		return coupon, err
	}

	return coupon, nil
}

func GetCouponByValueInRedisDb(couponValue string) (models.CouponCode, error) {
	return models.CouponCode{}, nil
}
