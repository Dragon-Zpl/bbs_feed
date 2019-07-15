package feed_conf

import (
	"bbs_feed/boot"
	"github.com/astaxie/beego/orm"
)

const (
	tablename = "feed_conf"
)

func init() {
	orm.RegisterModelWithPrefix("pre_", new(Model))
}

type Model struct {
	Id    int    `orm:"pk;column(id)" json:"id"`
	Name  string `orm:"column(name)" json:"name"`
	Conf  string `orm:"column(conf)" json:"conf"`
	IsUse int    `orm:"column(is_use)" json:"is_use"`
}

// 实现表名的接口
func (m *Model) TableName() string {
	return tablename
}

func GetAll() []*Model {
	o := boot.GetSlaveMySql()
	qs := o.QueryTable((*Model)(nil))
	m := make([]*Model, 0)
	qs.Filter("is_use", 1).All(&m)
	return m
}
