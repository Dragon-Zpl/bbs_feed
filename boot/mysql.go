package boot

import (
	"bbs_feed/conf"
	"bbs_feed/lib/stringi"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
	"time"
)

var bbs map[string]orm.Ormer

func init() {
	bbs = make(map[string]orm.Ormer, 0)
}

func ConnectMySQL() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	LoadingBbsCluster()
}

func LoadingBbsCluster() {
	var tpl = "{userName}:{password}@tcp({host}:{port})/{DbName}?charset=utf8mb4&loc=Local"
	hosts := strings.Split(conf.MySQLConf.Host, ",")
	var params = stringi.Form{
		"host":     hosts[0],
		"userName": conf.MySQLConf.UserName,
		"password": conf.MySQLConf.Password,
		"port":     conf.MySQLConf.Port,
		"DbName":   conf.MySQLConf.DbName,
	}
	var dsn = stringi.Build(tpl, params)
	orm.RegisterDataBase("default", "mysql", dsn)
	orm.SetMaxIdleConns("default", 200)
	orm.SetMaxOpenConns("default", 200)
	bbs["bbs-0"] = orm.NewOrm()

	for i := 1; i < len(hosts); i++ {
		params["host"] = hosts[i]
		dsn = stringi.Build(tpl, params)
		name := "bbs-" + strconv.Itoa(i)
		orm.RegisterDataBase(name, "mysql", dsn)
		orm.SetMaxIdleConns(name, 100)
		orm.SetMaxOpenConns(name, 200)
		bbs[name] = orm.NewOrm()
	}
}

// mode="w"或只有一个数据库时返回主库
func GetMySQL(mode ...string) orm.Ormer {
	if len(mode) > 0 && mode[0] == "w" {
		bbs["bbs-0"].Using("default")
		return bbs["bbs-0"]
	}
	if len(bbs) == 1 {
		bbs["bbs-0"].Using("default")
		return bbs["bbs-0"]
	}

	var index = time.Now().UnixNano()%int64(len(bbs)-1) + 1
	var key = "bbs-" + string(index)
	db, _ := bbs[key]
	db.Using(key)
	return db
}

// 获取从库方法
func GetSlaveMySql() orm.Ormer {
	return GetMySQL()
}

// 获取主库方法
func GetMasterMysql() orm.Ormer {
	return GetMySQL("w")
}
