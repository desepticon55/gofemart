package common

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress        string
	DatabaseConnString   string
	AccrualSystemAddress string
}

func ParseConfig() Config {
	defaultAddress := "localhost:8080"
	if envAddr, exists := os.LookupEnv("RUN_ADDRESS"); exists {
		defaultAddress = envAddr
	}
	address := flag.String("a", defaultAddress, "Server address")

	defaultDatabaseConnString := "postgres://postgres:postgres@localhost:5432/postgres"
	if envDatabaseConnString, exists := os.LookupEnv("DATABASE_URI"); exists {
		defaultDatabaseConnString = envDatabaseConnString
	}
	databaseConnString := flag.String("d", defaultDatabaseConnString, "Database connection string")

	defaultAccrualSystemAddress := "http://localhost:8081"
	if envAccrualSystemAddress, exists := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); exists {
		defaultAccrualSystemAddress = envAccrualSystemAddress
	}
	accrualSystemAddress := flag.String("r", defaultAccrualSystemAddress, "Accrual system address")

	flag.Parse()
	return Config{
		ServerAddress:        *address,
		DatabaseConnString:   *databaseConnString,
		AccrualSystemAddress: *accrualSystemAddress,
	}
}
