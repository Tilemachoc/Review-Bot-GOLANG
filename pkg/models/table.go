package models

import (
	"fmt"
	"log"
	"time"

	"github.com/Tilemachoc/TASK1/pkg/config"
	"github.com/jinzhu/gorm"
	// "gorm.io/gorm"
)

var db *gorm.DB

type info struct {
	gorm.Model
	// Id			string	`json:"id"`
	Name		string	`json:"name"`
	Email		string	`json:"email"`
	Region		string	`json:"region"`
}


type interaction struct {
	gorm.Model
	// Id			string		`json:"id"`
	Name			string		`json:"name"`
	History 		[]message	`json:"history"`
}

type message struct {
	Sender   	string    	`json:"sender"`
	Receiver 	string    	`json:"receiver"`
	Content  	string    	`json:"content"`
	Timestamp 	time.Time 	`json:"timestamp"`
}


type review struct {
	gorm.Model
	// Id			string	`json:"id"`
	Name		string	`json:"name"`
	Item		string	`json:"item"`
	Review		uint8	`json:"review"`
}


func init() {
	config.Connect()
	db = config.GetDB()
	err := db.AutoMigrate(&info{}, &interaction{}, &review{})
	if err != nil {
		fmt.Println("Error migrating models:", err)
		panic(err)
	}

	var count int32
	db.Model(&info{}).Count(&count)
	if count == 0 {
		err := insertData()
		if err != nil {
			log.Fatalf("Problem with insertData: %v", err)
		}
	}
}


func insertData() error{
	//INFO--
	infos := []info{
		{Name: "Tilemachos", Email: "tilemachos@gmail.com", Region: "EU"},
		{Name: "George", Email: "george@gmail.com", Region: "NA"},
	}

	//INTERACTIONS--
	messages := []message{
		{Sender: "CHATBOT", Receiver: "Tilemachos", Content: "Hello Tilemachos, we noticed you've recently received your iPhone 13. We'd love to hear about your experience Can you spare a few minutes to share your thoughts?", Timestamp: time.Now()},
		{Sender: "Tilemachos", Receiver: "CHATBOT", Content: "Sure, I can do that.", Timestamp: time.Now().Add(time.Minute * 5)},
	}

	interactions := []interaction{
		{Name: "Tilemachos", History: messages},
		{Name: "George", History: []message{}},
	}

	//REVIEWS--
	reviews := []review{
		{Name: "Tilemachos", Item: "iPhone 13", Review: 4},
		{Name: "Tilemachos", Item: "GPU", Review: 3},
		{Name: "George", Item: "Monitor", Review: 2},
	}

	if err := db.Create(&infos).Error; err != nil {
		return fmt.Errorf("error creating info record: %w", err)
	}

	if err := db.Create(&interactions).Error; err != nil {
		return fmt.Errorf("error creating interactions record: %w", err)
	}

	if err := db.Create(&reviews).Error; err != nil {
		return fmt.Errorf("error creating reviews record: %w", err)
	}

	return nil
}