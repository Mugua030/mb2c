package main

import (
	_ "mb2c/models"
	_ "mb2c/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
