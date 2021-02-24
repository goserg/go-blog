package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/goserg/microblog/server"
)

type homePageData struct {
	User   server.User
	Errors []string
	Posts  []server.Post
}

//HomePage - главная страница
func (c *Controller) HomePage(w http.ResponseWriter, r *http.Request) {
	data := homePageData{
		User:   c.getUserFromCookies(r),
		Errors: []string{"testErr", "testErr2"},
	}

	data.Posts = c.getPosts()

	files := []string{
		"templates/base.gohtml",
		"templates/home.content.gohtml",
		"templates/header.gohtml",
	}

	tmpl := template.Must(template.ParseFiles(files...))
	err := tmpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Controller) getPosts() []server.Post {
	rows, err := c.db.Query(`select * from posts`)
	if err != nil {
		fmt.Println(err)
	}
	posts := []server.Post{}
	for rows.Next() {
		post := server.Post{}
		authorID := 0
		rows.Scan(&post.ID, &post.Time, &post.Title, &post.Text, &authorID)
		post.Author = c.getAuthorFromDBByID(uint64(authorID))
		post.Time = strings.Replace(strings.Split(post.Time, ".")[0], "T", " в ", -1)
		posts = append(posts, post)
	}
	return posts
}

func (c *Controller) getAuthorFromDBByID(id uint64) server.User {
	author := server.User{}
	err := c.db.QueryRow(`select * from auth_user where id=$1`, id).Scan(&author.ID, &author.UserName, &author.PasswordHash)
	if err != nil {
		fmt.Println(err)
	}
	return author
}
