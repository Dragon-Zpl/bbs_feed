package call_block

/*
 热门调用块
*/
import (
	"bbs_feed/boot"
	"bbs_feed/model/topic_fid_relation"
	"bbs_feed/service"
	"bbs_feed/service/data_source"
	"bbs_feed/service/redis_ops"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

const CALl_BLOCK_HOT_THREAD = "call-block-hot-thread"
const CALl_BLOCK_HOT_THREAD_TRAIT = "call-block-hot-thread_trait"

type HotRules struct {
	day        int
	viewCount  int
	replyCount int
	cronExp    time.Duration // 周期时间
}

type Hot struct {
	name string

	reportChan chan []int
	hotRules   HotRules

	cancel context.CancelFunc
	Ctx    context.Context

	topicIds []string // 数据源

}

func NewHot(topicId int, topicIds []string) *Hot {
	return &Hot{
		name:           fmt.Sprintf("%d-hot", topicId),
		topicIds:       topicIds,
	}
}

// 从reids 去掉 hot thread
func (this *Hot) RemoveReportThread() {
	for {
		select {
		case tids := <-this.reportChan:
			hotThreads := data_source.GetHotThreadByTids(tids)
			for _, thread := range hotThreads {
				if trait, err := boot.InstanceRedisCli(boot.CACHE).HGet(CALl_BLOCK_HOT_THREAD_TRAIT, strconv.Itoa(thread.Thread.Tid)).Result(); err == nil {
					var callBlockTrait service.CallBlockTrait
					if err = json.Unmarshal([]byte(trait), &callBlockTrait); err == nil {
						thread.Trait = callBlockTrait
					}
				}
				if threadBytes, err := json.Marshal(thread); err == nil {
					fmt.Println(threadBytes)
					redis_ops.DelZAdd(CALl_BLOCK_HOT_THREAD, string(threadBytes))
				}
			}
		}
	}
}

func Remove(tids []int) {
	hotThreads := data_source.GetHotThreadByTids(tids)
	for _, thread := range hotThreads {
		if trait, err := boot.InstanceRedisCli(boot.CACHE).HGet(CALl_BLOCK_HOT_THREAD_TRAIT, strconv.Itoa(thread.Thread.Tid)).Result(); err == nil {
			var callBlockTrait service.CallBlockTrait
			if err = json.Unmarshal([]byte(trait), &callBlockTrait); err == nil {
				thread.Trait = callBlockTrait
			}
		}
		if threadBytes, err := json.Marshal(thread); err == nil {
			fmt.Println(string(threadBytes))
			redis_ops.DelZAdd(CALl_BLOCK_HOT_THREAD, string(threadBytes))
		}
	}
}

func (this *Hot) AcceptSign(tids []int) {
	this.reportChan <- tids
	return
}

func (this *Hot) GetThis() interface{} {
	return this
}

func (this *Hot) Init() {
	ctx, cancel := context.WithCancel(context.Background())
	this.Ctx = ctx
	this.cancel = cancel
	this.hotRules = HotRules{
		day:        7,
		viewCount:  2000,
		replyCount: 30,
		cronExp:    10 *time.Minute,
	}
	//this.hotRules = service_confs.Hot
}

func (this *Hot) ChangeConf(conf interface{}) {
	if conf, ok := conf.(HotRules); ok {
		this.hotRules = conf
		this.reStart()
	}
}

func (this *Hot) Start() {
	fmt.Println("start")
	this.start()
	t := time.NewTimer(this.hotRules.cronExp)
	for {
		select {
		case <-t.C:
			this.start()
			t.Reset(this.hotRules.cronExp)
		case <-this.Ctx.Done():
			return
		}
	}
}


func (this *Hot) start() {
	redisThreads := data_source.GetHotSortThread(topic_fid_relation.GetFids(this.topicIds), this.hotRules.day, this.hotRules.viewCount, this.hotRules.replyCount)
	redisTraits, _ := boot.InstanceRedisCli(boot.CACHE).HGetAll(CALl_BLOCK_HOT_THREAD_TRAIT).Result()

	datas := make([]interface{}, 0, len(redisThreads))
	for _, thread := range redisThreads {
		if redisTraits != nil {
			if threadTrait, ok := redisTraits[strconv.Itoa(thread.Thread.Tid)]; ok {
				var callBlockTrait service.CallBlockTrait
				if err := json.Unmarshal([]byte(threadTrait), &callBlockTrait); err == nil {
					thread.Trait = callBlockTrait
				}
			}
		}

		datas = append(datas, thread)
	}
	redis_ops.ZAddSort(CALl_BLOCK_HOT_THREAD, datas)
}

func (this *Hot) ChangeFids(topicIds []string) {
	this.topicIds = topicIds
}

func (this *Hot) Stop() {
	this.cancel()
}

func (this *Hot) GetName() string {
	return this.name
}

func (this *Hot) reStart() {
	this.Stop()
	this.Ctx, this.cancel = context.WithCancel(context.Background())
	this.Start()
}
