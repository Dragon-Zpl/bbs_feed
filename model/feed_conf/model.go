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
	IsUse int    `orm:"column(is_use)" json:"isUse"`
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

func GetOne(name string) (m Model, err error) {
	o := boot.GetSlaveMySql()
	qs := o.QueryTable((*Model)(nil))
	err = qs.Filter("name", name).One(&m)
	return
}

func Insert(m Model) (err error) {
	o := boot.GetMasterMysql()
	_, err = o.Insert(&m)
	return err
}

func UpdateConf(typ string, conf string) (err error) {
	o := boot.GetMasterMysql()
	qs := o.QueryTable((*Model)(nil))
	_, err = qs.Filter("name", typ).Update(orm.Params{
		"conf": conf,
	})
	return
}
