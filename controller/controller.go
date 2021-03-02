package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/goserg/microblog/logger"
	"github.com/goserg/microblog/models"
	"github.com/goserg/microblog/utils"
)

//Controller is a controller for page handling and sql connection
type Controller struct {
	db *sql.DB
	l  *logger.Logger
}

//NewController - конструктор контроллера
func NewController(db *sql.DB) *Controller {
	l := logger.NewLogger(db)
	return &Controller{db, l}
}

type loginPageData struct {
	User   models.User
	Errors []string
}

func (c *Controller) getUserFromCookies(r *http.Request) models.User {
	user := models.User{ID: 0, UserName: "Unauthorized", PasswordHash: ""}
	cookie, err := r.Cookie("token")
	if err == nil && cookie.Value != "" {
		username, passHash, exp, err := utils.ParseJWT(cookie.Value)
		if err != nil {
			fmt.Println("Error parsing token")
		}
		if exp > time.Now().Unix() {
			cUser, err := c.getUserFromDB(username, passHash)
			if err == nil {
				user = cUser
			}
		}
	}
	return user
}

func (c *Controller) addUserToDB(userName string, password string) error {
	fmt.Println("Adding user to db")
	_, err := c.db.Exec(`insert into auth_user ("name", "pass") values($1, $2)`, userName, password)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (c *Controller) getUserFromDB(userName string, passwordHash string) (models.User, error) {
	user := models.User{Authorized: false}
	err := c.db.QueryRow(
		`select * from "auth_user" where name=$1 and pass=$2`, userName, passwordHash,
	).Scan(&user.ID, &user.UserName, &user.PasswordHash)
	if err == nil {
		user.Authorized = true
	}
	return user, err
}
