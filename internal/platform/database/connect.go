package database

import (
	"fmt"

	"github.com/gusram01/linked-bookmarks/internal/platform/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Initialize(models ...interface{}) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=%s" ,
        config.Config("GC_MARK_DB_HOST"),
		config.Config("GC_MARK_DB_USER"),
		config.Config("GC_MARK_DB_PASS"),
		config.Config("GC_MARK_DB_NAME"),
		config.Config("GC_MARK_DB_SSL_MODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("DB::init::failed::%s", err.Error()))
	}


	fmt.Println("DB connection initialized")

    if models != nil {
        db.AutoMigrate(models...)
        fmt.Println("DB migrated")
    }

    DB = db
}
