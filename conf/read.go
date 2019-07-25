package conf

import (
	"github.com/astaxie/beego/logs"
	"github.com/go-ini/ini"
	"go.etcd.io/etcd/clientv3"
	"os"
	"path/filepath"
	"strings"
)

var (
	MySQLConf MySQL
	RedisConf Redis
	EtcdConf  clientv3.Config
	EsConf    Search
)

type Search struct {
	Host  string
}

type MySQL struct {
	Host     string
	Port     string
	UserName string
	Password string
	DbName   string
	Prefix   string
}

type Redis struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func GetRootPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logs.Error(err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func InitConf() {
	confPath := GetRootPath() + "/conf/conf.ini"
	cfg, err := ini.Load(confPath)
	if err != nil {
		panic(err)
	}

	err = cfg.Section("MySQL").MapTo(&MySQLConf)
	if err != nil {
		logs.Error("cfg.MapTo MySQL settings err: %v", err)
	}

	err = cfg.Section("Redis").MapTo(&RedisConf)
	if err != nil {
		logs.Error("cfg.MapTo Redis settings err: %v", err)
	}
	err = cfg.Section("Search").MapTo(&EsConf)
	if err != nil {
		logs.Error("cfg.MapTo Search settings err: %v", err)
	}
}
