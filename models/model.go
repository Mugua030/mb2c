package models

import (

	//"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Uid      int    `orm:"pk;auto"`
	UserName string `orm:"unique;size(100)"`
	Pwd      string `orm:"size(80)"`
	Email    string
	Power    int `orm:"default(0);comment(权限)"`
	Active   int `orm:"default(0);comment(是否激活)"`

	Receiver []*Receiver `orm:"reverse(many)"`
}

type Receiver struct {
	Rid     int `orm:"pk;auto"`
	Name    string
	ZipCode string
	Addr    string
	Phone   string
	Default bool `orm:"default(false)"`

	User *User `orm:"rel(fk)"`
}

func init() {
	//iniconf, err := NewConfig("ini", "app.conf")
	//beego.AppConfig.String()
	orm.RegisterDataBase("default", "mysql", "zhouping:telnetdb@tcp(localhost:3306)/mb2c?charset=utf8")
	orm.RegisterModel(new(User), new(Receiver))
	orm.RunSyncdb("default", false, true)
}
