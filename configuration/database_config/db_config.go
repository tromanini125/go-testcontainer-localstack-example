package databaseconfig

import (
	"context"
	"fmt"
	"log"

	"github.com/tromanini125/go-testcontainer-localstack-example/configuration"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func Connect(ctx context.Context) error {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		configuration.Config.DBConfig.User,
		configuration.Config.DBConfig.Password,
		configuration.Config.DBConfig.Host,
		configuration.Config.DBConfig.Port,
		configuration.Config.DBConfig.Database,
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

func GetConnection() (*gorm.DB, error) {
	if db == nil {
		err := Connect(context.Background())
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
			return nil, err
		}
	}
	return db, nil
}
