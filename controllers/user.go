package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"mb2c/models"
	"regexp"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils"
	"github.com/gomodule/redigo/redis"
)

// UserController for insert user
type UserController struct {
	beego.Controller
}

//ShowRegister show index page
func (c *UserController) ShowRegister() {
	//sum := sha256.Sum256([]byte("helloworda"))

	c.TplName = "register.html"
}

func (c *UserController) HandleRegister() {
	userName := c.GetString("user_name")
	if userName == "" {
		c.Data["err_tips"] = "用户名不能为空"
		c.TplName = "register.html"
		return
	}
	pwd := c.GetString("pwd")
	if pwd == "" {
		c.Data["err_tips"] = "密码不能空"
		c.TplName = "register.html"
		return
	}
	cpwd := c.GetString("cpwd")
	if cpwd == "" {
		c.Data["err_tips"] = "请再次填入密码"
		c.TplName = "register.html"
		return
	}
	if pwd != cpwd {
		c.Data["err_tips"] = "两次填写的密码不一致"
		c.TplName = "register.html"
		return
	}
	email := c.GetString("email")
	if email == "" {
		c.Data["err_tips"] = "邮箱不能为空"
		c.TplName = "register.html"
		return
	}
	re := regexp.MustCompile(`[a-zA-Z0-9]*\.[a-zA-Z0-9]+@[a-zA-Z0-9]+\.[a-zA-Z]?`)
	matched := re.MatchString(email)
	if !matched {
		c.Data["err_tips"] = "邮箱格式不正确"
		c.TplName = "register.html"
		return
	}

	var user models.User
	// 密码加密
	hash := sha256.New()
	hash.Write([]byte(pwd))
	md := hash.Sum(nil)
	pwdStr := hex.EncodeToString(md)
	//beego.Error("sha256v : ", mdStr)

	user.Name = userName
	user.Password = pwdStr
	user.Email = email
	o := orm.NewOrm()
	uid, err := o.Insert(&user)

	if err != nil {
		c.Data["err_tips"] = "注册失败"
		c.TplName = "register.html"
		return
	}
	// 发送激活邮件 send email for active
	eConfig := `{"username":"540992526@qq.com","password":"biadalkoghihbbda", "host":"smtp.qq.com", "port":587}`
	emailSend := utils.NewEMail(eConfig)
	emailSend.From = "540992526@qq.com"
	emailSend.To = []string{email}
	emailSend.Subject = "beego-go-web test"
	emailSend.HTML = "<a href='http://localhost:8080/active?uid='" + strconv.FormatInt(uid, 10) + ">激活</a>"
	beego.Info("emailContent:", emailSend.HTML)
	err = emailSend.Send()
	if err != nil {
		beego.Error("sendMailFail: ", err)
	}

	c.Ctx.WriteString("注册成功")
}

//HandleActiveUser 激活用户
func (c *UserController) HandleActiveUser() {
	uid, err := c.GetInt("uid")
	if err != nil {
		beego.Error("active user fail: ", err)
		c.Ctx.WriteString("激活用户失败")
		return
	}
	beego.Info("GetUid:: ", uid)
	var user models.User
	user.Id = uid
	user.Active = true
	o := orm.NewOrm()
	_, err = o.Update(&user)
	if err != nil {
		beego.Error("active user fail: ", err)
		c.Ctx.WriteString("激活用户失败")
		return
	}
	c.Ctx.WriteString("激活用户成功")
	//c.Redirect("/login", 302)
}

//ShowLogin Page For Login
func (c *UserController) ShowLogin() {
	userName := c.Ctx.GetCookie("username")
	if userName == "" {
		c.Data["checked"] = ""
		c.Data["username"] = ""
	} else {
		c.Data["checked"] = "checked"
		c.Data["username"] = userName
	}

	c.TplName = "login.html"
}

//HandleLogin action to login
func (c *UserController) HandleLogin() {
	userName := c.GetString("username")
	if userName == "" {
		beego.Error("HandleLogin:: empty username")
		c.TplName = "login.html"
		return
	}
	pwd := c.GetString("pwd")
	if pwd == "" {
		beego.Error("HandleLogin:: empty pwd")
		c.TplName = "login.html"
		return
	}
	hash := sha256.New()
	hash.Write([]byte(pwd))
	md := hash.Sum(nil)
	pwdStr := hex.EncodeToString(md)

	var user models.User
	user.Name = userName
	o := orm.NewOrm()
	err := o.Read(&user, "Name")
	if err != nil {
		beego.Error("HandleLogin:: read user fail", err)
		c.TplName = "login.html"
		return
	}
	if user.Password != pwdStr {
		beego.Error("HandleLogin:: pwd not equal")
		c.TplName = "login.html"
		return
	}

	c.SetSession("username", userName)
	//登录成功设置记录用户名
	remember := c.GetString("remember")
	beego.Info("remember=", remember)
	if remember == "on" {
		c.Ctx.SetCookie("username", userName, 3600)
	} else {
		c.Ctx.SetCookie("username", userName, -1)
	}

	c.Redirect("/index", 302)
}

