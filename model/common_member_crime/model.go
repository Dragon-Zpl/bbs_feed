package common_member_crime

import (
	"bbs_feed/boot"
	"bbs_feed/lib/helper"
	"github.com/astaxie/beego/orm"
	"time"
)

const (
	tablename          = "common_member_crime"
	crime_delpost      = iota //删除帖子
	crime_warnpost            //警告帖子
	crime_banpost             //屏蔽帖子
	crime_banspeak            //禁止发言
	crime_banvisit            //禁止访问
	crime_banstatus           //锁定用户
	crime_avatar              //清除头像
	crime_sightml             //清除签名
	crime_customstatus        //清除自定义头衔
	crime_unban               //解禁用户
)

func init() {
	orm.RegisterModelWithPrefix("pre_", new(Model))
}

type Model struct {
	Cid        int    `orm:"pk;column(cid)" json:"cid"`
	Uid        int    `orm:"column(uid)" json:"uid"`
	Operatorid int    `orm:"column(operatorid)" json:"operatorid"`
	Operator   string `orm:"column(operator)" json:"operator"`
	Action     int    `orm:"column(action)" json:"action"`
	Reason     string `orm:"column(reason)" json:"reason"`
	Dateline   int64  `orm:"column(dateline)" json:"dateline"`
}

// 实现表名的接口
func (m *Model) TableName() string {
	return tablename
}

func GetAll(duration time.Duration) []*Model {
	o := boot.GetSlaveMySql()
	qs := o.QueryTable((*Model)(nil))
	m := make([]*Model, 0)
	action := []int{crime_delpost, crime_banpost, crime_banspeak, crime_banvisit, crime_avatar, crime_sightml}
	qs.Filter("dateline__gte", helper.PreMinuteTime(duration+time.Minute)).Filter("action__in", action).All(&m)
	return m
}
