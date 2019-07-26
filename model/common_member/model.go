package common_member

import (
	"bbs_feed/boot"
	"github.com/astaxie/beego/orm"
)

const tablename = "pre_common_member"

func init() {
	orm.RegisterModelWithPrefix("", new(Model))
}

type Model struct {
	Uid      int    `orm:"pk;column(uid)" json:"uid"`
	Username string `orm:"column(username)" json:"username"`
}

func (m *Model) TableName() string {
	return tablename
}

func GetUserName(uid string) (m Model, err error) {
	o := boot.GetSlaveMySql()
	qs := o.QueryTable((*Model)(nil))
	err = qs.Filter("uid", uid).One(&m)
	return
}
