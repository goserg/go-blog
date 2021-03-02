package controller

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/goserg/microblog/utils"
)

//LoginPage Страница логина
func (c *Controller) LoginPage(w http.ResponseWriter, r *http.Request) {
	data := loginPageData{User: c.getUserFromCookies(r)}

	if data.User.Authorized {
		http.Redirect(w, r, "/logout/", http.StatusFound)
	}

	if r.Method == http.MethodPost {
		r.ParseForm()

		user, err := c.getUserFromDB(r.Form["userName"][0], utils.Hash([]byte(r.Form["password"][0])))
		if err != nil {
			c.l.Log(fmt.Sprintf("Failed attempt to log (username:'%s', password: '%s') with error: %s", r.Form["userName"][0], r.Form["password"][0], err.Error()), utils.ReadUserIP(r))
			data.Errors = append(data.Errors, "Нет такой комбинации имени пользователя и пароля")
		} else {
			c.l.Log(fmt.Sprintf("'%s' logged in.", user.UserName), utils.ReadUserIP(r))
			cookie := http.Cookie{Name: "token", Value: utils.GenerateNewToken(user), Path: "/"}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}

	files := []string{
		"templates/base.gohtml",
		"templates/login.content.gohtml",
		"templates/header.gohtml",
	}

	tmpl := template.Must(template.ParseFiles(files...))
	tmpl.Execute(w, data)
}
