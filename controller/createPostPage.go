package controller

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/goserg/microblog/server"
)

//CreatePostPage Страница логина
func (c *Controller) CreatePostPage(w http.ResponseWriter, r *http.Request) {
	data := loginPageData{User: c.getUserFromCookies(r)}

	if !data.User.Authorized {
		http.Redirect(w, r, "/login/", http.StatusFound)
	}

	if r.Method == http.MethodPost {
		r.ParseForm()

		title := r.Form["title"][0]
		text := r.Form["text"][0]
		fmt.Println(title)
		fmt.Println(text)

		post := server.Post{
			Title:  title,
			Text:   text,
			Author: data.User,
		}
		c.insertPostToDB(&post)
	}

	files := []string{
		"templates/base.gohtml",
		"templates/create.content.gohtml",
		"templates/header.gohtml",
	}

	tmpl := template.Must(template.ParseFiles(files...))
	tmpl.Execute(w, data)
}

func (c *Controller) insertPostToDB(post *server.Post) {
	_, err := c.db.Exec(`insert into posts ("time", "title", "text", "author") values(NOW(), $1, $2, $3)`, post.Title, post.Text, post.Author.ID)
	if err != nil {
		fmt.Println(err)
	}
}
