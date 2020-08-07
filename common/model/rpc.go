package model

import "fmt"

type SimpleRpcResponse struct {
	Code int `json:"code"`
}

func (u *SimpleRpcResponse) String() string {
	return fmt.Sprintf("<Code: %d>", u.Code)
}

type NullRpcRequest struct {
}

type BuiltinMetricResponse struct {
	Metrics   []BuiltinMetric
	Checksum  string
	Timestamp int64
}

func (u *BuiltinMetricResponse) String() string {
	return fmt.Sprintf(
		"<Metrics:%v, Checksum:%s, Timestamp:%v>",
		u.Metrics,
		u.Checksum,
		u.Timestamp,
	)
}

type AgentHeartbeatRequest struct {
	Hostname string
	Checksum string
}

func (u *AgentHeartbeatRequest) String() string {
	return fmt.Sprintf(
		"<Hostname: %s, Checksum: %s>",
		u.Hostname,
		u.Checksum,
	)
}

type AgentPluginsResponse struct {
	Plugins   []string
	Timestamp int64
}

func (u *AgentPluginsResponse) String() string {
	return fmt.Sprintf(
		"<Plugins:%v, Timestamp:%v>",
		u.Plugins,
		u.Timestamp,
	)
}

type BuiltinMetricSlice []BuiltinMetric

func (u BuiltinMetricSlice) Len() int {
	return len(u)
}

func (u BuiltinMetricSlice) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (u BuiltinMetricSlice) Less(i, j int) bool {
	return u[i].String() < u[j].String()
}

type ExpressionResponse struct {
	Expressions []Expression `json:"expressions"`
}

type HostStrategy struct {
	Hostname   string     `json:"hostname"`
	Strategies []Strategy `json:"strategies"`
}

type StrategiesResponse struct {
	HostStrategies []HostStrategy `json:"hostStrategies"`
}
