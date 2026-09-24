# SysMon 设计文档

## 1. 背景与目标

提供轻量级、零依赖的本机资源监控代理，作为 Prometheus / 日志系统的补充。

## 2. 架构

\`\`\`
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  /proc/stat  │────▶│  Collector   │────▶│  HTTP Server │
│ /proc/meminfo│     │  (goroutine) │     │  /metrics    │
│  statfs(/)   │     └──────┬───────┘     │  /health     │
└──────────────┘            │             └──────────────┘
                            ▼
                    ┌──────────────┐
                    │ atomic.Value │  ← 最新指标缓存
                    └──────────────┘
                            │
                            ▼
                    ┌──────────────┐
                    │ slog (JSON)  │  → stdout → journald
                    └──────────────┘
\`\`\`

## 3. 数据流

1. 启动时加载配置，校验参数。
2. 后台 goroutine 按 `interval` 周期采集指标。
3. 采集结果写入 `atomic.Value`，供 HTTP handler 读取。
4. HTTP handler 返回 JSON。
5. 收到 `SIGTERM` / `SIGINT` 时优雅关闭。

## 4. 模块划分

| 模块 | 职责 | 文件 |
|------|------|------|
| config | 配置解析、校验、默认值 | `config/config.go` |
| collector | 采集 CPU/内存/磁盘 | `collector/collector.go` |
| main | 生命周期、HTTP 路由、信号处理 | `main.go` |

## 5. 接口

### GET /metrics

\`\`\`json
{
  "cpu_percent": 12.5,
  "memory_percent": 67.2,
  "disk_percent": 45.0,
  "timestamp": 1710000000
}
\`\`\`

### GET /health

\`\`\`json
{"status":"ok"}
\`\`\`

## 6. 配置项

| 字段 | 类型 | 默认 | 说明 |
|------|------|------|------|
| interval | int | 5 | 采集间隔（秒） |
| listen_addr | string | 127.0.0.1:9100 | 监听地址 |
| disk_path | string | / | 磁盘挂载点 |

## 7. 错误处理

- 配置读取失败 → 记录日志，退出码 1。
- 采集失败 → 记录日志，保留上一次指标。
- HTTP 启动失败 → 记录日志，退出码 1。

## 8. 扩展点

- 增加 Prometheus 格式输出。
- 增加网络、负载、进程数指标。
- 增加指标上报到远端。
