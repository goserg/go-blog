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
	db := getDatabase()
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

	log.Fatal(http.ListenAndServe(getPort(), nil))
}

func getPort() string {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8000"
	}
	return ":" + port
}

func getDatabase() *sql.DB {
	db, err := sql.Open("postgres", getDatabaseURL())
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func getDatabaseURL() string {
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
	return databaseURL
}
