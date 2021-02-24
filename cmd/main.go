package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/goserg/microblog/server/controller"
)

const (
	user     = "postgres"
	password = "pass"
	host     = "db"
	port     = 5432
	dbname   = "postgres"
)

func main() {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		user,
		password,
		host,
		port,
		dbname))
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

	http.ListenAndServe(":8080", nil)
}
