package database

import (
	"fmt"
	"strconv"

	"github.com/gusram01/linked-bookmarks/internal/platform/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Initialize() {

	var err error

	p := config.Config("GC_MARK_PORT")
	port, parsePortErr := strconv.ParseUint(p, 10, 32)

	if parsePortErr != nil {
		panic("DB::PORT::parsing::err")
	}


	dsn := fmt.Sprintf(
		"host=db port=%d user=%s password=%s dbname=%s sslmode=disable",
		port,
		config.Config("DB_USER"),
		config.Config("DB_PASSWORD"),
		config.Config("DB_NAME"),
	)


	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("DB::init::failed::%s", err.Error()))
	}


	fmt.Println("DB connection initialized")

	DB.AutoMigrate()

	fmt.Println("DB migrated")
}
