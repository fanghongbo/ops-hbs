package model

import (
	"fmt"
	"github.com/fanghongbo/dlog"
	"time"
)

type AgentReportRequest struct {
	Hostname      string
	IP            string
	AgentVersion  string
	PluginVersion string
}

func (u *AgentReportRequest) String() string {
	return fmt.Sprintf(
		"<Hostname:%s, IP:%s, AgentVersion:%s, PluginVersion:%s>",
		u.Hostname,
		u.IP,
		u.AgentVersion,
		u.PluginVersion,
	)
}

type AgentUpdateInfo struct {
	LastUpdate    int64
	ReportRequest *AgentReportRequest
}

func UpdateAgent(agentInfo *AgentUpdateInfo) {
	var (
		host Host
		err  error
	)

	host = Host{
		Hostname:      agentInfo.ReportRequest.Hostname,
		Ip:            agentInfo.ReportRequest.IP,
		AgentVersion:  agentInfo.ReportRequest.AgentVersion,
		PluginVersion: agentInfo.ReportRequest.PluginVersion,
		UpdateAt:      time.Now(),
	}

	// 查询hostname是否存在，存在则更新记录，否则创建
	if host.IsExist() {
		// 更新
		if err = host.Update(); err != nil {
			dlog.Errorf("update agent err: %s", err)
		}
	} else {
		if err = host.Create(); err != nil {
			dlog.Errorf("create agent err: %s", err)
		}
	}
}
