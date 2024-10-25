package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nxenon/rc-h3-webapp/utils"
	"github.com/redis/go-redis/v9"
)

var mysqldb *sql.DB
var redisDb *redis.Client
var dbType string

func ConnectToMySqlDatabase(appData utils.AppData) {
	dbType = "mysql"
	// MySQL database connection parameters
	dbUser := appData.MySqlUser
	dbPass := appData.MySqlPass
	dbName := appData.MySqlDbName
	dbHost := appData.MySqlHost
	dbPort := appData.MySqlPort

	// Establish MySQL database connection
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUser, dbPass, dbHost, dbPort, dbName))
	if err != nil {
		panic(err)
	}
	mysqldb = db
	//defer db.Close()

}

func DoChecks(appData utils.AppData) {
	MakeCartProductsTableEmpty()
	MakeCartsTableEmpty()
	InsertDefaultUsers()
	InsertDefaultProducts()
	InsertDefaultCouponCodes()
}

func InsertDefaultUsers() {
	// password are same as usernames in MD5
	insertTestUserInDb("xenon", "43c1cb1c1cf84a689b551d8dd1b13190", 1)
	insertTestUserInDb("xenon2", "c122ecc4ab0a41f1418689b1c444de69", 2)
	insertTestUserInDb("xenon3", "1821cb9b461cb375cfeb9b6c90d8b4ec", 3)
	insertTestUserInDb("xenon4", "b5e3644d97062173a22b2f94de8b9a69", 4)
	insertTestUserInDb("xenon5", "f6fffba39d816899916697446a92fbf5", 5)
	insertTestUserInDb("xenon20", "c51240e40faf5f39e93d338cd83acf96", 20)
	insertTestUserInDb("xenon35", "d31cb3d338997f9a677dcc4018848929", 35)
	insertTestUserInDb("xenon50", "eacd4bbe82f0b37529fb62b9b5e4ed89", 50)
	insertTestUserInDb("xenon65", "7d766bfdafda7afbbc4b155793e992c4", 65)
	insertTestUserInDb("xenon80", "af48ddcfa19d8d27d086b0edbe3ace4f", 80)
}

func InsertDefaultCouponCodes() {
	// Insert Coupon Codes
	InsertOrUpdateCouponCode(true, 10, "Coupon1")
	InsertOrUpdateCouponCode(true, 10, "Coupon2")
	InsertOrUpdateCouponCode(true, 10, "Coupon3")
	InsertOrUpdateCouponCode(true, 10, "Coupon4")
	InsertOrUpdateCouponCode(true, 10, "Coupon5")
	InsertOrUpdateCouponCode(true, 10, "Coupon6")
	InsertOrUpdateCouponCode(true, 10, "Coupon7")
	InsertOrUpdateCouponCode(true, 10, "Coupon8")
	InsertOrUpdateCouponCode(true, 10, "Coupon9")
	InsertOrUpdateCouponCode(true, 10, "Coupon10")
}

func InsertDefaultProducts() {
	// Insert Products
	addProductIfNotExists("Laptop 1", 200, "https://images.unsplash.com/photo-1516321497487-e288fb19713f")
	addProductIfNotExists("Smartphone 1", 149, "https://images.unsplash.com/photo-1511707171634-5f897ff02aa9")
	addProductIfNotExists("Gaming Laptop 1", 470, "https://images.unsplash.com/photo-1517336714731-489689fd1ca8")
	addProductIfNotExists("Laptop 2", 210, "https://images.unsplash.com/photo-1593642532744-d377ab507dc8")
	addProductIfNotExists("Book", 20, "https://images.unsplash.com/photo-1512820790803-83ca734da794")
}

func ConnectToRedisDatabase(appData utils.AppData) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     appData.RedisAddr,
		Password: appData.RedisDbPass,
		DB:       appData.RedisDbID,
	})

	ctx := context.Background()

	// Ping the Redis server
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Could not connect to Redis:", err)
		panic(err)
	}
	fmt.Println("Connected to Redis:", pong)

	redisDb = rdb

	DoChecks(appData)
}
