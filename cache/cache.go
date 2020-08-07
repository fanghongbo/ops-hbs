package cache

import "time"

func InitCache() {
	for {
		InitExpressionCache()
		InitHostGroupsCache()
		InitHostCache()
		InitMonitorHostCache()
		InitGroupPluginsCache()
		InitTemplateCache()
		InitGroupTemplatesCache()
		InitHostTemplateIdsMeta()
		InitStrategiesCache()

		time.Sleep(time.Minute)
	}
}
