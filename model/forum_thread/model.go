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
	Tid          int    `orm:"pk;default(0);column(tid)" json:"tid"`
	Fid          int    `orm:"default(0);column(fid)" json:"fid"`
	Author       string `orm:"size(15);default('');column(author)" json:"author"`
	AuthorId     int    `orm:"default(0);column(authorid)" json:"authorId"`
	Subject      string `orm:"size(120);default('');column(subject)" json:"subject"`
	Dateline     int    `orm:"default(0);column(dateline)" json:"dateline"`
	LastPost     int    `orm:"default(0);column(lastpost)" json:"_"`
	LastPoster   string `orm:"default('');column(lastposter)" json:"_"`
	Views        int    `orm:"default(0);column(views)" json:"_"`
	Replies      int    `orm:"default(0);column(replies)" json:"_"`
	DisplayOrder int8   `orm:"default(0);column(displayorder)" json:"_"`
	Digest       int8   `orm:"default(0);column(digest)" json:"_"`
	Special      int8   `orm:"default(0);column(special)" json:"_"`
	Status       int    `orm:"default(0);column(status)" json:"_"`
	FavTimes     int    `orm:"default(0);column(favtimes)" json:"_"`
}

// 实现表名的接口
func (m *Model) TableName() string {
	return tablename
}

func GetHotThreads(fids []int, day, views, replys int) []*Model {
	o := boot.GetSlaveMySql()
	qs := o.QueryTable((*Model)(nil))
	ms := make([]*Model, 0)
	qs.Filter("fid__in", fids).Filter("dateline__gte", helper.PreNDayTime(day)).Filter("views__gte", views).Filter("replies__gte", replys).Limit(400).All(&ms)
	return ms
}

func GetByTids(tids []int) []*Model {
	o := boot.GetSlaveMySql()
	qs := o.QueryTable((*Model)(nil))
	ms := make([]*Model, 0)
	qs.Filter("tid__in", tids).All(&ms)
	return ms
}

func GetEssenceThreads(fids []int, day int) []*Model {
	o := boot.GetSlaveMySql()
	qs := o.QueryTable((*Model)(nil))
	ms := make([]*Model, 0)
	qs.Filter("digest", 1).Filter("fid__in", fids).Filter("dateline__gte", helper.PreNDayTime(day)).All(&ms)
	return ms
}
