package report_num

import (
	"bbs_feed/boot"
	"bbs_feed/lib/helper"
	"github.com/astaxie/beego/orm"
	"time"
)

const (
	tablename = "report_num"
)

func init() {
	orm.RegisterModelWithPrefix("pre_", new(Model))
}

type Model struct {
	ThreadId   int   `orm:"pk;column(thread_id)" json:"threadId"`
	Num        int   `orm:"column(num)" json:"num"`
	CreateTime int64 `orm:"column(create_time)" json:"createTime"`
	UpdateTime int64 `orm:"column(update_time)" json:"updateTime"`
}

// 实现表名的接口
func (m *Model) TableName() string {
	return tablename
}

func GetAll(duration time.Duration) []*Model {
	o := boot.GetSlaveMySql()
	qs := o.QueryTable((*Model)(nil))
	m := make([]*Model, 0)
	qs.Filter("update_time__gte", helper.PreMinuteTime(duration+time.Minute)).All(&m)
	return m
}
