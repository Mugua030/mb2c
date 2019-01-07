package main

import (
	_ "mb2c/routers"
	"github.com/astaxie/beego"
	_ "mb2c/models"
)

func main() {
	beego.Run()
}

