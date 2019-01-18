package controllers

import (
	"mb2c/models"
	"path"

	"github.com/astaxie/beego/orm"

	"github.com/astaxie/beego"
	"github.com/weilaihui/fdfs_client"
)

type AdminGoodsController struct {
	beego.Controller
}

//GetGoodsTypes 获取商品分类
func (c *AdminGoodsController) GetGoodsTypes(o orm.Ormer) []models.GoodsType {

	var goodsTypes []models.GoodsType
	_, tpsErr := o.QueryTable("GoodsType").All(&goodsTypes)
	if tpsErr != nil {
		beego.Error("ShowGoodsList:: Get goodsTypeList fail", tpsErr)
	}

	return goodsTypes
}

//GetGoodsSPU  获取商品
func (c *AdminGoodsController) GetGoodsSPUList(o orm.Ormer) []models.Goods {
	var goodsSPUList []models.Goods
	o.QueryTable("Goods").All(&goodsSPUList)

	return goodsSPUList
}

func (c *AdminGoodsController) ShowGoodsList() {

	//分类列表
	o := orm.NewOrm()
	var goodsTypes []models.GoodsType
	_, tpsErr := o.QueryTable("GoodsType").All(&goodsTypes)
	if tpsErr != nil {
		beego.Error("ShowGoodsList:: Get goodsTypeList fail", tpsErr)
		return
	}

	typeId, tErr := c.GetInt("select")
	if tErr != nil {
		typeId = 0 // so ..all types
	}
	beego.Error("typeId=", typeId)
	page, pErr := c.GetInt("page")
	if pErr != nil {
		page = 1
	}
	pageSize := 3
	start := (page - 1) * pageSize

	//商品列表  : 如果没有选择分类，那么就列出所有的分类商品
	var goodsSKUs []models.GoodsSKU
	qs := o.QueryTable("GoodsSKU").RelatedSel("GoodsType")
	if typeId != 0 {
		qs = qs.Filter("GoodsType__Id", typeId)
	}
	_, oErr := qs.Limit(pageSize, start).All(&goodsSKUs)
	if oErr != nil {
		beego.Error("ShowGoodsList:: query goodssku fail", oErr)
		return
	}

	c.Data["imgUrl"] = beego.AppConfig.String("imgURL::img_url")
	c.Data["typeId"] = typeId
	c.Data["goodsTypes"] = goodsTypes
	c.Data["goodsSKUs"] = goodsSKUs

	c.Layout = "admin_layout.html"
	c.TplName = "admin_index.html"
}

func (c *AdminGoodsController) ShowAddGoods() {
	o := orm.NewOrm()
	goodsTypes := c.GetGoodsTypes(o)
	goodsSPUList := c.GetGoodsSPUList(o)

	beego.Error("goodTypes=", goodsTypes)

	c.Data["goodsTypes"] = goodsTypes
	c.Data["goodsSPUList"] = goodsSPUList
	c.Layout = "admin_layout.html"
	c.TplName = "admin_add.html"
}
func (c *AdminGoodsController) HandleAddGoods() {
	goodsId, gerr := c.GetInt("selectGoodsSPU")
	if gerr != nil {
		beego.Error("HandleAddGoods:: no goods id")
		c.Data["errmsg"] = "empty goods id"
		c.TplName = "admin_add.html"
		return
	}
	goodsTypeId, terr := c.GetInt("selectType")
	if terr != nil {
		//beego.Error("HandleAddGoods:: empty goods_typeid")

		c.Data["errmsg"] = "empty goodsTypeid"
		c.TplName = "admin_add.html"
		return
	}
	goodsName := c.GetString("goodsName")
	if goodsName == "" {
		c.Data["errmsg"] = "empty goodsName"
		c.TplName = "admin_add.html"
		return
	}
	//库存
	stock, serr := c.GetInt("goodsStock")
	if serr != nil {
		c.Data["errmsg"] = "empty goodsStock"
		c.TplName = "admin_add.html"
		return
	}
	//商品价格
	price, err := c.GetInt("goodsPrice")
	if err != nil {
		c.Data["errmsg"] = "empty goods price"
		c.TplName = "admin_add.html"
		return
	}
	//图片
	file, head, err := c.GetFile("uploadname")
	defer file.Close()
	fileExtName := path.Ext(head.Filename)

	fdfsClient, err := fdfs_client.NewFdfsClient("/Users/Mugua/go/src/mb2c/conf/client.conf")
	if err != nil {
		beego.Error("upload file to fdfs fail:: ", err)
		c.Data["errmsg"] = "upload file fail"
		c.TplName = "admin_add.html"
		return
	}
	var buffer []byte
	_, err = file.Read(buffer)
	if err != nil {
		c.Data["errmsg"] = "file read to buffer fail"
		c.TplName = "admin_add.html"
		return
	}
	resp, err := fdfsClient.UploadByBuffer(buffer, fileExtName[1:])
	if err != nil {
		beego.Error("upfile to fdfs fail: ", err)
		c.Data["errmsg"] = "upfile to fdfs ail"
		c.TplName = "admin_add.html"
		return
	}
	//beego.Info("stock=", stock, "price=", price, "goodTypeId=", goodsTypeId)
	//beego.Error("uploadResult: ", resp.RemoteFileId)
	o := orm.NewOrm()
	var goodsSKU models.GoodsSKU
	var goods models.Goods
	var goodsType models.GoodsType
	//商品
	goods.Id = goodsId
	o.Read(&goods)
	goodsType.Id = goodsTypeId
	o.Read(&goodsType)

	goodsSKU.Name = goodsName
	goodsSKU.Stock = stock
	goodsSKU.Price = price
	goodsSKU.GoodsType = &goodsType
	goodsSKU.Goods = &goods
	goodsSKU.Image = resp.RemoteFileId
	_, err = o.Insert(&goodsSKU)
	if err != nil {
		c.Data["errmsg"] = "insert data to data fail"
		c.TplName = "admin_add.html"
		return
	}

	c.Redirect("/admin/goods/list", 302)
}

//ShowAddGoodsType 展示添加商品类型
func (c *AdminGoodsController) ShowAddGoodsType() {

	c.Layout = "admin_layout.html"
	c.TplName = "admin_addType.html"
}

//HandleAddGoodsType 添加商品
func (c *AdminGoodsController) HandleAddGoodsType() {

}

//ShowAddGoodsSKU 展示添加商品SKU
func (c *AdminGoodsController) ShowAddGoodsSPU() {
	c.Layout = "admin_layout.html"
	c.TplName = "admin_addGoodsSPU.html"

}

//HandleAddGoodsSKU  添加商品SKU
func (c *AdminGoodsController) HandleAddGoodsSPU() {

}
