package controllers

import (
	"github.com/revel/revel"
	"math/rand"
	"encoding/json"
	"time"
	"strconv"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	)

type Object struct {
	*revel.Controller
	id int
}
// Object from Request
type JsonObject struct {
	Id	string
	Name string
}
// Object to Write into DB
type Person struct {
	UniqId	int
	Id			string
	Name 		string
}
// GET /object/:id
func (c Object) Get() revel.Result {
	stringArray := [4]string{"RandomName 4711", "Lorem Ipsum Name", "Some random Name", "Mr. Obi Wan"}
	var id string
	c.Params.Bind(&id, "id")
	var jsonForResponse JsonObject
	var httpStatus int
	if rand.Intn(5) > 1 { // Randomize HTTP.200 or HTTP.404
		httpStatus = 200
		jsonForResponse.Id = id
		jsonForResponse.Name = stringArray[rand.Intn(len(stringArray))]
	} else {
		httpStatus = 404
	}
	// Creates HTTP Response
	c.Response.Out.Header().Add("Content-Type", "application/json")
	c.Response.Out.Header().Add("Location", "/object/"+jsonForResponse.Id)
	c.Response.Status = httpStatus
	return c.RenderJson(jsonForResponse)
}

// POST /object/:id
// JSON Example: {"Name": "alex"}
func (c Object) Post(object string) revel.Result {
	var jsonFromRequest JsonObject
	err := json.Unmarshal([]byte(object), &jsonFromRequest)
	var id string
	c.Params.Bind(&id, "id")
	jsonFromRequest.Id = id

	if err != nil || len(object) == 0 {
		c.Response.Status = 400
		return c.RenderJson(nil)
	} else {
		c.Response.Out.Header().Add("Content-Type", "application/json")
		c.Response.Out.Header().Add("Location", "/object/"+jsonFromRequest.Id)
		c.Response.Status = 201
		return c.RenderJson(jsonFromRequest)
	}
}
// POST /object/new
// JSON Example: {"Name": "alex"}
func (c Object) PostNew(object string) revel.Result {
	stringArray := []string{"Name1", "Name2", "Name3", "Name4"} // Defines some forbidden strings
	var jsonFromRequest JsonObject
	err := json.Unmarshal([]byte(object), &jsonFromRequest)
	newId := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	jsonFromRequest.Id = newId

	if err != nil || len(object) == 0 || stringInSlice(jsonFromRequest.Name, stringArray)  {
		c.Response.Status = 400
		return c.RenderJson(nil)
	} else {
		db, err := sql.Open("mysql", "root:pass@/express")
		if err != nil {
		    panic(err.Error())
		}
		defer db.Close()
		var tblPersonNameCount int
		// Prepare statement for reading data
		stmtOut, err := db.Prepare("SELECT Count(Name) FROM tbl_person WHERE Name = ?")
		if err != nil {
					panic(err.Error()) // proper error handling instead of panic in your app
			}
		defer stmtOut.Close()


		err = stmtOut.QueryRow(jsonFromRequest.Name).Scan(&tblPersonNameCount)
		if err != nil {
				panic(err.Error()) // proper error handling instead of panic in your app
		}

		if tblPersonNameCount > 0 {
			c.Response.Status = 400
			return c.RenderJson(nil)
		} else {
			c.Response.Out.Header().Add("Content-Type", "application/json")
			c.Response.Out.Header().Add("Location", "/object/"+newId)
			c.Response.Status = 201
			return c.RenderJson(jsonFromRequest)
		}
	}
}
// POST /object/write
// JSON Example: {"Name": "alex"}
func (c Object) WriteDb(object string) revel.Result {
	var jsonFromRequest JsonObject
	var person Person
	err := json.Unmarshal([]byte(object), &jsonFromRequest)

	newId := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	jsonFromRequest.Id = newId

	if err != nil || len(object) == 0   {
		c.Response.Status = 400
		return c.RenderJson("")
	} else {

		db, err := sql.Open("mysql", "root:pass@/express")
		if err != nil {
		    panic(err.Error())
		}
		defer db.Close()
		// Prepare Insert
		stmtIns, err := db.Prepare("INSERT INTO tbl_person VALUES(?, ?, ? )")
		if err != nil {
		    panic(err.Error())
		}
		defer stmtIns.Close()
		// Execute Insert
		_, err = stmtIns.Exec("", jsonFromRequest.Id, jsonFromRequest.Name)
		if err != nil {
		    panic(err.Error())
		}

		var tblPersonId int
		// Prepare statement for reading data
		stmtOut, err := db.Prepare("SELECT UniqId FROM tbl_person WHERE Id = ?")
		if err != nil {
					panic(err.Error()) // proper error handling instead of panic in your app
			}
		defer stmtOut.Close()

		err = stmtOut.QueryRow(jsonFromRequest.Id).Scan(&tblPersonId)
		if err != nil {
				panic(err.Error()) // proper error handling instead of panic in your app
		}

		person.Id = jsonFromRequest.Id
		person.Name = jsonFromRequest.Name
		person.UniqId = tblPersonId
		// Prepare Response
		c.Response.Out.Header().Add("Content-Type", "application/json")
		c.Response.Out.Header().Add("Location", "/object/"+jsonFromRequest.Id)
		c.Response.Status = 201

		return c.RenderJson(person)
	}
	return c.RenderJson(jsonFromRequest)
}



func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}
