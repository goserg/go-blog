package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/goserg/microblog/controller"
)

func main() {
	databaseURL, exists := os.LookupEnv("DATABASE_URL")
	if !exists {
		const (
			user     = "postgres"
			password = "pass"
			host     = "localhost"
			port     = 5432
			dbname   = "postgres"
		)
		databaseURL = fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
			user,
			password,
			host,
			port,
			dbname,
		)
	}
	db, err := sql.Open("postgres", databaseURL)
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

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8000"
	}

	http.ListenAndServe(":"+port, nil)
}
