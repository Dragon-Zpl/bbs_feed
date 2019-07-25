package boot

import (
	"bbs_feed/conf"
	"context"
	"github.com/astaxie/beego/logs"
	"github.com/olivere/elastic"
)

type SearchClient struct {
	Client *elastic.Client
	host   string
}

var searchCli *SearchClient

func InstanceSearchCli() *SearchClient {
	return searchCli
}

func InitSearchClient() {
	var err error
	searchCli = new(SearchClient)
	searchCli.host = conf.EsConf.Host
	searchCli.Client, err = elastic.NewClient(elastic.SetURL(searchCli.host))
	if err != nil {
		panic(err)
	}
	info, code, err := searchCli.Client.Ping(searchCli.host).Do(context.Background())
	if err != nil {
		panic(err)
	}
	logs.Info("Elasticsearch returned with code %d and version %s", code, info.Version.Number)
}

