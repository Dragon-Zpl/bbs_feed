package feed_permission

import (
	"bbs_feed/boot"
	"github.com/astaxie/beego/orm"
)

const (
	tablename = "feed_permission"
)

func init() {
	orm.RegisterModelWithPrefix("pre_", new(Model))
}

type Model struct {
	TopicId           int    `orm:"pk;column(topic_id)"`
	HotThread         int    `orm:"column(hot_thread)" json:"hotThread"`
	Essence           int    `orm:"column(essence)" json:"essence"`
	TodayIntroduction int    `orm:"column(today_introduction)" json:"todayIntroduction"`
	WeekPopularity    int    `orm:"column(week_popularity)" json:"weekPopularity"`
	WeekContribution  int    `orm:"column(week_contribution)" json:"weekContribution"`
	TopicIds          string `orm:"column(topicIds)"`
	IsUse             int    `orm:"column(isUse)"`
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

func GetOne(topicId string) (m Model, err error) {
	o := boot.GetSlaveMySql()
	qs := o.QueryTable((*Model)(nil))
	err = qs.Filter("topicId", topicId).One(&m)
	return
}
