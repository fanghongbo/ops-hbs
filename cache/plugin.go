package cache

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-hbs/common/model"
	"sort"
	"sync"
)

var GroupPluginsCache = NewGroupPluginsMeta()

type GroupPluginsMeta struct {
	sync.RWMutex
	Data map[int64][]string
}

func (u *GroupPluginsMeta) GetPlugins(gid int64) ([]string, bool) {
	var (
		plugins []string
		exist   bool
	)
	u.RLock()
	defer u.RUnlock()
	plugins, exist = u.Data[gid]
	return plugins, exist
}

func NewGroupPluginsMeta() *GroupPluginsMeta {
	return &GroupPluginsMeta{Data: make(map[int64][]string)}
}

func InitGroupPluginsCache() {
	var (
		data  []model.Plugin
		cache map[int64][]string
		err   error
	)

	data, err = model.Plugin{}.GetAll()
	if err != nil {
		dlog.Errorf("get group plugins err: %s", err)
		return
	}

	cache = make(map[int64][]string)
	for _, item := range data {
		if _, exists := cache[item.GrpId]; exists {
			cache[item.GrpId] = append(cache[item.GrpId], item.Dir)
		} else {
			cache[item.GrpId] = []string{item.Dir}
		}
	}

	GroupPluginsCache.Lock()
	defer GroupPluginsCache.Unlock()
	GroupPluginsCache.Data = cache
}

func GetPlugins(hostname string) []string {
	var (
		hid        int64
		gidList    []int64
		exist      bool
		pluginDirs map[string]struct{}
		size       int
		dirs       []string
	)

	hid, exist = HostCache.GetID(hostname)
	if !exist {
		return []string{}
	}

	gidList, exist = HostGroupsCache.GetGroupIds(hid)
	if !exist {
		return []string{}
	}

	pluginDirs = make(map[string]struct{})
	for _, gid := range gidList {
		plugins, exists := GroupPluginsCache.GetPlugins(gid)
		if !exists {
			continue
		}

		for _, plugin := range plugins {
			pluginDirs[plugin] = struct{}{}
		}
	}

	size = len(pluginDirs)
	if size == 0 {
		return []string{}
	}

	dirs = make([]string, size)
	i := 0
	for dir := range pluginDirs {
		dirs[i] = dir
		i++
	}

	sort.Strings(dirs)
	return dirs
}
