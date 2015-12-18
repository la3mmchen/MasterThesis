package controllers

import (
	"github.com/revel/revel"
	"fmt"
	//"encoding/json"
	//"database/sql"
	//_ "github.com/go-sql-driver/mysql"
	"time"
	"strconv"
	//"crypto/rand"
	"log"
	)

type Pipe struct {
	*revel.Controller
	PipeId string
	PipeName string
	PipeState string
	PipeLastupdate string
	PipeIdentifier string
	PipeUserUniqueId string
}

func (c Pipe) Pipes() revel.Result {
	if validateSessionPipe(c) {
		uUsr := getUserUniqueIdByName(GlobalUserName)
		fmt.Printf("Pipes()] $uUsr: %+v\n", uUsr)
		s := getPipes(uUsr)
		return c.Render(s, GlobalUserName, GlobalUserSignedIn)
	} else {
		return c.Redirect("/")
	}
}

func (c Pipe) Activate() revel.Result {
	if validateSessionPipe(c) {
		var pipeId string
		c.Params.Bind(&pipeId, "id")
		var cPipe Pipe
		cPipe = loadPipe(pipeId)
		if objectAuth(cPipe) {
			var sTimestamp string
			sTimestamp = strconv.FormatInt(time.Now().Unix(), 10)
			dbHandle = dbConnect()
			defer dbHandle.Close()
			sSql := "UPDATE tblPipe SET PipeLastupdate = '"+sTimestamp+"' WHERE PipeId = '"+pipeId+"';"
			fmt.Printf("Activate()] $sSql: %+v\n", sSql)
			stmtHandle = dbStatement(sSql)
			defer stmtHandle.Close()

			_, err := stmtHandle.Exec()

			if err == nil {
				return c.Redirect("/Pipes")
			} else {
				return c.Redirect("Pipes")
			}
		} else {
			return c.Redirect("/Pipes")
		}
	} else {
		return c.Redirect("/")
	}
}

func (c Pipe) Delete() revel.Result {
	if validateSessionPipe(c) {
		var pipeId string
		c.Params.Bind(&pipeId, "id")
		dbHandle = dbConnect()
		defer dbHandle.Close()
		sSql := "DELETE FROM tblPipe WHERE PipeId = '"+pipeId+"';"
		fmt.Printf("Delete()] $sSql: %+v\n", sSql)
		stmtHandle = dbStatement(sSql)
		defer stmtHandle.Close()

		_, err := stmtHandle.Exec()

		if err == nil {
			return c.Redirect("/Pipes")
		} else {
			return c.Redirect("Pipes")
		}
	} else {
		return c.Redirect("/")
	}
}

func (c Pipe) New() revel.Result {
	if validateSessionPipe(c) {
		return c.Render(GlobalUserName, GlobalUserSignedIn)
	}
	return c.Redirect("/")
}

func (c Pipe) NewPost() revel.Result {
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
		return c.RenderTemplate("Pipe/NewPipe.html")
	}
}

func createPipe(p Pipe) bool {
	p.PipeId = randString(10)
	uUsr := getUserUniqueIdByName(GlobalUserName)
	dbHandle = dbConnect()
	defer dbHandle.Close()
	// Prepare statement for reading data
	sSql := "INSERT INTO tblPipe (PipeId, PipeName, PipeState, PipeIdentifier, PipeUserUniqueId, PipeTypesUniqueId) VALUES ('"+p.PipeId+"', '"+p.PipeName+"', '"+p.PipeState+"', '"+p.PipeIdentifier+"', '"+uUsr+"', '"+typesContext+"');"
	fmt.Printf("createPipe()] $sSql: %+v\n", sSql)
	stmtHandle = dbStatement(sSql)
	defer stmtHandle.Close()

	_, err := stmtHandle.Exec()

	if err == nil {
		return true
	} else {
		return false
	}
}

func getPipes(i string) []Pipe {
		s := make([]Pipe, 0)
		dbHandle = dbConnect()
		defer dbHandle.Close()

		rows, err := dbHandle.Query("SELECT PipeId, PipeName, PipeState, PipeLastupdate, PipeIdentifier FROM tblPipe WHERE PipeUserUniqueId = ? order by PipeLastupdate DESC", i)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var p Pipe
			err := rows.Scan(&p.PipeId, &p.PipeName, &p.PipeState, &p.PipeLastupdate, &p.PipeIdentifier)
			if err != nil {
				log.Fatal(err)
			}
			s = append(s, p)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	return s
}

// GET /pipe/:id
func (c Pipe) GetPipe() revel.Result {
	var id string
	c.Params.Bind(&id, "id")
	var pipe Pipe
  var tblPipeId string
	var tblPipeName string
	var tblPipeState string
	var tblPipeLastupdate string
	var tblPipeIdentifier string

	dbHandle = dbConnect()
	defer dbHandle.Close()

	// Prepare statement for reading data
	sSql := "SELECT PipeId, PipeName, PipeState, PipeLastupdate, PipeIdentifier FROM tblPipe WHERE PipeId = ?"
	stmtHandle = dbStatement(sSql)
	defer stmtHandle.Close()

	// Construct output
	c.Response.Out.Header().Add("Content-Type", "application/json")

	err := stmtHandle.QueryRow(id).Scan(&tblPipeId,&tblPipeName,&tblPipeState,&tblPipeLastupdate,&tblPipeIdentifier)
	if err != nil { // Object not found
		c.Response.Status = 404
		return c.RenderJson("")
	} else { // Object found
		pipe.PipeId = tblPipeId
		pipe.PipeName = tblPipeName
		pipe.PipeState = tblPipeState
		pipe.PipeLastupdate = tblPipeLastupdate
		pipe.PipeIdentifier = tblPipeIdentifier
		c.Response.Out.Header().Add("Location", "/object/"+pipe.PipeId)
		c.Response.Status = 200
		return c.RenderJson(pipe)
	}
}

func loadPipe(i string) Pipe {
		var y Pipe
		fmt.Printf("loadPipe()] $i= %+v\n", i)
		dbHandle = dbConnect()
		defer dbHandle.Close()
		fmt.Printf("loadPipe()]: $i: %+v\n", i)
		// Prepare statement for reading data
		sSql := "Select PipeId, PipeName, PipeState, PipeLastupdate, PipeIdentifier, PipeUserUniqueId FROM tblPipe WHERE PipeId = ?"
		stmtHandle = dbStatement(sSql)
		defer stmtHandle.Close()

		err := stmtHandle.QueryRow(i).Scan(&y.PipeId, &y.PipeName, &y.PipeState, &y.PipeLastupdate, &y.PipeIdentifier, &y.PipeUserUniqueId)
		fmt.Printf("loadPipe()]: err: %+v\n", err)
		return y
}


func validateSessionPipe(c Pipe) bool {
	fmt.Printf("validateSessionPipe()] Session User: %+v\n", c.Session["user"])
	fmt.Printf("validateSessionPipe()] Session Time: %+v\n", c.Session["_TS"])

	if  c.Session["user"] != "" {
		GlobalUserName = c.Session["user"]
		GlobalUserSignedIn = true
	}
	return GlobalUserSignedIn
}

func objectAuth(x Pipe) bool {
	var bAuthed = false
	fmt.Printf("x.PipeUserUniqueId: %+v\n", x.PipeUserUniqueId)
	if x.PipeUserUniqueId == getUserUniqueIdByName(GlobalUserName) {
			bAuthed = true
	}
	return bAuthed
}
