package cache

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-hbs/common/model"
	"sync"
)

var HostCache = NewHostMeta()

type HostMeta struct {
	sync.RWMutex
	Data map[string]int64
}

func (u *HostMeta) GetID(hostname string) (int64, bool) {
	var (
		id    int64
		exist bool
	)

	u.RLock()
	defer u.RUnlock()
	id, exist = u.Data[hostname]
	return id, exist
}

func NewHostMeta() *HostMeta {
	return &HostMeta{Data: make(map[string]int64)}
}

func InitHostCache() {
	var (
		data  []model.Host
		cache map[string]int64
		err   error
	)

	data, err = model.Host{}.GetAll()
	if err != nil {
		dlog.Errorf("query host err: %s", err)
		return
	}

	cache = map[string]int64{}
	for _, item := range data {
		if item.Hostname == "" || item.Ip == "" {
			continue
		}
		cache[item.Hostname] = item.ID
	}

	HostCache.Lock()
	defer HostCache.Unlock()
	HostCache.Data = cache
}

type MonitorHostMeta struct {
	sync.RWMutex
	Data map[int64]model.Host
}

var MonitorHostCache = NewNotInMaintainHostMeta()

func (u *MonitorHostMeta) Get() map[int64]model.Host {
	u.RLock()
	defer u.RUnlock()
	return u.Data
}

func NewNotInMaintainHostMeta() *MonitorHostMeta {
	return &MonitorHostMeta{Data: make(map[int64]model.Host)}
}

func InitMonitorHostCache() {
	var (
		data  []model.Host
		cache map[int64]model.Host
		err   error
	)

	data, err = model.Host{}.GetHostNotInMaintain()
	if err != nil {
		dlog.Errorf("get host not in maintain err: %s", err)
		return
	}

	cache = map[int64]model.Host{}
	for _, item := range data {
		if item.Ip == "" || item.Hostname == "" {
			continue
		}
		cache[item.ID] = item
	}

	MonitorHostCache.Lock()
	defer MonitorHostCache.Unlock()
	MonitorHostCache.Data = cache
}
