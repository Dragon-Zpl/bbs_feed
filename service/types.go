package service

const (
	HOT             = "hot"               // 热门
	ESSENCE         = "essence"           // 精华
	WEEK_POPULARITY = "weekPopularity"    // 周人气
	CONTRIBUTION    = "contribution"      // 贡献
	TODAY_INTRO     = "todayIntroduction" // 今日导读
)

type CallBlockTrait struct {
	IsSetTop  bool   `json:"isSetTop"`  //是否置顶
	Subscript string `json:"subscript"` //下标
}

const Separator = "_"
