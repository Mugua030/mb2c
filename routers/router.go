package routers

//@author zhouping :: will perfect this coninue

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
	beego.Router("/ordering", &controllers.OrderController{}, "post:HandleOrdering")
	beego.Router("/addOrder", &controllers.OrderController{}, "post:HandleAddOrder")
	//购物车 cart
	beego.Router("/cart/add", &controllers.CartController{}, "post:HandleAddCart")
	beego.Router("/cart/list", &controllers.CartController{}, "get:ShowItemListInCart")
	//支付
	beego.Router("/payorder", &controllers.OrderController{}, "get:ShowPayOrder")
	beego.Router("/payOk", &controllers.OrderController{}, "get:PayOk")
	//用户中心
	beego.Router("/ucenter/userinfo", &controllers.UserController{}, "get:ShowUserInfo")
	beego.Router("/ucenter/orderlist", &controllers.UserController{}, "get:ShowUserOrderList")
	beego.Router("/ucenter/siteinfo", &controllers.UserController{}, "get:ShowUserSiteInfo")
	beego.Router("/ucenter/addsite", &controllers.UserController{}, "post:HandleAddSite")

	//后台管理 admin
	beego.Router("/admin", &controllers.AdminGoodsController{}, "get:ShowGoodsList")
	beego.Router("/admin/goods/list", &controllers.AdminGoodsController{}, "get:ShowGoodsList;post:ShowGoodsList")
	beego.Router("/admin/goods/addGoods", &controllers.AdminGoodsController{}, "get:ShowAddGoods;post:HandleAddGoods")
	beego.Router("/admin/goods/addType", &controllers.AdminGoodsController{}, "get:ShowAddGoodsType;post:HandleAddGoodsType")
	beego.Router("/admin/goods/addSPU", &controllers.AdminGoodsController{}, "get:ShowAddGoodsSPU;post:HandleAddGoodsSPU")
}
