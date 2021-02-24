package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/goserg/microblog/server"
)

type homePageData struct {
	User      server.User
	Errors    []string
	Posts     []server.Post
	Page      int
	NextPages []int
	PrevPages []int
}

//HomePage - главная страница
func (c *Controller) HomePage(w http.ResponseWriter, r *http.Request) {
	data := homePageData{
		User:   c.getUserFromCookies(r),
		Errors: []string{"testErr", "testErr2"},
	}
	var page int64
	page, _ = strconv.ParseInt(strings.Split(r.URL.Path, "/")[1], 10, 64)

	data.Page = int(page)
	data.Posts = c.getPosts(data.Page)
	var postsCount int
	_ = c.db.QueryRow(`select count(*) from posts`).Scan(&postsCount)
	totalPages := postsCount/10 + 1

	if data.Page <= 0 {
		data.Page = 1
	}
	if data.Page > totalPages {
		data.Page = totalPages + 1
	}

	for i := 1; i < data.Page; i++ {
		data.PrevPages = append(data.PrevPages, i)
	}
	for i := data.Page + 1; i <= totalPages; i++ {
		data.NextPages = append(data.NextPages, i)
	}

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

func (c *Controller) getPosts(page int) []server.Post {
	postsOnPage := 10
	if page < 1 {
		page = 1
	}
	rows, err := c.db.Query(`select * from posts order by "time" DESC OFFSET $1 LIMIT $2`, (page-1)*postsOnPage, postsOnPage)
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
		post.Strings = strings.Split(post.Text, "\n")
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
