package models

import (

	//"github.com/astaxie/beego/config"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

//User 用户表
type User struct {
	Id       int    `orm:"pk;auto"`
	Name     string `orm:"unique;size(20)"`
	Password string `orm:"size(80)"`
	Email    string `orm:"size(50)"`
	Power    int    `orm:"default(0)"`
	Active   bool   `orm:"default(false)"`

	Address   []*Address   `orm:"reverse(many)"`
	OrderInfo []*OrderInfo `orm:"reverse(many)"`
}

//Address 地址表
type Address struct {
	Id        int    `orm:"pk;auto"`
	Receiver  string `orm:"size(50)"`
	ZipCode   string `orm:"size(20)"`
	Addr      string
	Phone     string
	Default   bool         `orm:"default(false)"`
	User      *User        `orm:"rel(fk)"`
	OrderInfo []*OrderInfo `orm:"reverse(many)"`
}

//Goods 商品spu表
type Goods struct {
	Id       int         `orm:"pk;auto"`
	Name     string      `orm:"size(20)"`
	Detail   string      `orm:"size(200)"`
	GoodsSKU []*GoodsSKU `orm:"reverse(many)"`
}

//GoodsType 商品类型
type GoodsType struct {
	Id                   int `orm:pk;auto`
	Name                 string
	Logo                 string
	Image                string
	GoodsSKU             []*GoodsSKU             `orm:"reverse(many)"`
	IndexTypeGoodsBanner []*IndexTypeGoodsBanner `orm:"reverse(many)"` //banner广告
}

//GoodsSKU 商品SKU
type GoodsSKU struct {
	Id                   int        `orm:"pk;auto"`
	Goods                *Goods     `orm:"rel(fk)"`
	GoodsType            *GoodsType `orm:"rel(fk)"` //商品归属种类
	Name                 string     //商品名
	Desc                 string     //简介
	Price                int
	Unite                string //商品单位
	Image                string
	Stock                int                     `orm:"default(1)"`
	Sales                int                     `orm:"default(0)"` //商品销量
	Status               int                     `orm:"default(0)"` //商品状态
	Time                 time.Time               `orm:"auto_now_add"`
	GoodsImage           []*GoodsImage           `orm:"reverse(many)"`
	IndexGoodsBanner     []*IndexGoodsBanner     `orm:"reverse(many)"`
	IndexTypeGoodsBanner []*IndexTypeGoodsBanner `orm:"reverse(many)"`
	OrderGoods           []*OrderGoods           `orm:"reverse(many)"`
}

//GoodsImage 商品图片表
type GoodsImage struct {
	Id       int
	Image    string
	GoodsSKU *GoodsSKU `orm:"rel(fk)"`
}

//IndexGoodsBanner 首页轮播商品展示表
type IndexGoodsBanner struct {
	Id       int
	GoodsSKU *GoodsSKU `orm:"rel(fk)"`
	Image    string
	Index    int `orm:"default(0)"`
}

//IndexTypeGoodsBanner 首页商品分类展示表
type IndexTypeGoodsBanner struct {
	Id          int
	GoodsType   *GoodsType `orm:"rel(fk)"`
	GoodsSKU    *GoodsSKU  `orm:"rel(fk)"`
	DisplayType int        `orm:"default(1)"`
	Index       int        `orm:"default(0)"`
}

//IndexPromotionBanner 首页促销商品展示表
type IndexPromotionBanner struct {
	Id    int
	Name  string `orm:"size(20)"`
	Url   string `orm:"size(50)"`
	Image string
	Index int `orm:"default(0)"`
}

//OrderInfo 订单信息表
type OrderInfo struct {
	Id           int
	OrderId      string        `orm:"unique"`
	User         *User         `orm:"rel(fk)"`
	Address      *Address      `orm:"rel(fk)"`
	PayWay       string        //支付方式   //wx ; zfb
	PayMethod    int           //付款方式 支付方法 银行卡 信用卡
	TotalPrice   int           //商品总价
	TransitPrice int           //运费
	TraceNo      string        `orm:"default('')"`
	Time         time.Time     `orm:"auto_now_add"`
	OrderGoods   []*OrderGoods `orm:"reverse(many)"`
	Orderstatus  int           `orm:"default(0)"` //订单状态
}

//OrderGoods 订单商品表
type OrderGoods struct {
	Id        int
	OrderInfo *OrderInfo `orm:"rel(fk)"`
	GoodsSKU  *GoodsSKU  `orm:"rel(fk)"`
	Count     int        `orm:"default(1)"`
	Price     int
	Comment   string `orm:"default('')"`
}

//Receiver 收件人
type Receiver struct {
	Id      int
	Name    string //收件人名字
	ZipCode string //收件人邮编
	Addr    string //地址
	Phone   string //收件人联系方式
	Default bool   `orm:"default(false)"` //是否未默认收件人
	User    *User  `orm:"rel(fk)"`
}

func init() {
	//iniconf, err := NewConfig("ini", "app.conf")
	//beego.AppConfig.String()
	orm.RegisterDataBase("default", "mysql", "zhouping:telnetdb@tcp(localhost:3306)/mb2c?charset=utf8")
	orm.RegisterModel(new(User), new(Address), new(OrderGoods), new(OrderInfo), new(IndexPromotionBanner), new(IndexTypeGoodsBanner), new(IndexGoodsBanner), new(GoodsImage), new(GoodsSKU), new(GoodsType), new(Goods))
	orm.RunSyncdb("default", false, true)
}
