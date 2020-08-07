package model

import (
	"bytes"
	"fmt"
	"github.com/fanghongbo/ops-hbs/utils"
	"sort"
)

type BuiltinMetric struct {
	Metric string
	Tags   string
}

func (u *BuiltinMetric) String() string {
	return fmt.Sprintf(
		"%s/%s",
		u.Metric,
		u.Tags,
	)
}

func DigestBuiltinMetrics(items []BuiltinMetric) string {
	var buf bytes.Buffer

	sort.Sort(BuiltinMetricSlice(items))

	for _, m := range items {
		buf.WriteString(m.String())
	}

	return utils.Md5(buf.String())
}
