package dal

import (
	config "api/src/config"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// Initialize ...Function to initialize to DB
func Initialize() map[string]*gorm.DB {
	fmt.Println("Connecting to db-----------------")
	configuration := config.GetConfig()
	fmt.Println("configuration : ", configuration)
	dbConnectionMap := make(map[string]*gorm.DB)
	for k, v := range configuration.DbConnections {
		db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", v.DBUser, v.DBPassword, v.DBHost, v.DBName))
		if err != nil {
			panic(err.Error())
		}

		// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
		db.DB().SetMaxIdleConns(10)

		// SetMaxOpenConns sets the maximum number of open connections to the database.
		db.DB().SetMaxOpenConns(100)

		// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
		db.DB().SetConnMaxLifetime(time.Hour)
		db.LogMode(true)
		dbConnectionMap[k] = db
	}
	return dbConnectionMap
}
