package controller

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/goserg/microblog/models"
)

//CreatePostPage Страница логина
func (c *Controller) CreatePostPage(w http.ResponseWriter, r *http.Request) {
	data := loginPageData{User: c.getUserFromCookies(r)}

	if !data.User.Authorized {
		http.Redirect(w, r, "/login/", http.StatusFound)
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
			return
		}
		title := r.Form["title"][0]
		text := r.Form["text"][0]
		fmt.Println(title)
		fmt.Println(text)

		post := models.Post{
			Title:  title,
			Text:   text,
			Author: data.User,
		}
		c.insertPostToDB(&post)

		http.Redirect(w, r, "/", http.StatusFound)
	}

	files := []string{
		"templates/base.gohtml",
		"templates/create.content.gohtml",
		"templates/header.gohtml",
	}

	tmpl := template.Must(template.ParseFiles(files...))
	err := tmpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Controller) insertPostToDB(post *models.Post) {
	_, err := c.db.Exec(`insert into post ("time", "title", "text", "author") values(NOW(), $1, $2, $3)`, post.Title, post.Text, post.Author.ID)
	if err != nil {
		fmt.Println(err)
	}
}
