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
	//商品相关 Goods
	beego.Router("/goods/detail", &controllers.GoodsController{}, "get:ShowGoodsDetail")
	beego.Router("/goods/type", &controllers.GoodsController{}, "get:ShowGoodsListByType")
	//购物车 cart
	beego.Router("/cart/add", &controllers.CartController{}, "post:HandleAddCart")
	beego.Router("/cart/list", &controllers.CartController{}, "get:ShowItemListInCart")
	//用户中心
	beego.Router("/ucenter/userinfo", &controllers.UserController{}, "get:ShowUserInfo")
	beego.Router("/ucenter/orderlist", &controllers.UserController{}, "get:ShowUserOrderList")
	beego.Router("/ucenter/siteinfo", &controllers.UserController{}, "get:ShowUserSiteInfo")
	beego.Router("/ucenter/addsite", &controllers.UserController{}, "post:HandleAddSite")

	//后台管理 admin
	beego.Router("/admin", &controllers.MainAdminController{}, "get:ShowAdminIndex")
	beego.Router("/admin/addGoods", &controllers.GoodsController{}, "get:ShowAddGoods;post:HandleAddGoods")
}
