package main

import(
	"fmt"
	"log"
	"net/http"
	"html/template"
)

var tpl *template.Template

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	var err error
	tpl, err = template.ParseGlob(".../static/templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/information", informationHandler)
	http.HandleFunc("/interactions", interactionHandler)
	http.HandleFunc("/reviews", reviewHandler)

	fmt.Println("Listening on port 8080...")
	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}


func mainHandler(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "main.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func informationHandler(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "information.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


func interactionHandler(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "interactions.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


func reviewHandler(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "reviews.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

