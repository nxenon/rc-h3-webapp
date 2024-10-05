package main

import (
	"fmt"
	"os"
	"rc-h3-webapp/apps"
	"rc-h3-webapp/db"
	"rc-h3-webapp/routes"
	"rc-h3-webapp/utils"
)

var AppData utils.AppData

func main() {
	AppData = utils.LoadEnvFile(".env")
	utils.GenerateSecretKey(32)
	routes.SetRoutes()

	if AppData.DbType == "mysql" {
		db.ConnectToMySqlDatabase(AppData)
	} else {
		fmt.Println("Invalid Database Type: ", AppData.DbType)
		fmt.Println("Valid Database types: mysql")
		os.Exit(1)
	}
	go apps.StartHttp2Server(AppData)
	apps.StartHttp3Server(AppData)

}
