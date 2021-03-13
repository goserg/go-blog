package controller

import (
	"fmt"
	"html/template"
	"net/http"
)

//ForgotPassPage Страница логина
func (c *Controller) ForgotPassPage(w http.ResponseWriter, r *http.Request) {
	data := loginPageData{User: c.getUserFromCookies(r)}

	if data.User.Authorized {
		http.Redirect(w, r, "/logout/", http.StatusFound)
	}

	files := []string{
		"templates/base.gohtml",
		"templates/forgot_pass.content.gohtml",
		"templates/header.gohtml",
	}

	tmpl := template.Must(template.ParseFiles(files...))
	err := tmpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
}
