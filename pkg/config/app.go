package config

import(
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db *gorm.DB
)


func Connect() {
	dns := fmt.Sprintf("root:%s@tcp(localhoost:3306/productbot?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_PASSWORD"))
	d, err := gorm.Open("mysql", dns)
	if err != nil {
		panic(err)
	}
	db = d
}

func GetDB() *gorm.DB {
	return db
}