package controllers

import (
	"github.com/revel/revel"
	"fmt"
	"log"
	"golang.org/x/crypto/bcrypt"
	)

type User struct {
	*revel.Controller
	id int

	UserId string
	UserName string
	UserLastPipe Pipe
	UserHashedPassword []byte
}

type UserOut struct {
	*revel.Controller
	id int

	UserId string
	UserName string
	UserLastPipe Pipe
}

func (c User) List() revel.Result {
	if validateSessionUser(c) {
		users := getUsers(50)
		return c.Render(users, GlobalUserName, GlobalUserSignedIn)
	} else {
		return c.Redirect("/")
	}
}

// GET /user/:id
func (c User) GetUser(apiKey string) revel.Result {
	var id string
	c.Params.Bind(&id, "id")
	c.Params.Bind(&apiKey, "apiKey")

	fmt.Printf("GetUser()] $apiKey: %+v\n", apiKey)
	fmt.Printf("GetUser()] $riddlApiKeys[apiKey]: %+v\n", riddlApiKeys["8r23h5rutgnefduib"])

	if riddlApiKeys[apiKey] == 1 || true {
		var userObject UserOut
		var tblUserName string
		var tblUserId string

		dbHandle = dbConnect()
		defer dbHandle.Close()

		// Prepare statement for reading data
		sSql := "SELECT UserName, UserId FROM tblUser WHERE UserId = ?"
		stmtHandle = dbStatement(sSql)
		defer stmtHandle.Close()

		// Construct output
		c.Response.Out.Header().Add("Content-Type", "application/json")

		err := stmtHandle.QueryRow(id).Scan(&tblUserName,&tblUserId)
		if err != nil { // User not found
			c.Response.Status = 404
			return c.RenderJson("")
		} else { // User found
			userObject.UserName = tblUserName
			userObject.UserId = id
			userObject.UserLastPipe = getActivePipeByUser(getUserUniqueIdByUserId(userObject.UserId))
			c.Response.Out.Header().Add("Location", "/object/"+userObject.UserId)
			c.Response.Status = 200
			return c.RenderJson(userObject)
		}
	} else {
		c.Response.Status = 403
		return c.RenderJson("")
	}
}

func getUserById(i string) User {
	var u User
	fmt.Printf("$i= %+v\n", i)
	dbHandle = dbConnect()
	defer dbHandle.Close()

	// Prepare statement for reading data
	sSql := "SELECT UserName, UserId FROM tblUser WHERE UserUniqueId = ?"
	stmtHandle = dbStatement(sSql)
	defer stmtHandle.Close()

	err := stmtHandle.QueryRow(i).Scan(&u.UserName,&u.UserId)
	fmt.Printf("err: %+v\n", err)
	return u
}

func getUserByName(i string) User {
	var u User
	fmt.Printf("$i= %+v\n", i)
	dbHandle = dbConnect()
	defer dbHandle.Close()

	// Prepare statement for reading data
	sSql := "SELECT UserName, UserId, UserHashedPassword FROM tblUser WHERE UserName = ?"
	stmtHandle = dbStatement(sSql)
	defer stmtHandle.Close()

	err := stmtHandle.QueryRow(i).Scan(&u.UserName,&u.UserId,&u.UserHashedPassword)
	fmt.Printf("err: %+v\n", err)
	return u
}

func getUserUniqueIdByName(i string) string {
	var str string
	fmt.Printf("$i= %+v\n", i)
	dbHandle = dbConnect()
	defer dbHandle.Close()

	// Prepare statement for reading data
	sSql := "SELECT UserUniqueId FROM tblUser WHERE UserName = ?"
	stmtHandle = dbStatement(sSql)
	defer stmtHandle.Close()

	err := stmtHandle.QueryRow(i).Scan(&str)
	fmt.Printf("err: %+v\n", err)
	return str
}

func getUserUniqueIdByUserId(i string) string {
	var userUniqueId string
	dbHandle = dbConnect()
	defer dbHandle.Close()

	// Prepare statement for reading data
	sSql := "SELECT UserUniqueId FROM tblUser WHERE UserId = ?"
	stmtHandle = dbStatement(sSql)
	defer stmtHandle.Close()

	err := stmtHandle.QueryRow(i).Scan(&userUniqueId)
	fmt.Printf("err: %+v\n", err)
	return userUniqueId
}

func getUsers(sIn int) []User {
		s := make([]User, 0)
		uUsr := getUserUniqueIdByName(GlobalUserName)

		dbHandle = dbConnect()
		defer dbHandle.Close()

		rows, err := dbHandle.Query("SELECT UserId, UserName FROM tblUser WHERE userUniqueId NOT LIKE '"+uUsr+"' LIMIT ?", sIn)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var uUserLocal User
			err := rows.Scan(&uUserLocal.UserId, &uUserLocal.UserName)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("noContactsYet(uUserLocal.UserId): %+v\n", noContactsYet(uUserLocal.UserId))

			if noContactsYet(uUserLocal.UserId) {
				s = append(s, uUserLocal)
			}
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	return s
}

func getActivePipeByUser(sIn string) Pipe {
	var userPipe Pipe
	dbHandle = dbConnect()
	defer dbHandle.Close()

	// Prepare statement for reading data
	sSql := "SELECT PipeId, PipeName, PipeState, PipeLastupdate, PipeIdentifier FROM tblPipe WHERE PipeUserUniqueId = '"+sIn+"' ORDER BY `tblPipe`.`PipeLastupdate` DESC LIMIT 1;"
	fmt.Printf("sSql: %+v\n", sSql)
	stmtHandle = dbStatement(sSql)
	defer stmtHandle.Close()

	err := stmtHandle.QueryRow().Scan(&userPipe.PipeId,&userPipe.PipeName,&userPipe.PipeState,&userPipe.PipeLastupdate,&userPipe.PipeIdentifier,)
	fmt.Printf("err: %+v\n", err)
	return userPipe
}

func createUser(n string, p string) bool {
	if GlobalUserSignedIn {
		return false
	} else {
		uUserId := randString(10)
		var pHashed []byte
		pHashed, _ = bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
		dbHandle = dbConnect()
		defer dbHandle.Close()
		// Prepare statement for reading data
		sSql := "INSERT INTO tblUser (UserId, UserName, UserHashedPassword) VALUES ('"+uUserId+"', '"+n+"', '"+string(pHashed)+"');"
		fmt.Printf("sSql: %+v\n", sSql)
		stmtHandle = dbStatement(sSql)
		defer stmtHandle.Close()

		_, err := stmtHandle.Exec()

		if err == nil {
			return true
		} else {
			return false
		}
	}
}


func validateSessionUser(c User) bool {
	fmt.Printf("validateSessionUser()] Session User: %+v\n", c.Session["user"])
	fmt.Printf("validateSessionUser()] Session Time: %+v\n", c.Session["_TS"])

	if  c.Session["user"] != "" {
		GlobalUserName = c.Session["user"]
		GlobalUserSignedIn = true
	}
	return GlobalUserSignedIn
}
