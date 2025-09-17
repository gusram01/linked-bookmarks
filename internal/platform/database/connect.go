package database

import (
	"fmt"

	"github.com/gusram01/linked-bookmarks/internal/platform/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Initialize(models ...interface{}) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=%s port=%s",
		config.ENVS.DbHost,
		config.ENVS.DbUser,
		config.ENVS.DbPass,
		config.ENVS.DbName,
		config.ENVS.DbSSLMode,
		config.ENVS.DbPort,
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