// HandleLogout  退出登录
func (c *UserController) HandleLogout() {
	c.DelSession("username")
	c.Redirect("/login", 302)
}

// ShowUserInfo 用户中心:: member info
func (c *UserController) ShowUserInfo() {
	userName := c.GetSession("username")

	var goodsSKUs []models.GoodsSKU
	if userName == nil {
		c.Data["username"] = ""
	} else {
		c.Data["username"] = userName.(string)
		//浏览历史
		redisHost := beego.AppConfig.String("redis::redis_host")
		redisPort := beego.AppConfig.String("redis::redis_port")
		conn, err := redis.Dial("tcp", redisHost+":"+redisPort+"")

		if err == nil {

			defer conn.Close()
			resp, err := conn.Do("lrange", "his_"+userName.(string), 0, 4)
			res, err := redis.Ints(resp, err)
			beego.Error(res)
			o := orm.NewOrm()
			for _, skuId := range res {

				var goodsSKU models.GoodsSKU
				goodsSKU.Id = skuId
				err = o.Read(&goodsSKU)
				if err != nil {
					continue
				}
				goodsSKUs = append(goodsSKUs, goodsSKU)
			}
		}

	}
	c.Data["goodsSKUs"] = goodsSKUs
	c.Data["img_url"] = beego.AppConfig.String("imgURL::img_url")
	beego.Info("imgurl:", c.Data["img_url"])

	c.Layout = "layout.html"
	c.TplName = "user_center_info.html"
}

//ShowUserOrderList 用户中心:: order list
func (c *UserController) ShowUserOrderList() {
	userNameS := c.GetSession("username")
	if userNameS == nil {
		c.Redirect("/login", 302)
	}

	page := 1
	pageSize := 2
	start := (page - 1) * pageSize
	userName := userNameS.(string)
	var orderlist []models.OrderInfo
	o := orm.NewOrm()

	var itemOrderList []map[string]interface{}
	o.QueryTable("OrderInfo").RelatedSel("User").Filter("User__Name", userName).Limit(pageSize, start).All(&orderlist)
	beego.Error("orderlist:", orderlist)
	for _, v := range orderlist {
		var orderGoodsList []models.OrderGoods
		o.QueryTable("OrderGoods").RelatedSel("OrderInfo", "GoodsSKU").Filter("OrderInfo__Id", v.Id).All(&orderGoodsList)

		tmp := make(map[string]interface{})
		tmp["orderInfo"] = v
		tmp["orderGoodsList"] = orderGoodsList
		itemOrderList = append(itemOrderList, tmp)
	}

	c.Data["imgUrl"] = beego.AppConfig.String("imgURL::img_url")
	c.Data["orderList"] = itemOrderList
	c.Layout = "layout.html"
	c.TplName = "user_center_order.html"
}

//ShowUserSiteInfo  用户中心:: addr receive something
func (c *UserController) ShowUserSiteInfo() {

	c.Layout = "layout.html"
	c.TplName = "user_center_site.html"
}

func (c *UserController) HandleAddSite() {
	receiver := c.GetString("receiver")
	addr := c.GetString("addr")
	phone := c.GetString("phone")
	zipCode := c.GetString("zipCode")
	if receiver == "" || addr == "" || phone == "" || zipCode == "" {
		beego.Error("HandeAddSite:: parama err: can not empty ")
		c.Redirect("/ucenter/siteinfo", 302)
		return
	}

	o := orm.NewOrm()
	var user models.User
	//查询用户信息 member info
	userNameE := c.GetSession("username")
	userName := userNameE.(string)
	user.Name = userName
	err := o.Read(&user)
	if err != nil {
		beego.Error("HandleAddSite:: no user"+userNameE.(string), err)
		c.Redirect("/ucenter/siteinfo", 302)
		return
	}

	//查询曾经的默认地址 query the default addr  once
	var receiverMayDefault models.Receiver
	qs := o.QueryTable("receiver").RelatedSel("User").Filter("User__UserName", userName)
	err = qs.Filter("Default", true).One(&receiverMayDefault)
	if err == nil {
		receiverMayDefault.IsDefault = false
		_, err := o.Update(&receiverMayDefault)
		if err != nil {
			beego.Error("HandleAddSite:: updte receiverMayDefault fail,", err)
			c.Redirect("/ucenter/siteinfo", 302)
			return
		}
	}

	var receiverOb models.Receiver
	receiverOb.Name = userName
	receiverOb.Addr = addr
	receiverOb.Phone = phone
	receiverOb.ZipCode = zipCode
	receiverOb.User = &user
	_, err = o.Insert(&receiverOb)
	if err != nil {
		beego.Error("HandleAddSite:: insert fail, ", err)
		c.Redirect("/ucenter/siteinfo", 302)
		return
	}

	//数据校验
}
