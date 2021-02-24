package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/goserg/microblog/server/controller"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	c := controller.NewController(db)
	http.HandleFunc("/", c.HomePage)
	http.HandleFunc("/login/", c.LoginPage)
	http.HandleFunc("/register/", c.RegisterPage)
	http.HandleFunc("/logout/", c.LogoutPage)
	http.HandleFunc("/forgot/", c.ForgotPassPage)
	http.HandleFunc("/create/", c.CreatePostPage)

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	fmt.Println("Server started")

	port := os.Getenv("PORT")

	http.ListenAndServe(":"+port, nil)
}
