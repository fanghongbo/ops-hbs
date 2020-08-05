package g

import "context"

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

	return nil
}
