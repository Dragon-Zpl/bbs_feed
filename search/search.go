package search

import (
	"bbs_feed/boot"
	"bbs_feed/conf"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic"
	"io"
	"strings"
	"sync"
	"time"
)

const (
	ScrollSize = 2000
)

// 添加索引前缀
func addIndexPrefix(index string) string {
	return conf.EsConf.Index + "_" + index
}

type UserAction struct {
	ThreadSupported int `json:"thread_supported"` //帖子被加分
	ThreadReplied   int `json:"thread_replied"` //帖子被回复
	ThreadCollected int `json:"thread_collected"` //帖子被收藏
	PublishThread   int `json:"publish_thread"`   //发帖
	PublishPost     int `json:"publish_post"`    //评论
	PostSupported   int `json:"post_supported"`  //评论被加分
}

type User struct {
	Uid string
	Action UserAction
}



func Search(index string) (map[string][]User, error) {
	index = addIndexPrefix(index)
	once := sync.Once{}
	var searchResults []*elastic.SearchResult
	//游标查询
	svc := boot.InstanceSearchCli().Client.Scroll(index).Size(ScrollSize)
	for {
		searchResult, err := svc.Do(context.TODO())
		if err == io.EOF {
			break
		}
		once.Do(func() {
			searchResults = make([]*elastic.SearchResult, 0, searchResult.Hits.TotalHits)
		})
		searchResults = append(searchResults, searchResult)
	}
	dataMap := make(map[string][]User)
	for _, searchResult := range searchResults {
		for _, hit := range searchResult.Hits.Hits {
			id := strings.Split(hit.Id, "_")
			uid := id[0]
			fid := id[1]
			item := new(UserAction)
			err := json.Unmarshal([]byte(*hit.Source), &item)
			if err != nil {
				return nil, err
			}
			user := User{
				Uid: uid,
				Action: *item,
			}
			if _, ok := dataMap[fid]; ok {
				dataMap[fid] = append(dataMap[fid], user)
				continue
			}
			data := make([]User,0)
			data = append(data, user)
			dataMap[fid] = data
		}
	}
	svc.Clear(context.TODO())
	svc.Do(context.TODO())
	return dataMap, nil
}

//条件查询
func FactorSearch(index string, source interface{}) (map[string]map[string]interface{}, error) {
	index = addIndexPrefix(index)
	result := make(map[string]map[string]interface{})
	searchResults := make([]*elastic.SearchResult, 0, ScrollSize)
	//游标查询
	svc := boot.InstanceSearchCli().Client.Scroll(index).Body(source).Size(ScrollSize)
	for {
		searchResult, err := svc.Do(context.TODO())
		if err == io.EOF || searchResult == nil {
			break
		}
		searchResults = append(searchResults, searchResult)
	}
	for _, searchResult := range searchResults {
		for _, hit := range searchResult.Hits.Hits {
			item := make(map[string]interface{})
			err := json.Unmarshal([]byte(*hit.Source), &item)
			if err != nil {
				return nil, err
			}
			result[hit.Id] = item
		}
	}
	svc.Clear(context.TODO())
	svc.Do(context.TODO())
	return result, nil
}

// 批量创建数据
func CreateBulkIndex(index string, docs map[string]interface{}) error {
	index = addIndexPrefix(index)
	begin := time.Now()
	bulk := boot.InstanceSearchCli().Client.Bulk().Index(index).Type("doc")
	for key, value := range docs {
		bulk.Add(elastic.NewBulkIndexRequest().Id(key).Doc(value))
	}
	res, err := bulk.Do(context.TODO())
	if err != nil {
		return err
	}
	if res.Errors {
		return errors.New("bulk commit failed")
	}
	dur := time.Since(begin).Seconds()
	sec := int(dur)
	pps := int64(float64(len(docs)) / dur)
	fmt.Printf("%-30s %10d | %10d req/s | %02d:%02d\n", "Insert Data To ES", len(docs), pps, sec/60, sec%60)
	return nil
}

// 创建索引
func CreateIndex(index string) (err error) {
	_, err = boot.InstanceSearchCli().Client.CreateIndex(addIndexPrefix(index)).
		Do(context.TODO())
	return
}

//删除索引
func DeleteIndex(index string) (err error) {
	_, err = boot.InstanceSearchCli().Client.DeleteIndex(addIndexPrefix(index)).
		Do(context.TODO())
	return
}

//判断索引是否存在
func IsExistIndex(index string) (b bool) {
	b, _ = boot.InstanceSearchCli().Client.IndexExists(addIndexPrefix(index)).
		Do(context.TODO())
	return
}
