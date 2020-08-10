# ops-hbs

基于open-falcon二次开发的hbs组件

## 特性
- 使用gorm来操作数据库;
- 新增运行日志配置, 支持日志滚动;
- 新增cpu核心绑定、内存阈值配置; 当agent内存达到阈值的50%时, 打印告警信息；当内存达到阈值的100%, 程序直接退出;


## 编译

it is a golang classic project

``` shell
cd $GOPATH/src/github.com/fanghongbo/ops-hbs/
./control build
./control start
```

## 配置
Refer to `cfg.example.json`, modify the file name to `cfg.json` :

```config
{
  "debug": false,
  "log": {
    "log_level": "INFO",
    "log_path": "./logs",
    "log_file_name": "run.log",
    "log_keep_hours": 3
  },
  "database": {
    "host": "127.0.0.1",
    "user": "falcon",
    "password": "test123456",
    "port": 3306,
    "db": "falcon_portal",
    "max_conn": 20,
    "max_idle": 100
  },
  "rpc": {
    "enabled": true,
    "listen": ":6030"
  },
  "http": {
    "enabled": true,
    "listen": ":6031"
  },
  "max_cpu_rate": 0.2,
  "max_mem_rate": 0.3
}

```

## License

This software is licensed under the Apache License. See the LICENSE file in the top distribution directory for the full license text.
