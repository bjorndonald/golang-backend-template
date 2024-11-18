package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

func Connect(config *Config) {
	var (
		err error
		// port, _ = strconv.ParseUint(config.Port, 10, 32)
		dsn = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", config.User, config.Password, config.Host, config.Port, config.DBName, config.SSLMode)
	)

	log.Println(dsn)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	DB.Logger.LogMode(logger.Silent)

	if err != nil {
		fmt.Println(
			err.Error(),
		)
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")
}

var DB *gorm.DB