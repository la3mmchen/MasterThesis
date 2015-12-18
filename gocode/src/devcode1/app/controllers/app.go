package controllers

import "github.com/revel/revel"

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Echo(sInput string) revel.Result{
	return c.Render(sInput)
}

func (c App) PostIndex() revel.Result {
	return c.Render()
}
