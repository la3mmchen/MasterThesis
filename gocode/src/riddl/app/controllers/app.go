package controllers

import (
	"golang.org/x/crypto/bcrypt"
	"github.com/revel/revel"
	"fmt"
	)

// Generic Values
var GlobalUserName = ""
var GlobalUserSignedIn = false
var typesContext = "1"

// API Keys
var riddlApiKeys = map[string]int{"8r23h5rutgnefduib": 1, "fsejfbszr13423rfb": 1}

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	validateSession(c)
	uUsr := getUserUniqueIdByName(GlobalUserName)
	return c.Render(uUsr, GlobalUserSignedIn, GlobalUserName)
}

func (c App) Register() revel.Result {
	if validateSession(c) {
		return c.Redirect("/")
	} else {
		return c.RenderTemplate("user/register.html")
	}
}

func (c App) UserLogin(username, password string, remember bool) revel.Result {
	user := getUserByName(username)
	fmt.Printf("user.UserId: %+v\n", user.UserId)
	fmt.Printf("user.UserName: %+v\n", user.UserName)
	if user.UserId != "" { // User found
			/*if user.UserHashedPassword == password {
				c.Session["user"] = username
				c.Session.SetDefaultExpiration()
				c.Flash.Success("Welcome back, " + username)
				return c.Redirect("/Contacts")
			}
			*/
		err := bcrypt.CompareHashAndPassword(user.UserHashedPassword, []byte(password))
		if err == nil {
			c.Session["user"] = username
			c.Session.SetDefaultExpiration()
			c.Flash.Success("Welcome back, " + username)
			return c.Redirect("/Contacts")
		}
		c.Flash.Out["username"] = username
		c.Flash.Error("Login failed")
		return c.Redirect("/")
	} else if password != "" {
		bUserCreated := createUser(username, password)
		if bUserCreated {
			GlobalUserName = c.Session["user"]
			GlobalUserSignedIn = true
			c.Flash.Success("Hej, " + username + " welcome to riddl.me")
			return c.Redirect("/Contacts")
		}
	}
	c.Flash.Out["username"] = username
	c.Flash.Error("Login failed")
	return c.Redirect("/")
}

func (c App) UserLogout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	GlobalUserSignedIn = false
	GlobalUserName = ""
	fmt.Printf("Session ended: %+v\n", c.Session)
	return c.Redirect("/")
}

func validateSession(c App) bool {
	fmt.Printf("Session User: %+v\n", c.Session["user"])
	fmt.Printf("Session Time: %+v\n", c.Session["_TS"])

	if  c.Session["user"] != "" {
		GlobalUserName = c.Session["user"]
		GlobalUserSignedIn = true
	}
	return GlobalUserSignedIn
}
