package cache

import (
	"github.com/fanghongbo/ops-hbs/common/model"
	"sync"
	"time"
)

type AgentsMeta struct {
	sync.RWMutex
	Data map[string]*model.AgentUpdateInfo
}

var AgentsCache = NewAgentsMeta()

func NewAgentsMeta() *AgentsMeta {
	return &AgentsMeta{Data: make(map[string]*model.AgentUpdateInfo)}
}

func (u *AgentsMeta) Put(req *model.AgentReportRequest) {
	val := &model.AgentUpdateInfo{
		LastUpdate:    time.Now().Unix(),
		ReportRequest: req,
	}

	if agentInfo, exists := u.Get(req.Hostname); !exists ||
		agentInfo.ReportRequest.AgentVersion != req.AgentVersion ||
		agentInfo.ReportRequest.IP != req.IP ||
		agentInfo.ReportRequest.PluginVersion != req.PluginVersion {

		model.UpdateAgent(val)
	}

	// 更新 hbs 时间
	u.Lock()
	u.Data[req.Hostname] = val
	u.Unlock()
}

func (u *AgentsMeta) Get(hostname string) (*model.AgentUpdateInfo, bool) {
	u.RLock()
	defer u.RUnlock()
	val, exists := u.Data[hostname]
	return val, exists
}

func (u *AgentsMeta) Delete(hostname string) {
	u.Lock()
	defer u.Unlock()
	delete(u.Data, hostname)
}

func (u *AgentsMeta) Keys() []string {
	u.RLock()
	defer u.RUnlock()
	count := len(u.Data)
	keys := make([]string, count)
	i := 0
	for hostname := range u.Data {
		keys[i] = hostname
		i++
	}
	return keys
}

func DeleteStaleAgents() {
	duration := time.Hour * time.Duration(24)
	for {
		time.Sleep(duration)
		deleteStaleAgents()
	}
}

func deleteStaleAgents() {
	// 一天都没有心跳的Agent，从内存中干掉
	before := time.Now().Unix() - 3600*24
	keys := AgentsCache.Keys()
	count := len(keys)
	if count == 0 {
		return
	}

	for i := 0; i < count; i++ {
		curr, _ := AgentsCache.Get(keys[i])
		if curr.LastUpdate < before {
			AgentsCache.Delete(curr.ReportRequest.Hostname)
		}
	}
}
