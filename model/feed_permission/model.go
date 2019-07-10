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
	HotThread         int    `orm:"column(hot_thread)" json:"hot_thread"`
	Essence           int    `orm:"column(essence)" json:"essence"`
	TodayIntroduction int    `orm:"column(today_introduction)" json:"today_introduction"`
	WeekPopularity    int    `orm:"column(week_popularity)" json:"week_popularity"`
	WeekContribution  int    `orm:"column(week_contribution)" json:"week_contribution"`
	TopicIds          string `orm:"column(topic_ids)"`
	IsUse             int    `orm:"column(is_use)"`
}

// 实现表名的接口
func (m *Model) TableName() string {
	return tablename
}

func GetAll() []*Model {
	o := boot.GetMasterMysql()
	qs := o.QueryTable((*Model)(nil))
	m := make([]*Model, 0)
	qs.Filter("is_use", 1).All(&m)
	return m
}

func GetOne(topicId string) (m Model, err error) {
	o := boot.GetMasterMysql()
	qs := o.QueryTable((*Model)(nil))
	err = qs.Filter("topicId", topicId).One(&m)
	return
}
