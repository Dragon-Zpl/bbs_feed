package elastic_ops

import (
	"bbs_feed/boot"
	"context"
	"fmt"
)


func GetEsData()  {
	client := boot.InstanceSearchCli().Client
	svc, err := client.Scroll("bbs_user_action_2019-07-22").Type("doc").Do(context.Background())
	if err != nil {
		panic(err)
	}
	for _, v := range svc.Hits.Hits {
		fmt.Println(v)
	}
}