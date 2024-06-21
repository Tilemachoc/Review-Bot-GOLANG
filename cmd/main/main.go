package main

import(
	"fmt"
	"log"
	"net/http"
	"html/template"
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