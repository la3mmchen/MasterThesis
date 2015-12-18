package controllers

import (
	"github.com/revel/revel"
	"log"
	"fmt"
	"strconv"
	"time"
	)

type Contact struct {
	*revel.Controller
	id int
}

// ** business logic **
func (c Contact) List() revel.Result {
	if validateSessionContact(c) {
		uUsr := getUserUniqueIdByName(GlobalUserName)
		s := getContacts(uUsr)
		return c.Render(s, GlobalUserName, GlobalUserSignedIn)
	} else {
		return c.Redirect("/")
	}
}

func (c Contact) Uncontact() revel.Result {
	if validateSessionContact(c) {
		var userId string
		c.Params.Bind(&userId, "id")
		uUsr := getUserUniqueIdByName(GlobalUserName)
		toUnique := getUserUniqueIdByUserId(userId)

		dbHandle = dbConnect()
		defer dbHandle.Close()

		sSql := "DELETE FROM tblContacts WHERE ContactUserUniqueIdFrom = '"+uUsr+"' and ContactUserUniqueIdTo = '"+toUnique+"';"
		fmt.Printf("sSql: %+v\n", sSql)
		stmtHandle = dbStatement(sSql)
		defer stmtHandle.Close()

		_, err := stmtHandle.Exec()

		if err == nil {
			return c.Redirect("/Contacts")
		} else {
			return c.Redirect("/Contacts")
		}
		return c.Redirect("/Contacts")
	} else {
		return c.Redirect("/")
	}
}

func (c Contact) Add() revel.Result {
	if validateSessionContact(c) {
		var userId string
		c.Params.Bind(&userId, "id")

		var sTimestamp string
		sTimestamp = strconv.FormatInt(time.Now().Unix(), 10)
		toUnique := getUserUniqueIdByUserId(userId)
		uUsr := getUserUniqueIdByName(GlobalUserName)
		dbHandle = dbConnect()
		defer dbHandle.Close()

		sSql := "INSERT INTO tblContacts (ContactUserUniqueIdTo, ContactUserUniqueIdFrom, ContactLastupdate) VALUES ('"+toUnique+"', '"+uUsr+"', '"+sTimestamp+"');"
		fmt.Printf("sSql: %+v\n", sSql)
		stmtHandle = dbStatement(sSql)
		defer stmtHandle.Close()

		_, err := stmtHandle.Exec()

		if err == nil {
			return c.Redirect("/Contacts")
		} else {
			return c.Redirect("/Users")
		}
		return c.Redirect("/Users")
	} else {
		return c.Redirect("/")
	}
}

func getContacts(sIn string) []User {
		s := make([]User, 0)
		dbHandle = dbConnect()
		defer dbHandle.Close()

		rows, err := dbHandle.Query("SELECT ContactUserUniqueIdTo FROM `tblContacts` where ContactUserUniqueIdFrom = ?", sIn)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var uUserLocal User
			var sUserId string
			err := rows.Scan(&sUserId)
			uUserLocal = getUserById(sUserId)
			if err != nil {
				log.Fatal(err)
			}
			uUserLocal.UserLastPipe = getActivePipeByUser(getUserUniqueIdByUserId(uUserLocal.UserId))
			s = append(s, uUserLocal)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	return s
}

func noContactsYet (UserIdIn string) bool {

	var ContactUniqueId string
	uUserLocal := getUserUniqueIdByUserId(UserIdIn)
	uUsr := getUserUniqueIdByName(GlobalUserName)
	dbHandle = dbConnect()
	defer dbHandle.Close()

	// Prepare statement for reading data
	sSql := "SELECT ContactUniqueId FROM tblContacts WHERE  ContactUserUniqueIdFrom = '"+uUsr+"' and ContactUserUniqueIdTo = '"+uUserLocal+"';"
	fmt.Printf("sSql: %+v\n", sSql)
	stmtHandle = dbStatement(sSql)
	defer stmtHandle.Close()
	err := stmtHandle.QueryRow().Scan(&ContactUniqueId)

	if err != nil { // Object not found
		return true
	} else { // Object found
		return false
	}
}

func validateSessionContact(c Contact) bool {
	fmt.Printf("Session User: %+v\n", c.Session["user"])
	fmt.Printf("Session Time: %+v\n", c.Session["_TS"])

	if  c.Session["user"] != "" {
		GlobalUserName = c.Session["user"]
		GlobalUserSignedIn = true
	}
	return GlobalUserSignedIn
}
