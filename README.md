# HTTP/2 & HTTP/3 Race Condition Lab (Web Application)
HTTP/1/2/3 Web Application for testing Last Frame Synchronization (also known as Single Packet Attack) on HTTP/2 and HTTP/3.
For see the exploitation tool and reference go to [H3SpaceX Library](https://github.com/nxenon/h3spacex).

Features:
- Apply Coupon Code on User Cart
- Transfer User Balance to another User
- Login Rate Limit Per User

## Run the App

```shell
go run ./webapp.go
```

## .env File
You have to fill .env file before run the app. The app creates .env file when it gets run at first time.
```env
CertPath=CERT_PATH
KeyPath=KEY_PATH
H2ListenAddr=127.0.0.1:443
H3ListenAddr=127.0.0.1:443
KeyLogFile=h3_quic.log
MySqlHost=IP_ADDRESS
MySqlPort=3306
MySqlUser=USERNAME
MySqlPass=PASSWORD
MySqlDbName=DB_NAME
DbType=mysql
RedisAddr=127.0.0.1:6379
RedisDbPass=
RedisDbID=0
MaxIncomingStreams=100
```

## Login
Username & Passwords are same! e.g `xenon2:xenon2`. The third is balance of each user. See [db.go](./db/db.go).
```go
insertTestUserInMySqlDb("xenon", "43c1cb1c1cf84a689b551d8dd1b13190", 1)
insertTestUserInMySqlDb("xenon2", "c122ecc4ab0a41f1418689b1c444de69", 2)
insertTestUserInMySqlDb("xenon3", "1821cb9b461cb375cfeb9b6c90d8b4ec", 3)
insertTestUserInMySqlDb("xenon4", "b5e3644d97062173a22b2f94de8b9a69", 4)
insertTestUserInMySqlDb("xenon5", "f6fffba39d816899916697446a92fbf5", 5)
insertTestUserInMySqlDb("xenon20", "c51240e40faf5f39e93d338cd83acf96", 20)
insertTestUserInMySqlDb("xenon35", "d31cb3d338997f9a677dcc4018848929", 35)
insertTestUserInMySqlDb("xenon50", "eacd4bbe82f0b37529fb62b9b5e4ed89", 50)
insertTestUserInMySqlDb("xenon65", "7d766bfdafda7afbbc4b155793e992c4", 65)
insertTestUserInMySqlDb("xenon80", "af48ddcfa19d8d27d086b0edbe3ace4f", 80)
```

## Coupon Codes
Coupon Codes: e.g `Coupon1` with `10%` discount for cart. See [db.go](./db/db.go).
```go
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
```

