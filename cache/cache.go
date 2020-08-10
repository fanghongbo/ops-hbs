package cache

import (
	"time"
)

func InitCache() {
	var t1 *time.Timer
	t1 = time.NewTimer(time.Second * 10)

	for {
		select {
		case <-t1.C:
			InitExpressionCache()
			InitHostGroupsCache()
			InitHostCache()
			InitMonitorHostCache()
			InitGroupPluginsCache()
			InitTemplateCache()
			InitHostTemplateIdsMeta()
			InitGroupTemplatesCache()
			InitStrategiesCache()

			// reset timer
			t1.Reset(time.Minute * 1)
		}
	}
}
