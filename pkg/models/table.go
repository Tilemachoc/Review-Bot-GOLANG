package models

import (
	"fmt"
	"log"
	"time"

	"github.com/Tilemachoc/TASK1/pkg/config"
	// "github.com/jinzhu/gorm"
	"gorm.io/gorm"
)

var db *gorm.DB
// need a user variable to send it to main so we can send it to templates
// var users []User

// func GetUsers() []User {
// 	return users
// }


//Philosophy behind second table structure:
//Every User has conversations, orders and reviews
//each conversation has messages
//and each order could have many items
//and each item has a review
//we make a connection with foreignkeys and we also add constraint onupdate-ondelete
type User struct {
	UserID 			uint 			`gorm:"primaryKey"`
	Name			string			`gorm:"size:100;not null"`
	Email			string			`gorm:"size:150;not null;unique"`
	CreatedAt		time.Time		`gorm:"autoCreateTime"`

	Conversations 	[]Conversation	`gorm:"foreignKey:UserID"`
	Orders			[]Order			`gorm:"foreignKey:UserID"`
	Reviews			[]Review		`gorm:"foreignKey:UserID"`
}

// We create auto the time the conversation started and the endedat will be renewed on each message so we do that "manually" with time.Now()
type Conversation struct {
	ConversationID 		uint		`gorm:"primaryKey"`
	UserID 				uint		`gorm:"not null"`
	StartedAt 			time.Time	`gorm:"autoCreateTime"`
	EndedAt 			time.Time
	Messages 			[]Message	`gorm:"foreignKey:ConversationID"`
}

// Since message will be between user-bot sender will be either "user" or "bot"
type Message struct {
	MessageID 			uint		`gorm:"primaryKey"`
	ConversationID		uint		`gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Sender				string		`gorm:"size:10"`
	MessageText			string		`gorm:"type:text"`
	SentAt				time.Time	`gorm:"autoCreateTime"`
}

type Order struct {
	OrderID			uint		`gorm:"primaryKey"`
	UserID			uint		`gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	OrderDate		time.Time	`gorm:"autoCreateTime"`
	OrderItems []OrderItem		`gorm:"foreignKey:OrderID"`
}


type OrderItem struct {
	OrderItemID		uint		`gorm:"primaryKey"`
	OrderID			uint		`gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ItemName		string		`gorm:"size:50"`
	ItemPrice		float64		`gorm:"type:decimal(10,2)"`
}

type Review struct {
	ReviewID		uint		`gorm:"primaryKey"`
	OrderItemID		uint		`gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID			uint		`gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Rating			int			`gorm:"type:int;check:rating >= 1 AND rating <= 5"`
	ReviewText		string		`gorm:"type:text"`
	ReviewDate		time.Time	`gorm:"autoCreateTime"`
}

func Init() {
	config.Connect()
	db = config.GetDB()
	err := db.AutoMigrate(&User{}, &Conversation{}, &Message{}, &Order{}, &OrderItem{}, &Review{})
	if err != nil {
		fmt.Println("Error migrating models:", err)
		panic(err)
	}

	var count int64
	db.Model(&User{}).Count(&count)
	if count == 0 {
		err := insertData()
		if err != nil {
			log.Fatalf("Problem with insertData: %v", err)
		}
	}
}


func insertData() error{
	user := User{
		Name: "John Doe",
		Email: "john.doe@example.com",
	}
	if err := db.Create(&user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	fmt.Printf("Create user with ID: %d\n", user.UserID)

	conversation := Conversation{
		UserID: user.UserID,
		EndedAt: time.Now(),
	}
	if err := db.Create(&conversation).Error; err != nil {
        return fmt.Errorf("failed to create conversation: %w", err)
    }

	//messages to conversations:
	messages := []Message{
		{ConversationID: conversation.ConversationID, Sender: "system", MessageText: "Hello John, can you please tell us a bit about the iPhone13 you recently received?"},
		{ConversationID: conversation.ConversationID, Sender: "user", MessageText: "Good overall not many problems"},
		{ConversationID: conversation.ConversationID, Sender: "system", MessageText: "Thank you for your feedback! If you need help with anything please let us now!"},
	}

	if err := db.Create(&messages).Error; err != nil {
		return fmt.Errorf("failed to create messages: %w", err)
	}

	order := Order{
		UserID: user.UserID,
	}
	if err := db.Create(&order).Error; err != nil {
        return fmt.Errorf("failed to create order: %w", err)
    }

	orderItems := []OrderItem{
		{OrderID: order.OrderID,
		ItemName: "iPhone-13",
		ItemPrice: 649.00,},

		{OrderID: order.OrderID,
		ItemName: "GPU",
		ItemPrice: 2247.49,},
	}
	if err := db.Create(&orderItems).Error; err != nil {
        return fmt.Errorf("failed to create order items: %w", err)
    }

	reviews := []Review{
		{OrderItemID: orderItems[0].OrderItemID,
		UserID: user.UserID,
		Rating: 4,
		ReviewText: "Good overall not many problems",},
	}
	if err := db.Create(&reviews).Error; err != nil {
        return fmt.Errorf("failed to create reviews: %w", err)
    }

	return nil
}