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
	TopicId        int    `orm:"pk;column(topicId)"`
	Hot            string `orm:"column(hot)"`
	Essence        string `orm:"column(essence)"`
	WeekPopularity string `orm:"column(week_popularity)"`
	Contribution   string `orm:"column(contribution)"`
	TopicIds       string `orm:"column(topicIds)"`
	IsUse          int    `orm:"column(is_use)"`
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
