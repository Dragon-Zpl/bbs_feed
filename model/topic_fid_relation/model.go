package topic_fid_relation

import (
	"bbs_feed/boot"
	"github.com/astaxie/beego/orm"
)

const (
	tablename = "topic_fid_relation"
)

func init() {
	orm.RegisterModelWithPrefix("pre_", new(Model))
}

type Model struct {
	Id      int `orm:"pk;column(id)"`
	TopicId int `orm:"column(topic_id)"`
	Fid     int `orm:"column(fid)"`
}

// 实现表名的接口
func (m *Model) TableName() string {
	return tablename
}

func GetFids(topicIds []string) []int {
	o := boot.GetMasterMysql()
	qs := o.QueryTable((*Model)(nil))
	ms := make([]*Model, 0)
	res := make([]int, 0)
	qs.Filter("topic_id__in", topicIds).All(&ms)
	for _, m := range ms {
		res = append(res, m.Fid)
	}
	return res
}
