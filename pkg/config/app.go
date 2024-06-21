package config

import(
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"

	//"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var (
	db *gorm.DB
)


func Connect() {
	err := godotenv.Load()
  	if err != nil {
    	log.Fatal("Error loading .env file")
  	}
	dsn := fmt.Sprintf("root:%s@tcp(localhost:3306)/productbot?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("DB_PASSWORD"))
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db = d
}

func GetDB() *gorm.DB {
	return db
}