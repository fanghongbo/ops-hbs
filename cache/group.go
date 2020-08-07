package cache

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-hbs/common/model"
	"sync"
)

var HostGroupsCache = NewHostGroupsMeta()

type HostGroupsMeta struct {
	sync.RWMutex
	Data map[int64][]int64
}

func (u *HostGroupsMeta) GetGroupIds(hid int64) ([]int64, bool) {
	var (
		gidList []int64
		exist   bool
	)

	u.RLock()
	defer u.RUnlock()

	gidList, exist = u.Data[hid]
	return gidList, exist
}

func NewHostGroupsMeta() *HostGroupsMeta {
	return &HostGroupsMeta{Data: make(map[int64][]int64)}
}

func InitHostGroupsCache() {
	var (
		data  []model.HostGroup
		cache map[int64][]int64
		err   error
	)

	data, err = model.HostGroup{}.GetAll()
	if err != nil {
		dlog.Errorf("query host group err: %s", err)
		return
	}

	cache = map[int64][]int64{}
	for _, item := range data {
		if _, exists := cache[item.HostId]; exists {
			cache[item.HostId] = append(cache[item.HostId], item.GrpId)
		} else {
			cache[item.HostId] = []int64{item.GrpId}
		}
	}

	HostGroupsCache.Lock()
	defer HostGroupsCache.Unlock()
	HostGroupsCache.Data = cache
}
