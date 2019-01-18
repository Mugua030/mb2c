package controllers

import (
	"fmt"
	"mb2c/models"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"github.com/smartwalle/alipay"
)

type OrderController struct {
	beego.Controller
}

func (c *OrderController) HandleOrdering() {

	goodsSKUIds := c.GetStrings("goodsSKU_id")
	beego.Info("goodsSKUIds", goodsSKUIds)
	userNameI := c.GetSession("username")
	if userNameI == nil {
		beego.Warning("no login user access")
		c.Redirect("/login", 302)
	}
	userName := userNameI.(string)
	beego.Info("userName=", userName)

	o := orm.NewOrm()
	//收获地址列表
	var receivers []models.Receiver
	//	o.QueryTable("Receiver").RelatedSel("User").Filter("User__Name", userName).All(&receivers)
	_, err := o.QueryTable("Receiver").RelatedSel("User").Filter("User__Name", userName).All(&receivers)
	if err != nil {
		beego.Error("queryReceiver fail")
	}
	beego.Error("receivers:", receivers)

	redisAddr := beego.AppConfig.String("redis::redis_host") + ":" + beego.AppConfig.String("redis::redis_port")
	conn, err := redis.Dial("tcp", redisAddr)
	if err != nil {
		beego.Error("HandleOrdering:: fail connect redis:", err)
		return
	}
	//商品列表
	var goodsSKUs []map[string]interface{}
	orderTotalNum := 0
	orderTotalPrice := 0
	for _, id := range goodsSKUIds {
		var goodsSKU models.GoodsSKU
		goodsSKUId, err := strconv.Atoi(id)
		if err != nil {
			continue
		}
		goodsSKU.Id = goodsSKUId
		err = o.Read(&goodsSKU)
		if err != nil {
			beego.Error("HandleOrdering:: read data from database fail", err)
			continue
		}
		//从redis中读取 商品对应的数量h
		resp, err := conn.Do("hget", "cart_"+userName, id)
		itemGoodsNumTotal, err := redis.Int(resp, err)
		if err != nil {
			beego.Error("hget cart_"+userName, err)
			continue
		}
		//某个商品某个数量后的总价
		itemTotalPrice := goodsSKU.Price * itemGoodsNumTotal

		tmp := make(map[string]interface{})

		tmp["goodsSKU"] = goodsSKU
		tmp["itemTotalNum"] = itemGoodsNumTotal
		tmp["itemTotalPrice"] = itemTotalPrice
		goodsSKUs = append(goodsSKUs, tmp)

		//总订单数目
		orderTotalNum++
		//总订单价格
		orderTotalPrice += itemTotalPrice
	}
	//添加成功后，删除购物车中的商品

	c.Data["username"] = userName
	c.Data["img_url"] = beego.AppConfig.String("imgURL::img_url")
	c.Data["receivers"] = receivers
	c.Data["goodsSKUs"] = goodsSKUs
	c.Data["goodsSKUIds"] = goodsSKUIds
	c.Data["orderTotalNum"] = orderTotalNum
	c.Data["orderTotalPrice"] = orderTotalPrice
	c.TplName = "place_order.html"
}

