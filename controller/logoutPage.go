package controller

import (
	"fmt"
	"html/template"
	"net/http"
)

//LogoutPage Страница логаут
func (c *Controller) LogoutPage(w http.ResponseWriter, r *http.Request) {
	data := loginPageData{User: c.getUserFromCookies(r)}
	if data.User.UserName == "Unauthorized" {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if r.Method == http.MethodPost {
		cookie := http.Cookie{Name: "token", Value: "", Path: "/"}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusFound)
	}

	files := []string{
		"templates/base.gohtml",
		"templates/logout.content.gohtml",
		"templates/header.gohtml",
	}

	tmpl := template.Must(template.ParseFiles(files...))
	err := tmpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
}
