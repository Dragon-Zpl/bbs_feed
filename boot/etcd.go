package boot

import (
	"bbs_feed/conf"
	"go.etcd.io/etcd/clientv3"
)

var client *clientv3.Client

func ConnectEtcd() {
	var err error
	client, err = clientv3.New(conf.EtcdConf)
	if err != nil {
		panic(err)
	}
}