func (c *OrderController) HandleAddOrder() {
	receiverId := c.GetString("addr")
	payway := c.GetString("payway")
	goodsIdsStr := c.GetString("goodsIds")
	beego.Error("goodsStr=", goodsIdsStr)
	str := goodsIdsStr[1 : len(goodsIdsStr)-1]
	sliceStr := strings.Split(str, " ")

	beego.Error("str=", sliceStr)

	resp := make(map[string]interface{})
	if receiverId == "" || payway == "" {
		resp["errmsg"] = "params empty"
		resp["errno"] = 1
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	//用户信息h
	userNameS := c.GetSession("username")
	if userNameS == nil {
		resp["errmsg"] = "can not get user info"
		resp["errno"] = 2
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}
	var user models.User
	userName := userNameS.(string)
	user.Name = userName
	errN := o.Read(&user, "Name")
	if errN != nil {
		beego.Error("userName=", userName)
		resp["errmsg"] = "fail Read user info from db "
		resp["errno"] = 2
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}
	//Receiver
	var receiver models.Receiver
	receiver.Id, _ = strconv.Atoi(receiverId)
	o.Read(&receiver)

	//插入数据到订单表
	var order models.OrderInfo
	order.Receiver = &receiver
	order.PayWay = payway
	order.OrderId = time.Now().Format("20060102150405") + strconv.Itoa(user.Id)
	order.User = &user
	_, inoErr := o.Insert(&order)
	if inoErr != nil {
		beego.Error("HandleAddOrder:: insert orderinfo data to db fail", inoErr)
		//c.Redirect("/cart/list", 302)
		resp["errmsg"] = "order fail"
		resp["errno"] = 2
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}

	//从redis中获取商品数量
	redisAddr := beego.AppConfig.String("redis::redis_host") + ":" + beego.AppConfig.String("redis::redis_port")
	conn, err := redis.Dial("tcp", redisAddr)
	if err != nil {
		resp["errmsg"] = "connect redis fail"
		resp["errno"] = 2
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}

	//插入数据到订单商品表
	//orderTotal
	for _, v := range sliceStr {
		goodsSKUId, _ := strconv.Atoi(v)
		//GoodsSKU
		var goodsSKU models.GoodsSKU
		goodsSKU.Id = goodsSKUId
		o.Read(&goodsSKU)
		//OrderGoods
		resp, err := conn.Do("hget", "cart_"+userName, goodsSKUId)
		goodCount, err := redis.Int(resp, err)

		var orderGoods models.OrderGoods
		orderGoods.Count = goodCount

		orderGoods.Price = goodsSKU.Price * goodCount
		orderGoods.OrderInfo = &order
		orderGoods.GoodsSKU = &goodsSKU
		_, iErr := o.Insert(&orderGoods)
		if iErr != nil {
			beego.Error("HandleAddOrder:: insert orderGoods data to db fail", goodsIdsStr, "errinfo:", iErr)
		}
	}
	//更新订单总价
	//order.TotalPrice = orderTotalPrice
	//o.Update(&order)

	resp["errmsg"] = "add order success"
	resp["errno"] = 0
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *OrderController) ShowPayOrder() {

	var privateKey = `MIIEpQIBAAKCAQEAszC3Yqadn9PbvjHX8r7ZaVVSBUosOIGYK0oOj5ET2hYX4AZQ
	9Rc///99yo68hsjDD90NzM3DuOqfkE8qpZsoKo2WGn9bTCOfMBXX3n/Tx9VaP01n
	QeAMzfi/ib55cK94a8PY2FvYorwQW6V/l459HDbVkb+ctZZB0YKwbTASdQ5VQs1z
	EE2i/PGVkmSAuj776POnVbdGU4LN9Dz1q/qWoicrriSkzINlA9GqQG6vtDXvVabs
	5nr2SU8S3vRNIuiA+FOtnpDDnAgpF0P0ZZm4MdR9WTtbR45t4NzuZH4Q9COiuBP4
	0fbM4Zo5E7IRvCN5EnGh9YaBY2w1AAR5UCOmDQIDAQABAoIBAQCBeMy4efzgM9rN
	lQQcgCtlAWHvMoW7GmBRvwPAVioK5PXBR68NOAxlMzy3s+SiWsMeXjGPbolhvh0m
	zxzYZcBi5sSzRpw36nEl9FJykNf7xrubi5j1LybxWC9FHpxugEq5StwOkGZ6Rvpm
	zbDgV/MsBK7Rzao0RmouMIi7jAV6D3d9+jEU0nhkggPnU4UAFB2mgeEGCbpqMSMq
	Rkl+f2U7eO/Zr7qO78IaHZtIqCFrsuVAgpMV8EHYE5qrH261NhEJvOEuhcGHzQjs
	mt8BPjRRar1DF4IkinX8tNz7sSI8MUUx7TXiu/ZNf42XIoXbomAjwLsYRTEu/3QB
	K9EGaFhBAoGBAOmKgwDGt7GLTJS1lz39RvIhYgLOA0OdJRqIWejvVdtKsPoDVLWQ
	G4GSH7mPlSxdWiJpqAy/gLjMoOKx2Ltivtzomi4Pzu8Dl+alkQIcINTAD7hU1Lts
	M1wwv0H4JbQeY+8XSdldZkywSz1TfRYBLV0gbhPoDGikR/S0cwTu2D3dAoGBAMRs
	KC0aqu3KGUGUDbpY4hnuEEp9fgAKyUxyfahbC746scIBDeYr66t9lx4r+qh28VVq
	Xzom58YVc7HGdEJQVMUs7IiE2nM0gHPennFb9dbKYEJZoucAkv5RDdRiBMh546f/
	juCWPZMk7zvPhEgUz6wZ4qXtKvw3knpuHWNDRv3xAoGBAIs6v5UT42mekV5K0Ert
	l4E8s7DWXw3NPtSNm4SKQxZEdjPnDnZb3nolwnIfDqDvWpAPi1dmR/hkTjo4Kuy0
	FvOeXGS/me/WpZWk+UlXuZ78jaKoOFcwT4JTsYJDzT6Pq6ZbrPRAgX+QzppWDsmy
	k/fkIJwPiG5OGnPhrHyxZAulAoGAQ1elipk7Ax2n/RDKiBoTIrq5ASD3QwJvs3MJ
	W+AjLYwoB5Ce+EwUl44Ocny3imyFHzjB/0j5a7NNICUfFOE/vv5A9ik+UAMvCwrH
	Haxeo85spDLhI/vRabnWWPtmEmmfwKhgjuVTpRAjqUjjXGcuMB4L08F1XFWdNbZt
	Auw8+bECgYEAqWSB0V/bjLv/d7oAtyprrtUlBio2J1djJbLiG4UiQOIDvqo0wa6y
	NuuzsLhD9wYsNvxRWzZD5YnFCEzZGI6nRp8ux7fAgxx0vegRmW6jEZeSudqW+Cp4
	nmLOV5eg5GYUiR/HqZmxiPX7s7tUmNyQOgKE/0VwRfisbs5nh9dnTZE=`
	var appId = "2016092500591329"
	var aliPublicKey = `MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAszC3Yqadn9PbvjHX8r7Z
	aVVSBUosOIGYK0oOj5ET2hYX4AZQ9Rc///99yo68hsjDD90NzM3DuOqfkE8qpZso
	Ko2WGn9bTCOfMBXX3n/Tx9VaP01nQeAMzfi/ib55cK94a8PY2FvYorwQW6V/l459
	HDbVkb+ctZZB0YKwbTASdQ5VQs1zEE2i/PGVkmSAuj776POnVbdGU4LN9Dz1q/qW
	oicrriSkzINlA9GqQG6vtDXvVabs5nr2SU8S3vRNIuiA+FOtnpDDnAgpF0P0ZZm4
	MdR9WTtbR45t4NzuZH4Q9COiuBP40fbM4Zo5E7IRvCN5EnGh9YaBY2w1AAR5UCOm
	DQIDAQAB`

	var client = alipay.New(appId, aliPublicKey, privateKey, false)

	var p = alipay.AliPayTradePagePay{}
	p.NotifyURL = "http://172.20.10.2:8080/payOk"
	p.ReturnURL = "http://172.20.10.2:8080/payOk"
	p.Subject = "you are mall"
	p.OutTradeNo = "123452123455"
	p.TotalAmount = "1000.21"
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	var url, err = client.TradePagePay(p)
	if err != nil {
		fmt.Println(err)
	}

	var payUrl = url.String()
	c.Redirect(payUrl, 302)
}

func (c *OrderController) PayOk() {
	trade_no := c.GetString("trade_no")
	if trade_no != "" {

	}
	beego.Info("payok invoke..ing")
	c.Redirect("/ucenter/orderlist", 302)
}

//短信SDK 接入
