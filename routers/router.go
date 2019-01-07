package routers

import (
	"mb2c/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/register", &controllers.UserController{}, "get:ShowRegister;post:HandleRegister")
	beego.Router("/active", &controllers.UserController{}, "get:HandleActiveUser")
	beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")
	beego.Router("/index", &controllers.GoodsController{}, "get:ShowIndex")
	beego.Router("/logout", &controllers.UserController{}, "get:HandleLogout")
	//用户中心
	beego.Router("/ucenter/userinfo", &controllers.UserController{}, "get:ShowUserInfo")
	beego.Router("/ucenter/orderlist", &controllers.UserController{}, "get:ShowUserOrderList")
	beego.Router("/ucenter/siteinfo", &controllers.UserController{}, "get:ShowUserSiteInfo")
	beego.Router("/ucenter/addsite", &controllers.UserController{}, "post:HandleAddSite")
}
