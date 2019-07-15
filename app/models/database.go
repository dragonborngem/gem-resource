package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var db, conn *gorm.DB
var err error

func init() {
	godotenv.Load()
	userName := os.Getenv("DATABASE_USER")
	databaseName := os.Getenv("DATABASE_NAME")
	pass := os.Getenv("DATABASE_PASS")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if os.Getenv("DATABASE_URL") != "" {
		conn, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			fmt.Println(err)
		}
	} else {
		conn, err = gorm.Open("postgres", "host="+host+" port="+port+" user="+userName+" dbname="+databaseName+" password="+pass+" sslmode=disable")
		fmt.Println("postgres", "host="+host+" port="+port+" user="+userName+" dbname="+databaseName+" password="+pass+" sslmode=disable")
		if err != nil {
			fmt.Println(err)
		}
	}
	db = conn
	db.LogMode(false)
	//defer db.Close()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(os.Getenv("DATABASE_URL"))
	db.Debug().AutoMigrate(&LoginSession{})
}

func OpenDB() *gorm.DB {
	return db
}
