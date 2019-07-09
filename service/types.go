package service

const (
	HOT             = "hot"            // 热门
	ESSENCE         = "essence"        // 精华
	WEEK_POPULARITY = "weekPopularity" // 周人气
	CONTRIBUTION    = "contribution"   // 贡献
)

type CallBlockTrait struct {
	IsSetTop  bool   `json:"is_set_top"`
	Subscript string `json:"subscript"`
}

const Separator = "-"
