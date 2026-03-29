package config

import (
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitDB() *gorm.DB {
	dbDir := "./db"

	_, err := os.Stat(dbDir)
	if os.IsNotExist(err) {
		if err := os.Mkdir(dbDir, 0755); err != nil {
			panic("build db failed: " + err.Error())
		}
	}

	dsn := filepath.Join(dbDir, "blog.db")

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		// Logger configuration
		Logger: logger.Default.LogMode(logger.Info),

		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",    // Prefix for all table names (e.g., "app_")
			SingularTable: false, // Use singular table names (User -> user instead of users)
			NoLowerCase:   false, // Disable automatic lowercasing
			NameReplacer:  nil,   // Custom name replacer function
		},
	})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
