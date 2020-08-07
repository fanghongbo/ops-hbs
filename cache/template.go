package cache

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-hbs/common/model"
	"sync"
)

var GroupTemplatesCache = NewGroupTemplatesMeta()

type GroupTemplatesMeta struct {
	sync.RWMutex
	Data map[int64][]int64
}

func (u *GroupTemplatesMeta) GetTemplateIds(gid int64) ([]int64, bool) {
	u.RLock()
	defer u.RUnlock()
	templateIds, exists := u.Data[gid]
	return templateIds, exists
}

func NewGroupTemplatesMeta() *GroupTemplatesMeta {
	return &GroupTemplatesMeta{Data: make(map[int64][]int64)}
}

func InitGroupTemplatesCache() {
	var (
		data  []model.GroupTemplate
		cache map[int64][]int64
		err   error
	)

	data, err = model.GroupTemplate{}.GetAll()
	if err != nil {
		dlog.Errorf("get group templates err: %s", err)
		return
	}

	cache = map[int64][]int64{}
	for _, item := range data {
		if _, exist := cache[item.GrpId]; exist {
			cache[item.GrpId] = append(cache[item.GrpId], item.TplId)
		} else {
			cache[item.GrpId] = []int64{item.TplId}
		}
	}

	GroupTemplatesCache.Lock()
	defer GroupTemplatesCache.Unlock()
	GroupTemplatesCache.Data = cache
}

var TemplateCache = NewTemplateMeta()

type TemplateMeta struct {
	sync.RWMutex
	Data map[int64]model.Template
}

func (u *TemplateMeta) GetMap() map[int64]model.Template {
	u.RLock()
	defer u.RUnlock()
	return u.Data
}

func (u *TemplateMeta) Get(templateId int64) (model.Template, bool) {
	u.RLock()
	defer u.RUnlock()
	val, exists := u.Data[templateId]
	return val, exists
}

func NewTemplateMeta() *TemplateMeta {
	return &TemplateMeta{Data: make(map[int64]model.Template)}
}

func InitTemplateCache() {
	var (
		data  []model.Template
		cache map[int64]model.Template
		err   error
	)

	data, err = model.Template{}.GetAll()
	if err != nil {
		dlog.Errorf("get template err: %s", err)
		return
	}

	cache = map[int64]model.Template{}
	for _, item := range data {
		cache[item.ID] = item
	}

	TemplateCache.Lock()
	defer TemplateCache.Unlock()
	TemplateCache.Data = cache
}


func Template2Strategies(strategies map[int64]model.Strategy) map[int64][]model.Strategy {
	var result map[int64][]model.Strategy

	result = make(map[int64][]model.Strategy)
	for _, s := range strategies {
		if s.TplId == 0 {
			continue
		}

		// 查询模版id是否存在
		_, exist := TemplateCache.Get(s.TplId)
		if !exist {
			continue
		}

		if _, exists := result[s.TplId]; exists {
			result[s.TplId] = append(result[s.TplId], s)
		} else {
			result[s.TplId] = []model.Strategy{s}
		}
	}
	return result
}

var HostTemplateIdsCache = NewHostTemplateIdsMeta()

type HostTemplateIdsMeta struct {
	sync.RWMutex
	Data map[int64][]int64
}

func (u *HostTemplateIdsMeta) GetMap() map[int64][]int64 {
	u.RLock()
	defer u.RUnlock()
	return u.Data
}

func NewHostTemplateIdsMeta() *HostTemplateIdsMeta {
	return &HostTemplateIdsMeta{Data: make(map[int64][]int64)}
}

func InitHostTemplateIdsMeta() {
	var (
		data  []model.HostTemplate
		cache map[int64][]int64
		err   error
	)

	data, err = model.HostTemplate{}.GetAll()
	if err != nil {
		return
	}

	cache = map[int64][]int64{}
	for _, item := range data {
		if _, ok := cache[item.HostId]; ok {
			cache[item.HostId] = append(cache[item.HostId], item.TplId)
		} else {
			cache[item.HostId] = []int64{item.TplId}
		}
	}

	HostTemplateIdsCache.Lock()
	defer HostTemplateIdsCache.Unlock()
	HostTemplateIdsCache.Data = cache
}
