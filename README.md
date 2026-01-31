# gRPC 微服务示例

这是一个统一的 gRPC 微服务项目，同时具备服务端和客户端功能。它可以接收来自其他服务的请求，也可以向其他微服务发起请求。

## 项目结构

```
server_go/
├── proto/                 # Protocol Buffer 定义
│   └── hello.proto        # 服务定义
├── server/                # 服务端实现（接收请求）
│   ├── service.go         # 服务端逻辑，包含客户端功能
├── client/                # 客户端实现（发起请求）
│   └── client.go          # 客户端逻辑，用于调用其他微服务
├── monitoring/            # 监控相关代码
│   └── otel.go            # OpenTelemetry 配置
├── cmd/                   # 主程序入口
│   └── server/            # 统一服务主程序
│       └── main.go
├── docker-compose.yml     # Docker Compose 配置
├── monitoring/            # 监控配置
│   ├── prometheus.yml     # Prometheus 配置
│   └── otel-collector-config.yaml  # OpenTelemetry 收集器配置
└── go.mod, go.sum        # Go 模块文件
```

## 功能说明

### 统一微服务架构
- **服务端功能**: 接收并处理来自其他服务的 gRPC 请求
  - `SayHello`: 一元 RPC 示例
  - `SayHelloStream`: 流式 RPC 示例
- **客户端功能**: 向其他微服务发起 gRPC 请求
  - 连接到其他微服务
  - 转发或代理请求到其他服务

### OpenTelemetry 集成
- 链路追踪：使用 OpenTelemetry 追踪请求路径
- 指标收集：导出到 Prometheus

### 监控系统
- Prometheus：收集和存储指标数据
- Grafana：可视化监控数据
- OpenTelemetry 收集器：接收追踪和指标数据

## 快速开始

### 1. 启动监控系统

```bash
cd /data/workspace/server_go
docker-compose up -d
```

这将启动 Prometheus、Grafana 和 OpenTelemetry 收集器。

### 2. 运行微服务

```bash
cd /data/workspace/server_go
go run cmd/server/main.go
```

## 环境要求

- Go 1.25+
- Docker 和 Docker Compose
- Protobuf 编译器 (protoc)

## 监控访问地址

- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (默认用户名/密码: admin/admin)
- OpenTelemetry 收集器指标: http://localhost:8888

