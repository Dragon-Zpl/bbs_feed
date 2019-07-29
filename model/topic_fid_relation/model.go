package topic_fid_relation

import (
	"bbs_feed/boot"
	"github.com/astaxie/beego/logs"
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
	o := boot.GetSlaveMySql()
	qs := o.QueryTable((*Model)(nil))
	ms := make([]*Model, 0)
	res := make([]int, 0)
	qs.Filter("topic_id__in", topicIds).All(&ms)
	for _, m := range ms {
		res = append(res, m.Fid)
	}
	return res
}

func GetAllFids() (map[string]string, error) {
	var maps []orm.Params
	o := boot.GetSlaveMySql()
	sql := "select pref.fid, pref.topic_id from pre_topic_fid_relation pref, pre_topic pret where pref.topic_id = pret.id and pret.is_use = 'yes'"
	_, err := o.Raw(sql).Values(&maps)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	results := make(map[string]string)
	for _, v := range maps {
		results[v["topic_id"].(string)] = v["fid"].(string)
	}
	return results, nil

}
