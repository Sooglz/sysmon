# SysMon

轻量级 Linux 系统资源监控代理。周期性采集 CPU、内存、磁盘使用率，通过 HTTP 接口暴露 JSON 指标。

## 功能

- 采集 CPU / 内存 / 磁盘使用率
- HTTP 接口：`/metrics`、`/health`
- YAML 配置文件 + 命令行覆盖
- JSON 结构化日志
- systemd 服务管理
- deb 包与 Docker 镜像交付

## 快速开始

### 从源码构建

\`\`\`bash
git clone https://github.com/yourname/sysmon.git
cd sysmon
make build
./sysmon --config configs/sysmon.yaml
\`\`\`

### 验证

\`\`\`bash
curl http://127.0.0.1:9100/health
curl http://127.0.0.1:9100/metrics | jq
\`\`\`

## 安装

### deb 包

\`\`\`bash
sudo dpkg -i sysmon_1.0.0_amd64.deb
sudo systemctl status sysmon
\`\`\`

### Docker

\`\`\`bash
docker build -t sysmon:1.0.0 .
docker run -d -p 9100:9100 --name sysmon sysmon:1.0.0
\`\`\`

## 配置

配置文件路径：`/etc/sysmon/sysmon.yaml`

\`\`\`yaml
interval: 5
listen_addr: "127.0.0.1:9100"
disk_path: "/"
\`\`\`

命令行参数：

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--config` | 配置文件路径 | `/etc/sysmon/sysmon.yaml` |
| `--listen` | 监听地址 | 配置文件值 |
| `--interval` | 采集间隔（秒） | 配置文件值 |
| `--version` | 显示版本 | - |

## 接口

- `GET /health` → `{"status":"ok"}`
- `GET /metrics` → 指标 JSON

## 开发

\`\`\`bash
make lint    # 静态检查
make test    # 单元测试
make build   # 编译
make deb     # 打包
\`\`\`

## 许可证

MIT
