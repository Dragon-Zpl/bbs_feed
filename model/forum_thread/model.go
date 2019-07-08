package forum_thread

import (
	"bbs_feed/boot"
	"bbs_feed/lib/helper"
	"github.com/astaxie/beego/orm"
)

const (
	tablename = "forum_thread"
)

func init() {
	orm.RegisterModelWithPrefix("pre_", new(Model))
}

type Model struct {
	Tid int `orm:"pk;default(0);column(tid)"`
	Fid int `orm:"default(0);column(fid)"`
	Author       string `orm:"size(15);default('');column(author)"`
	AuthorId     int    `orm:"default(0);column(authorid)"`
	Subject      string `orm:"size(120);default('');column(subject)"`
	Dateline     int    `orm:"default(0);column(dateline)" json:"_"`
	LastPost     int    `orm:"default(0);column(lastpost)" json:"_"`
	LastPoster   string `orm:"default('');column(lastposter)"`
	Views        int    `orm:"default(0);column(views)" json:"_"`
	Replies      int    `orm:"default(0);column(replies)" json:"_"`
	DisplayOrder int8   `orm:"default(0);column(displayorder)" json:"_"`
	Digest int8 `orm:"default(0);column(digest)" json:"_"`
	Special int8 `orm:"default(0);column(special)" json:"_"`
	Status int `orm:"default(0);column(status)" json:"_"`
	FavTimes int `orm:"default(0);column(favtimes)" json:"_"`
}

// 实现表名的接口
func (m *Model) TableName() string {
	return tablename
}


func GetByCondition (fids []int, day, views, replys int) ([]*Model){
	o := boot.GetMasterMysql()
	qs := o.QueryTable((*Model)(nil))
	ms := make([]*Model, 0)
	qs.Filter("fid__in", fids).Filter("dateline__gte", helper.PreNDayTime(day)).Filter("views__gte", views).Filter("replies__gte", replys).Limit(400).All(&ms)
	return ms
}

func GetByTids(tids []int) []*Model{
	o := boot.GetMasterMysql()
	qs := o.QueryTable((*Model)(nil))
	ms := make([]*Model, 0)
	qs.Filter("tid__in", tids).All(&ms)
	return ms
}