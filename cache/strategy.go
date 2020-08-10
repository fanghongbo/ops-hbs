package cache

import (
	"bytes"
	"fmt"
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-hbs/common/model"
	"github.com/fanghongbo/ops-hbs/utils"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var StrategiesCache = NewStrategiesMeta()

type StrategiesMeta struct {
	sync.RWMutex
	Data map[int64]model.Strategy
}

func (u *StrategiesMeta) GetMap() map[int64]model.Strategy {
	u.RLock()
	defer u.RUnlock()
	return u.Data
}

func NewStrategiesMeta() *StrategiesMeta {
	return &StrategiesMeta{Data: make(map[int64]model.Strategy)}
}

func InitStrategiesCache() {
	var (
		data  []model.Strategy
		cache map[int64]model.Strategy
		err   error
	)

	data, err = model.Strategy{}.GetAll()
	if err != nil {
		dlog.Errorf("get strategy err: %s", err)
		return
	}

	cache = map[int64]model.Strategy{}
	for _, item := range data {
		cache[item.ID] = item
	}

	StrategiesCache.Lock()
	defer StrategiesCache.Unlock()
	StrategiesCache.Data = cache
}

func GetBuiltinMetrics(hostname string) ([]model.BuiltinMetric, error) {
	var (
		hId          int64
		ret          []model.BuiltinMetric
		tIds         []int64
		pIds         []int64
		gIds         []int64
		exist        bool
		allTemplates map[int64]model.Template
		count        int
		tidStrArr    []string
	)

	hId, exist = HostCache.GetID(hostname)
	if !exist {
		return ret, nil
	}

	gIds, exist = HostGroupsCache.GetGroupIds(hId)
	if !exist {
		return ret, nil
	}

	// 根据gid，获取绑定的所有tid
	tIds = []int64{}
	for _, gid := range gIds {
		tIds, exist = GroupTemplatesCache.GetTemplateIds(gid)
		if !exist {
			continue
		}
	}

	if len(tIds) == 0 {
		return ret, nil
	}

	// 继续寻找这些tid的ParentId
	allTemplates = TemplateCache.GetMap()
	pIds = []int64{}
	for _, tid := range tIds {
		pIds = ParentIds(allTemplates, tid)
	}

	// 终于得到了最终的tid列表
	tIds = append(tIds, pIds...)

	// 把tid列表用逗号拼接在一起
	count = len(tIds)
	tidStrArr = make([]string, count)
	for i := 0; i < count; i++ {
		tidStrArr[i] = strconv.FormatInt(tIds[i], 10)
	}

	return model.Strategy{}.GetBuiltinMetrics(strings.Join(tidStrArr, ","))
}

func ParentIds(allTemplates map[int64]model.Template, tid int64) (ret []int64) {
	depth := 0
	for {
		if tid <= 0 {
			break
		}

		if t, exists := allTemplates[tid]; exists {
			ret = append(ret, tid)
			tid = t.ParentId
		} else {
			break
		}

		depth++
		if depth == 10 {
			dlog.Errorf("template inherit cycle. id:", tid)
			return []int64{}
		}
	}

	sz := len(ret)
	if sz <= 1 {
		return
	}

	desc := make([]int64, sz)
	for i, item := range ret {
		j := sz - i - 1
		desc[j] = item
	}

	return desc
}

func SortedTags(tags map[string]string) string {
	var bufferPool = sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}

	if tags == nil {
		return ""
	}

	size := len(tags)

	if size == 0 {
		return ""
	}

	ret := bufferPool.Get().(*bytes.Buffer)
	ret.Reset()
	defer bufferPool.Put(ret)

	if size == 1 {
		for k, v := range tags {
			ret.WriteString(k)
			ret.WriteString("=")
			ret.WriteString(v)
		}
		return ret.String()
	}

	keys := make([]string, size)
	i := 0
	for k := range tags {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	for j, key := range keys {
		ret.WriteString(key)
		ret.WriteString("=")
		ret.WriteString(tags[key])
		if j != size-1 {
			ret.WriteString(",")
		}
	}

	return ret.String()
}

func ParseTag(tag string) map[string]string {
	newTag := make(map[string]string)
	if tag != "" {
		arr := strings.Split(tag, ",")
		for _, tag := range arr {
			kv := strings.SplitN(tag, "=", 2)
			if len(kv) != 2 {
				continue
			}
			newTag[kv[0]] = kv[1]
		}
	}

	return newTag
}

func CalcInheritStrategies(allTemplates map[int64]model.Template, templateIds []int64, tpl2Strategies map[int64][]model.Strategy) []model.Strategy {
	var (
		templateData       [][]int64
		count              int
		strategies         []model.Strategy
		ids                map[int64]struct{}
		uniqueTemplateData [][]int64
	)

	templateData = [][]int64{}
	for _, tid := range templateIds {
		ids := ParentIds(allTemplates, tid)
		if len(ids) <= 0 {
			continue
		}
		templateData = append(templateData, ids)
	}

	count = len(templateData)
	uniqueTemplateData = [][]int64{}
	for i := 0; i < count; i++ {
		var valid bool = true
		for j := 0; j < count; j++ {
			if i == j {
				continue
			}

			if utils.IsSameSlice(templateData[i], templateData[j]) {
				break
			}

			if utils.HasContainSlice(templateData[i], templateData[j]) {
				valid = false
				break
			}
		}

		if valid {
			uniqueTemplateData = append(uniqueTemplateData, templateData[i])
		}
	}

	// 继承覆盖父模板策略，得到每个模板聚合后的策略列表
	strategies = []model.Strategy{}
	ids = make(map[int64]struct{})

	for _, bucket := range uniqueTemplateData {
		// 开始计算一个桶，先计算老的tid，再计算新的，所以可以覆盖
		// 该桶最终结果
		strategiesData := make(map[string][]model.Strategy)
		for _, tid := range bucket {

			// 一个tid对应的策略列表
			idStrategies := make(map[string][]model.Strategy)

			if strategy, ok := tpl2Strategies[tid]; ok {
				for _, s := range strategy {
					uuid := fmt.Sprintf("metric:%s/tags:%v", s.Metric, SortedTags(ParseTag(s.Tags)))
					if _, ok2 := idStrategies[uuid]; ok2 {
						idStrategies[uuid] = append(idStrategies[uuid], s)
					} else {
						idStrategies[uuid] = []model.Strategy{s}
					}
				}
			}

			// 覆盖父模板
			for uuid, ss := range idStrategies {
				strategiesData[uuid] = ss
			}
		}

		lastTemplateId := bucket[len(bucket)-1]

		// 替换所有策略的模板为最新的模板
		for _, ss := range strategiesData {
			for _, s := range ss {
				valStrategy := s
				if _, exist := ids[valStrategy.ID]; !exist {
					if valStrategy.TplId != lastTemplateId {
						valStrategy.TplId = lastTemplateId
					}
					strategies = append(strategies, valStrategy)
					ids[valStrategy.ID] = struct{}{}
				}
			}
		}
	}

	return strategies
}
