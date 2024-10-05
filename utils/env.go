package utils

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type AppData struct {
	CertPath     string
	KeyPath      string
	H2ListenAddr string
	H3ListenAddr string
	H1ListenAddr string
	KeyLogFile   string
	MySqlHost    string
	MySqlPort    int
	MySqlUser    string
	MySqlPass    string
	MySqlDbName  string
	DbType       string
}

func LoadEnvFile(fileName string) AppData {
	// Load .env file
	err := godotenv.Load(fileName)

	if os.IsNotExist(err) {
		CreateEnvFile(fileName)
		os.Exit(0)
	}

	if err != nil {
		panic(err)
	}

	CertPath := os.Getenv("CertPath")
	KeyPath := os.Getenv("KeyPath")
	H1ListenAddr := os.Getenv("H1ListenAddr")
	H2ListenAddr := os.Getenv("H2ListenAddr")
	H3ListenAddr := os.Getenv("H3ListenAddr")
	KeyLogFile := os.Getenv("KeyLogFile")
	MySqlHost := os.Getenv("MySqlHost")
	MySqlPort := os.Getenv("MySqlPort")
	MySqlUser := os.Getenv("MySqlUser")
	MySqlPass := os.Getenv("MySqlPass")
	MySqlDbName := os.Getenv("MySqlDbName")
	DbType := os.Getenv("DbType")

	MySqlPortInt, err3 := strconv.Atoi(MySqlPort)
	if err3 != nil {
		fmt.Println("MySqlPort is not integer")
		panic(err3)
	}

	return AppData{
		CertPath:     CertPath,
		KeyPath:      KeyPath,
		H1ListenAddr: H1ListenAddr,
		H2ListenAddr: H2ListenAddr,
		H3ListenAddr: H3ListenAddr,
		KeyLogFile:   KeyLogFile,
		MySqlHost:    MySqlHost,
		MySqlPort:    MySqlPortInt,
		MySqlUser:    MySqlUser,
		MySqlPass:    MySqlPass,
		MySqlDbName:  MySqlDbName,
		DbType:       DbType,
	}

}

func CreateEnvFile(fileName string) error {
	// Open .env file for writing
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write default environment variables to .env file
	envContent := "CertPath=\nKeyPath=\nH2ListenAddr=127.0.0.1:443\nH3ListenAddr=127.0.0.1:443\nKeyLogFile=h3_quic.log\nMySqlHost=127.0.0.1\nMySqlPort=3306\nMySqlUser=\nMySqlPass=\nMySqlDbName=\nDbType=mysql\n"
	_, err = file.WriteString(envContent)
	if err != nil {
		return err
	}

	fmt.Println(".env file created successfully\nPlease fill it!")
	return nil
}
