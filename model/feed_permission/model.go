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
	TopicId           int    `orm:"pk;column(topic_id)" json:"topicId"`
	Fid               int    `orm:"column(fid)" json:"fid"`
	Hot               int    `orm:"column(hot)" json:"hot"`
	NewHot            int    `orm:"column(newHot)" json:"newHot"`
	Essence           int    `orm:"column(essence)" json:"essence"`
	TodayIntroduction int    `orm:"column(todayIntroduction)" json:"todayIntroduction"`
	WeekPopularity    int    `orm:"column(weekPopularity)" json:"weekPopularity"`
	WeekContribution  int    `orm:"column(weekContribution)" json:"weekContribution"`
	TopicIds          string `orm:"column(topic_ids)" json:"topicIds"`
	IsUse             int    `orm:"column(is_use)" json:"isUse"`
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

//1代表使用，0代表未使用
func UpdateIsUse(topicId int, isUse int) (err error) {
	o := boot.GetMasterMysql()
	qs := o.QueryTable((*Model)(nil))
	_, err = qs.Filter("topicId", topicId).Update(orm.Params{
		"is_use": isUse,
	})
	return
}

func UpdateFeedType(topicId int, feedTyp string, isUse int) (err error) {
	o := boot.GetMasterMysql()
	qs := o.QueryTable((*Model)(nil))
	_, err = qs.Filter("topicId", topicId).Update(orm.Params{
		feedTyp: isUse,
	})
	return
}

func UpdateTopicIds(topicId string, topicIds string) (err error) {
	o := boot.GetMasterMysql()
	qs := o.QueryTable((*Model)(nil))
	_, err = qs.Filter("topicId", topicId).Update(orm.Params{
		"topic_ids": topicIds,
	})
	return
}

func Insert(m Model) error {
	o := boot.GetMasterMysql()
	_, err := o.Insert(&m)
	return err
}
