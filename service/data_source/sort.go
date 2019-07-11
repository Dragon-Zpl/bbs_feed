package data_source

import "bbs_feed/model/forum_thread"

// 热门贴的排序规则
type HotThread []*forum_thread.Model


func (this HotThread) Len() int {
	return len(this)
}

func (this HotThread) Less(i, j int) bool {
	if this[i].Replies < this[j].Replies {
		return true
	} else if this[i].Replies == this[j].Replies {
		if this[i].FavTimes < this[j].FavTimes {
			return true
		} else if this[i].FavTimes == this[j].FavTimes {
			if this[i].Dateline < this[j].Dateline {
				return true
			}
		}
	}
	return false
}

func (this HotThread) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

// 精华贴的排序规则
type EssenceThread []*forum_thread.Model


func (this EssenceThread) Len() int {
	return len(this)
}

func (this EssenceThread) Less(i, j int) bool {
	if this[i].Replies < this[j].Replies {
		return true
	} else if this[i].Replies == this[j].Replies {
		if this[i].FavTimes < this[j].FavTimes {
			return true
		} else if this[i].FavTimes == this[j].FavTimes {
			if this[i].Dateline < this[j].Dateline {
				return true
			}
		}
	}
	return false
}

func (this EssenceThread) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}


// 今日导读的排序规则
type TodayIntroThread []*forum_thread.Model


func (this TodayIntroThread) Len() int {
	return len(this)
}

func (this TodayIntroThread) Less(i, j int) bool {
	if this[i].Dateline < this[j].Dateline {
		return true
	}
	return false
}

func (this TodayIntroThread) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}


// 最新最热的排序规则
type NewHotThread []*forum_thread.Model

func (this NewHotThread) Len() int {
	return len(this)
}

func (this NewHotThread) Less(i, j int) bool {
	if this[i].Replies < this[j].Replies {
		return true
	} else if this[i].Replies == this[j].Replies {
		if this[i].Views < this[j].Views {
			return true
		}
	}
	return false
}

func (this NewHotThread) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

