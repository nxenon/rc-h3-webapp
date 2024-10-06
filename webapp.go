package main

import (
	"fmt"
	"github.com/nxenon/rc-h3-webapp/apps"
	"github.com/nxenon/rc-h3-webapp/db"
	"github.com/nxenon/rc-h3-webapp/routes"
	"github.com/nxenon/rc-h3-webapp/utils"
	"os"
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
