package controllers

import (
	"github.com/revel/revel"
	"fmt"
	//"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"strconv"
	"crypto/rand"
	//"log"
	)

type Object struct {
	*revel.Controller
	id int
}

// :Types
type Types struct {
	TypesId string
	TypesName string
}

var dbHandle *sql.DB
var stmtHandle *sql.Stmt

// GET /types/:id
func (c Object) GetTypes() revel.Result {
	var id string
	c.Params.Bind(&id, "id")
	var types Types
  var tblTypesName string
	var tblTypesId string

	dbHandle = dbConnect()
	defer dbHandle.Close()

	// Prepare statement for reading data
	sSql := "SELECT TypesId, TypesName FROM tblTypes WHERE TypesId = ?"
	stmtHandle = dbStatement(sSql)
	defer stmtHandle.Close()

	// Construct output
	c.Response.Out.Header().Add("Content-Type", "application/json")

	err := stmtHandle.QueryRow(id).Scan(&tblTypesId,&tblTypesName)
	if err != nil { // User not found
		c.Response.Status = 404
		return c.RenderJson("")
	} else { // User found
		types.TypesId = tblTypesId
		types.TypesName = tblTypesName
		c.Response.Out.Header().Add("Location", "/object/"+types.TypesId)
		c.Response.Status = 200
		return c.RenderJson(types)
	}
}

// POST /new
func (c Object) New() revel.Result {
	var ObjectType string
	c.Params.Bind(&ObjectType, "ObjectType")

	switch ObjectType {
		case "pipe":
			var pipe Pipe
			var PipeName string
			c.Params.Bind(&PipeName, "PipeName")
			var PipeIdentifier string
			c.Params.Bind(&PipeIdentifier, "PipeIdentifier")
			pipe.PipeName = PipeName
			pipe.PipeState = "unknown"
			pipe.PipeLastupdate = strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
			pipe.PipeIdentifier = PipeIdentifier
			bCreate := createPipe(pipe)
			if bCreate {
				return c.Redirect("Pipes")
			} else {
				return c.RenderTemplate("app/NewPipe.html")
			}

		default:
			ObjectType = "nil"
			return c.RenderJson("")
	}
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}


// ** Generic globals **/
// creates a random string
func randString(n int) string {
    const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    var bytes = make([]byte, n)
    rand.Read(bytes)
    for i, b := range bytes {
        bytes[i] = alphanum[b % byte(len(alphanum))]
    }
		fmt.Printf("randString: %+v\n", string(bytes))
    return string(bytes)
}
// ** db calls **
func dbConnect() *sql.DB {
	db, err := sql.Open("mysql", "root:pass@/riddl")
	if err != nil {
			panic(err.Error())
	}
	return db
}
func dbStatement(s string) *sql.Stmt {
	stmtOut, err := dbHandle.Prepare(s)
	if err != nil {
				panic(err.Error()) // proper error handling instead of panic in your app
	}

	return stmtOut
}
