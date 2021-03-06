package controllers

import (
	"math"
	"mb2c/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
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
	//全部商品分类
	o := orm.NewOrm()
	var goodsTypes []*models.GoodsType
	_, err := o.QueryTable("GoodsType").All(&goodsTypes)
	if err != nil {
		beego.Error("Get goodsType fail")
		return
	}
	//轮播广告 advertisement
	var advsOnIndex []*models.IndexGoodsBanner
	_, err = o.QueryTable("IndexGoodsBanner").All(&advsOnIndex)
	if err != nil {
		beego.Error("GetAdvsOnIndex fail")
		return
	}
	//促销广告h
	var advsOfPromotion []*models.IndexPromotionBanner
	_, err = o.QueryTable("IndexPromotionBanner").All(&advsOfPromotion)
	if err != nil {
		beego.Error("GetAdvsOfPromotion fail", err)
		return
	}
	//商品列表
	var goodsList []map[string]interface{}
	for _, v := range goodsTypes {
		tmp := make(map[string]interface{})
		tmp["goodsType"] = v
		goodsList = append(goodsList, tmp)
	}

	for _, item := range goodsList {
		//此分类下的
		qs := o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType", "GoodsSKU").Filter("GoodsType", item["goodsType"])

		//文字广告h
		var textTypeBanner []*models.IndexTypeGoodsBanner
		var imgTypeBanner []*models.IndexTypeGoodsBanner
		qs.Filter("DisplayType", 0).All(&textTypeBanner)
		//图片广告
		qs.Filter("DisplayType", 1).All(&imgTypeBanner)
		item["goodsTextBanner"] = textTypeBanner
		item["goodsImgBanner"] = imgTypeBanner
	}

	c.Data["goodsTypes"] = goodsTypes
	c.Data["advsOnIndex"] = advsOnIndex
	c.Data["advsOfPromotion"] = advsOfPromotion
	c.Data["goodsList"] = goodsList

	//c.TplName = "index.html"
	//	c.Layout = "layout.html"
	c.TplName = "index.html"
}
func (c *GoodsController) ShowGoodsDetail() {

	skuId, err := c.GetInt("sku_id")
	if err != nil {
		beego.Error("ShowGoodsDetail:: Get sku_id fail", err)
		return
	}
	//进入详情就意味着有浏览记录了
	redisHost := beego.AppConfig.String("redis::redis_host")
	redisPort := beego.AppConfig.String("redis::redis_port")
	raddr := redisHost + ":" + redisPort
	conn, err := redis.Dial("tcp", raddr)

	if err == nil {
		defer conn.Close()
		userName := c.GetSession("username")
		if userName != nil {
			conn.Do("lrem", "his_"+userName.(string), skuId)
			re, err := conn.Do("lpush", "his_"+userName.(string), skuId)

			beego.Error("write success ho!!", re)
			beego.Error("err", err)
		}
	} else {
		beego.Error("connect Redis fail:", err)
	}

	var goodsSKU models.GoodsSKU
	o := orm.NewOrm()
	goodsSKU.Id = skuId
	err = o.Read(&goodsSKU)
	if err != nil {
		beego.Error("ShowGoodsDetail:: no this sku info")
		return
	}
	//新品推荐
	var goodsSKUs []*models.GoodsSKU
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType", goodsSKU.GoodsType).OrderBy("Time").Limit(2, 0).All(&goodsSKUs)

	//beego.Info("skuinfo:", goodsSKU.Image)
	c.Data["goodsSKUInfo"] = goodsSKU
	c.Data["goodsSKUs"] = goodsSKUs

	c.Layout = "layout.html"
	c.TplName = "detail.html"
}

//GetNewGoodsByRecommend  新品推荐
func (c *GoodsController) GetNewGoodsByRecommend(typeId int) (goodsSKUs []*models.GoodsSKU) {
	//var goodsSKUs []*models.GoodsSKU
	o := orm.NewOrm()
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", typeId).OrderBy("-Time").Limit(2, 0).All(&goodsSKUs)

	return goodsSKUs
}

//ShowGoodsListByType 分类下的商品列表
func (c *GoodsController) ShowGoodsListByType() {
	typeId, err := c.GetInt("type_id")
	if err != nil {
		beego.Error("ShowGoodsListByType:: fail", err)
		return
	}
	curPage, err := c.GetInt("page")
	if err != nil {
		curPage = 1
	}

	var goodsType models.GoodsType
	var goodsSKUs []*models.GoodsSKU
	o := orm.NewOrm()
	goodsType.Id = typeId
	o.Read(&goodsType)

	pageSize := 2
	start := (curPage - 1) * pageSize

	qs := o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", goodsType.Id)
	count, err := qs.Count()
	qs.Limit(pageSize, start).All(&goodsSKUs)
	pageCount := math.Ceil(float64(count) / float64(pageSize))

	c.Data["img_url"] = beego.AppConfig.String("imgURL::img_url")
	c.Data["goodsSKUs"] = goodsSKUs
	c.Data["newGoodsRecommend"] = c.GetNewGoodsByRecommend(typeId)
	c.Data["typeId"] = typeId

	//分页 一个页数的切片
	pages := PageDecorate(int(pageCount), int(curPage))
	c.Data["curPage"] = curPage
	prePage := curPage - 1
	nextPage := curPage + 1

	c.Data["prePage"] = prePage
	c.Data["nextPage"] = nextPage
	c.Data["pages"] = pages
	beego.Error("typeId=", typeId)

	c.TplName = "list.html"
}

//PageDecorate 分页装饰列表
func PageDecorate(pageCount int, curPage int) []int {
	curPageLen := 5
	var curPageList []int
	beego.Info("pageCount=", pageCount)
	beego.Info("curPage=", curPage)

	//处在第几个片中（一个片是5个）
	inSliceNo := math.Ceil(float64(curPage) / float64(curPageLen))
	curCeil := int(inSliceNo) * curPageLen
	start := curCeil - 5
	if start <= 0 {
		start = 1
	}

	for i := start; i <= curPageLen; i++ {
		if i <= pageCount {
			curPageList = append(curPageList, i)
		}
	}

	return curPageList
}
