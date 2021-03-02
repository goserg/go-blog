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

	if r.Method == http.MethodPost {
		r.ParseForm()
		fmt.Println(r.Form)
		comment := models.Comment{
			Text:    r.Form["text"][0],
			RefPost: r.Form["postID"][0],
		}
		c.postComment(comment, data.User)
	}

	var page int64
	page, _ = strconv.ParseInt(strings.Split(r.URL.Path, "/")[1], 10, 64)

	data.Page = int(page)
	data.Posts = c.getPosts(data.Page)
	var postsCount int
	_ = c.db.QueryRow(`select count(*) from post`).Scan(&postsCount)
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
		"templates/post.gohtml",
		"templates/comment.gohtml",
	}

	tmpl := template.Must(template.New("base.gohtml").Funcs(template.FuncMap{
		"isAuthorized": func() bool {
			return data.User.Authorized
		},
	}).ParseFiles(files...))

	err := tmpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Controller) postComment(comment models.Comment, user models.User) {
	c.db.Exec(`insert into comment ("time", "text", "author", "post") values(NOW(), $1, $2, $3)`, comment.Text, user.ID, comment.RefPost)
}

func (c *Controller) getPosts(page int) []models.Post {
	postsOnPage := 10
	if page < 1 {
		page = 1
	}
	rows, err := c.db.Query(`select * from post order by "time" DESC OFFSET $1 LIMIT $2`, (page-1)*postsOnPage, postsOnPage)
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
			unit := models.PostUnit{}
			firstLink := strings.Index(post.Text, "[Link]")
			firstBr := strings.Index(post.Text, "\n")

			if firstLink == -1 && firstBr == -1 {
				unit.Text = post.Text
				unit.UnitType = "text"
				post.Text = ""
				post.Units = append(post.Units, unit)
				continue
			}
			if firstBr >= 0 && firstBr < firstLink || firstLink == -1 {
				unit.Text = post.Text[0:firstBr]
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
				indexOfLinkEnd := strings.Index(post.Text, "[/Link]")
				unit.Link = post.Text[6:indexOfText]
				unit.Text = post.Text[indexOfText+10 : indexOfLinkEnd]
				post.Units = append(post.Units, unit)

				post.Text = post.Text[indexOfLinkEnd+7 : len(post.Text)]
				continue
			}
		}

		post.Comments = c.getComments(post)

		posts = append(posts, post)
	}
	return posts
}

func (c *Controller) getComments(post models.Post) []models.Comment {
	comments := []models.Comment{}

	res, err := c.db.Query(`select * from comment where post=$1 order by time desc`, post.ID)
	if err != nil {
		fmt.Print("error getting comments")
		fmt.Println(err)
	}
	for res.Next() {
		comment := models.Comment{}
		var authorID int64
		res.Scan(&comment.ID, &comment.Time, &comment.Text, &authorID, &comment.RefPost, &comment.RefComemnt)
		comment.Author = c.getAuthorFromDBByID(uint64(authorID))
		comment.Time = strings.Replace(strings.Split(comment.Time, ".")[0], "T", " в ", -1)
		comments = append(comments, comment)
	}

	return comments
}

func (c *Controller) getAuthorFromDBByID(id uint64) models.User {
	author := models.User{}
	err := c.db.QueryRow(`select * from auth_user where id=$1`, id).Scan(&author.ID, &author.UserName, &author.PasswordHash)
	if err != nil {
		fmt.Println(err)
	}
	return author
}
