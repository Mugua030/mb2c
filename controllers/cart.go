package controllers

import (
	"mb2c/models"
	"strconv"

	"github.com/astaxie/beego/orm"

	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
)

type CartController struct {
	beego.Controller
}

//HandleAddCart  添加购物车
func (c *CartController) HandleAddCart() {

	resp := make(map[string]interface{})

	goodId, err := c.GetInt("good_id")
	if err != nil {
		resp["errMsg"] = "add to cart fail"
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}
	count, err := c.GetInt("count")

	if err != nil {
		resp["errMsg"] = "add to cart fail"
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}

	redisAddr := beego.AppConfig.String("redis::redis_host") + ":" + beego.AppConfig.String("redis::redis_port")
	conn, err := redis.Dial("tcp", redisAddr)
	if err != nil {
		resp["errMsg"] = "数据连接出错redis"
		resp["errNo"] = "2"
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}

	userName := c.GetSession("username")
	if userName == nil {
		resp["errMsg"] = "请登录"
		resp["errNo"] = 3
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}
	beego.Error(userName)
	hkey := "cart_" + userName.(string)
	_, err = conn.Do("hset", hkey, goodId, count)
	if err != nil {
		resp["errMsg"] = "插入数据失败"
		resp["errNo"] = "2"
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}

	resp["errMsg"] = "添加成功"
	resp["errNo"] = 0
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *CartController) ShowItemListInCart() {
	userName := c.GetSession("username")
	if userName == nil {
		c.Redirect("/index", 302)
		return
	}
	rKey := "cart_" + userName.(string)
	redisAddr := beego.AppConfig.String("redis::redis_host") + ":" + beego.AppConfig.String("redis::redis_port")
	conn, err := redis.Dial("tcp", redisAddr)
	resp, err := conn.Do("hgetall", rKey)
	res, err := redis.IntMap(resp, err)
	o := orm.NewOrm()
	var goodsList []map[string]interface{}
	for goodId, nums := range res {
		//读取商品信息
		var goods models.GoodsSKU
		goods.Id, _ = strconv.Atoi(goodId)
		err := o.Read(&goods)
		if err != nil {
			beego.Warning("cartlist:: no this goods:", goodId)
			continue
		}
		tmp := make(map[string]interface{})
		tmp["goods"] = goods
		tmp["totalPrice"] = nums * goods.Price
		goodsList = append(goodsList, tmp)
	}

	c.Data["img_url"] = beego.AppConfig.String("imgURL::img_url")
	c.Data["username"] = userName
	c.Data["goodsList"] = goodsList

	c.Layout = "layout.html"
	c.TplName = "cart.html"
}
