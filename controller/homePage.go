package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/goserg/microblog/models"
)

type homePageData struct {
	User      models.User
	Errors    []string
	Posts     []models.Post
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

func (c *Controller) getPosts(page int) []models.Post {
	postsOnPage := 10
	if page < 1 {
		page = 1
	}
	rows, err := c.db.Query(`select * from posts order by "time" DESC OFFSET $1 LIMIT $2`, (page-1)*postsOnPage, postsOnPage)
	if err != nil {
		fmt.Println(err)
	}
	posts := []models.Post{}
	for rows.Next() {
		post := models.Post{}
		authorID := 0
		rows.Scan(&post.ID, &post.Time, &post.Title, &post.Text, &authorID)
		post.Author = c.getAuthorFromDBByID(uint64(authorID))
		post.Time = strings.Replace(strings.Split(post.Time, ".")[0], "T", " в ", -1)

		for post.Text != "" {
			fmt.Println(post.Text)
			unit := models.PostUnit{}
			firstLink := strings.Index(post.Text, "[Link]")
			firstBr := strings.Index(post.Text, "\n")
			fmt.Println(firstBr)
			fmt.Println(firstLink)
			if firstLink == -1 && firstBr == -1 {
				fmt.Println("чистый текст")
				unit.Text = post.Text
				unit.UnitType = "text"
				post.Text = ""
				post.Units = append(post.Units, unit)
				continue
			}
			if firstBr >= 0 && firstBr < firstLink || firstLink == -1 {
				fmt.Println("текстовый элемент")
				unit.Text = post.Text[0:firstBr]
				fmt.Println(unit.Text)
				post.Text = post.Text[firstBr+1 : len(post.Text)]
				post.Units = append(post.Units, unit, models.PostUnit{UnitType: "BR"})
				continue
			}
			if firstLink < firstBr || firstBr == -1 {
				if firstLink != 0 {
					unit.Text = post.Text[0:firstLink]
					unit.UnitType = "text"
					post.Units = append(post.Units, unit)
					post.Text = post.Text[firstLink:len(post.Text)]
				}

				unit = models.PostUnit{UnitType: "link"}
				indexOfText := strings.Index(post.Text, "[LinkText]")
				fmt.Printf(`indexOfText: %d`, indexOfText)
				indexOfLinkEnd := strings.Index(post.Text, "[/Link]")
				unit.Link = post.Text[6:indexOfText]
				fmt.Printf(`ссылка: %s`, unit.Link)
				unit.Text = post.Text[indexOfText+10 : indexOfLinkEnd]
				post.Units = append(post.Units, unit)

				post.Text = post.Text[indexOfLinkEnd+7 : len(post.Text)]
				continue
			}
		}
		posts = append(posts, post)
	}
	return posts
}

func (c *Controller) getAuthorFromDBByID(id uint64) models.User {
	author := models.User{}
	err := c.db.QueryRow(`select * from auth_user where id=$1`, id).Scan(&author.ID, &author.UserName, &author.PasswordHash)
	if err != nil {
		fmt.Println(err)
	}
	return author
}
