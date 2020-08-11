package g

import (
	"context"
	"github.com/fanghongbo/dlog"
)

func InitAll() {
	InitConfig()
	InitLog()
	InitRuntime()
	InitDB()
}

func Shutdown(ctx context.Context) error {
	defer ctx.Done()

	// 关闭数据库连接
	if db != nil {
		_ = db.Close()
	}

	// 刷新日志
	dlog.Close()

	return nil
}
