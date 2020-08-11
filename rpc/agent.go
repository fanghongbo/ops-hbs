package rpc

import (
	"github.com/fanghongbo/ops-hbs/cache"
	"github.com/fanghongbo/ops-hbs/common/model"
	"time"
)

type Agent int

func (t *Agent) MinePlugins(args model.AgentHeartbeatRequest, reply *model.AgentPluginsResponse) error {
	if args.Hostname == "" {
		return nil
	}

	reply.Plugins = cache.GetPlugins(args.Hostname)
	reply.Timestamp = time.Now().Unix()

	return nil
}

func (t *Agent) ReportStatus(args *model.AgentReportRequest, reply *model.SimpleRpcResponse) error {
	if args.Hostname == "" {
		reply.Code = 1
		return nil
	}

	cache.AgentsCache.Put(args)

	return nil
}

func (t *Agent) BuiltinMetrics(args *model.AgentHeartbeatRequest, reply *model.BuiltinMetricResponse) error {
	var (
		metrics  []model.BuiltinMetric
		checksum string
		err      error
	)

	if args.Hostname == "" {
		return nil
	}

	metrics, err = cache.GetBuiltinMetrics(args.Hostname)
	if err != nil {
		return nil
	}

	checksum = ""

	if len(metrics) > 0 {
		checksum = model.DigestBuiltinMetrics(metrics)
	}

	if args.Checksum == checksum {
		reply.Metrics = []model.BuiltinMetric{}
	} else {
		reply.Metrics = metrics
	}

	reply.Checksum = checksum
	reply.Timestamp = time.Now().Unix()

	return nil
}
