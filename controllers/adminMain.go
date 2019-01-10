package controllers

import (
	"github.com/astaxie/beego"
)

type MainAdminController struct {
	beego.Controller
}

func (c *MainAdminController) ShowAdminIndex() {

	//c.Data[]
	c.Layout = "admin_layout.html"
	c.TplName = "admin_index.html"
}
