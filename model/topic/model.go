package topic

import (
	"bbs_feed/boot"
	"github.com/astaxie/beego/orm"
)

const (
	tablename = "topic"
)

func init() {
	// 需要在init中注册定义的model
	orm.RegisterModelWithPrefix("pre_", new(Model))
}

// 定义表中的数据模型
type Model struct {
	Id           int    `orm:"pk;column(id)"`
	Title        string `orm:"size(20);column(title)"`
	Introduction string `orm:"size(255);column(introduction)"`
	FollowCount  int    `orm:"default(0);column(follow_count)"`
	ThreadCount  int    `orm:"null;default(0);column(thread_count)"`
	IsUse        string `orm:"size(10);default('yes');column(is_use)"`
}

// 实现表名的接口
func (m *Model) TableName() string {
	return "topic"
}

func GetAll() []*Model {
	o := boot.GetSlaveMySql()
	qs := o.QueryTable((*Model)(nil))
	ms := make([]*Model, 0)
	qs.Filter("is_use", 1).All(&ms)
	return ms
}
