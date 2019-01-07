package controllers

import (
	"github.com/astaxie/beego"
)

type GoodsController struct {
	beego.Controller
}

func (c *GoodsController) ShowIndex() {
	userName := c.GetSession("username")
	if userName == nil {
		c.Data["username"] = ""
	} else {
		c.Data["username"] = userName.(string)
	}
	c.TplName = "index.html"
}
