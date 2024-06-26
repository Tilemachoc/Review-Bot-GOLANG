package main

import(
	"fmt"
	"log"
	"net/http"
	"html/template"
	"encoding/json"
	"gorm.io/gorm"
	"github.com/Tilemachoc/TASK1/pkg/models"
	"github.com/Tilemachoc/TASK1/pkg/config"
)

var tpl *template.Template
//need db to send the data to templates
var db *gorm.DB


func main() {
	models.Init()
	db = config.GetDB()
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	var err error
	tpl, err = template.ParseGlob("static/templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/users", userHandler)
	http.HandleFunc("/buy", buyHandler)
	http.HandleFunc("/api/reviews", apiReviewHandler)
	http.HandleFunc("/api/orderitems", apiOrderHandler)

	fmt.Println("Listening on port 8080...")
	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}


func mainHandler(w http.ResponseWriter, r *http.Request) {
	product := r.URL.Query().Get("product")
	if product == "" {
		http.Redirect(w,r,"/buy",http.StatusSeeOther)
	}

	if err := tpl.ExecuteTemplate(w, "main.html", product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	//Gorm doesn't put foreignkeys columns data automatically so we can do it like that if we want:
	err := db.Preload("Conversations.Messages").
			Preload("Orders.OrderItems").
			Preload("Reviews").
			Find(&users).Error
	if err != nil {
		log.Fatalf("Error finding users: %v", err)
	}

	//fmt.Printf("Users found: %v\n", users)

	// Since we have one user we just pass the user
	if err := tpl.ExecuteTemplate(w, "users.html", users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


func buyHandler(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "buy.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


func apiReviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var review models.Review
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&review); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := db.Create(&review)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}


	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review)
}


func apiOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var orderitem models.OrderItem
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&orderitem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := db.Create(&orderitem)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(orderitem)
}