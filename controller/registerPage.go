package controller

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/goserg/microblog/models"
	"github.com/goserg/microblog/utils"
)

type registerPageData struct {
	User     models.User
	UserName string
	Errors   []string
}

//RegisterPage страница регистрации
func (c *Controller) RegisterPage(w http.ResponseWriter, r *http.Request) {
	data := registerPageData{User: c.getUserFromCookies(r)}

	if data.User.Authorized {
		http.Redirect(w, r, "/logout/", http.StatusFound)
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
		}
		userName := r.Form["userName"][0]
		password1 := r.Form["password1"][0]
		password2 := r.Form["password2"][0]

		data.Errors = append(data.Errors, validateRegisterForm(userName, password1, password2)...)

		if len(data.Errors) == 0 {
			err := c.addUserToDB(userName, utils.Hash([]byte(password1)))
			if err != nil {
				if err.Error() == `pq: duplicate key value violates unique constraint "auth_user_name_key"` {
					data.Errors = append(data.Errors, "Имя пользователя не доступно. Выберите другое.")
				}
				err := c.l.Log(fmt.Sprintf("User '%s' not created with error: %s.", userName, err.Error()), utils.ReadUserIP(r))
				if err != nil {
					fmt.Println(err)
				}
			} else {
				err := c.l.Log(fmt.Sprintf("User '%s' created.", userName), utils.ReadUserIP(r))
				if err != nil {
					fmt.Println(err)
				}

				files := []string{
					"templates/base.gohtml",
					"templates/register_success.content.gohtml",
					"templates/header.gohtml",
				}

				data.UserName = userName
				tmpl := template.Must(template.ParseFiles(files...))
				err = tmpl.Execute(w, data)
				if err != nil {
					fmt.Println(err)
				}
				return
			}
		}
	}

	files := []string{
		"templates/base.gohtml",
		"templates/register.content.gohtml",
		"templates/header.gohtml",
	}

	tmpl := template.Must(template.ParseFiles(files...))
	err := tmpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
}

func validateRegisterForm(userName string, password1 string, password2 string) []string {
	var errors []string
	if password1 != password2 {
		errors = append(errors, "Пароли не совпадают")
	}
	if userName == "" || password1 == "" || password2 == "" {
		errors = append(errors, "Заполните все поля")
	}
	minLength := 5
	if len(password1) < minLength {
		errors = append(errors, fmt.Sprintf("Пароль должен быть не менее %d символов", minLength))
	}

	return errors
}
