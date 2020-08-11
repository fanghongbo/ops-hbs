package rpc

import (
	"github.com/fanghongbo/ops-hbs/cache"
	"github.com/fanghongbo/ops-hbs/common/model"
)

type Hbs int

func (t *Hbs) GetExpressions(req model.NullRpcRequest, reply *model.ExpressionResponse) error {
	reply.Expressions = cache.ExpressionCache.Get()
	return nil
}

func (t *Hbs) GetStrategies(req model.NullRpcRequest, reply *model.StrategiesResponse) error {
	var (
		hostTemplateIds map[int64][]int64
		sz              int
		allTemplates    map[int64]model.Template
		hosts           map[int64]model.Host
		strategies      map[int64]model.Strategy
		tpl2Strategies  map[int64][]model.Strategy
		hostStrategies  []model.HostStrategy
	)

	reply.HostStrategies = []model.HostStrategy{}

	hostTemplateIds = cache.HostTemplateIdsCache.GetMap()
	sz = len(hostTemplateIds)
	if sz == 0 {
		return nil
	}

	hosts = cache.MonitorHostCache.Get()
	if len(hosts) == 0 {
		return nil
	}

	allTemplates = cache.TemplateCache.GetMap()
	if len(allTemplates) == 0 {
		return nil
	}

	strategies = cache.StrategiesCache.GetMap()
	if len(strategies) == 0 {
		return nil
	}

	tpl2Strategies = cache.Template2Strategies(strategies)

	hostStrategies = make([]model.HostStrategy, 0, sz)
	for hostId, tplIds := range hostTemplateIds {
		h, exists := hosts[hostId]
		if !exists {
			continue
		}

		ss := cache.CalcInheritStrategies(allTemplates, tplIds, tpl2Strategies)
		if len(ss) <= 0 {
			continue
		}

		hostStrategies = append(hostStrategies, model.HostStrategy{
			Hostname:   h.Hostname,
			Strategies: ss,
		})
	}

	reply.HostStrategies = hostStrategies
	return nil
}
