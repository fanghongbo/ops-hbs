package cache

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-hbs/common/model"
	"sync"
)

var ExpressionCache = NewExpressionCache()

type ExpressionCacheMeta struct {
	sync.RWMutex
	Data []model.Expression
}

func (u *ExpressionCacheMeta) Get() []model.Expression {
	u.RLock()
	defer u.RUnlock()

	return u.Data
}

func NewExpressionCache() *ExpressionCacheMeta {
	return &ExpressionCacheMeta{}
}

func InitExpressionCache() {
	var (
		data  []model.Expression
		cache []model.Expression
		err   error
	)

	data, err = model.Expression{}.GetAll()
	if err != nil {
		dlog.Errorf("query expression err: %s", err)
		return
	}

	cache = []model.Expression{}
	for _, item := range data {
		item.Metric, item.Tags, err = item.ParseExpression()
		if err != nil {
			dlog.Errorf("parse expression err: %s", err.Error())
			continue
		}
		cache = append(cache, item)
	}

	ExpressionCache.Lock()
	defer ExpressionCache.Unlock()

	ExpressionCache.Data = cache
}
