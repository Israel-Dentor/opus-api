# Opus API

Claude API 格式到 MorphLLM API 的代理服务，支持流式响应和完整的请求/响应日志记录。

## 功能特性

- **API 格式转换**: 接收 Claude API 格式请求，转换为 MorphLLM API 格式
- **流式响应**: 支持 SSE (Server-Sent Events) 流式输出
- **工具调用**: 完整支持 Claude 的 tool use 功能
- **Token 计数**: 自动计算输入/输出 token 数量
- **详细日志**: 可选的请求/响应完整日志记录，便于调试
- **Docker 支持**: 提供 Dockerfile 和构建脚本

## 快速开始

### 环境要求

- Go 1.21+
- Docker (可选)

### 本地运行

```bash
# 安装依赖
go mod download

# 运行服务
go run cmd/server/main.go
```

服务将在 `http://localhost:3002` 启动。

### Docker 运行

```bash
# 构建镜像
./build.sh

# 运行容器
docker run -p 3002:3002 opus-api
```

## 配置

在 `internal/types/common.go` 中配置上游 API：

```go
const (
    MorphAPIURL = "https://your-morph-api-endpoint"
    DebugMode   = true  // 启用详细日志
    LogDir      = "./logs"
)

var MorphHeaders = map[string]string{
    "Authorization": "Bearer your-token",
    "Content-Type":  "application/json",
}
```

## API 端点

### POST /v1/messages

接收 Claude API 格式的消息请求。

**请求示例**:

```json
{
  "model": "claude-opus-4-5-20251101",
  "max_tokens": 1024,
  "messages": [
    {
      "role": "user",
      "content": "Hello"
    }
  ],
  "stream": true
}
```

**响应**: SSE 流式响应，格式与 Claude API 兼容。

### GET /health

健康检查端点。

## 日志记录

启用 `DebugMode` 后，每个请求会在 `logs/` 目录下创建独立文件夹，包含：

1. `1_claude_request.json` - 原始 Claude 格式请求
2. `2_morph_request.json` - 转换后的 Morph 格式请求
3. `3_upstream_request.txt` - 发送到上游的完整 HTTP 请求
4. `4_upstream_response.txt` - 上游返回的原始响应
5. `5_client_response.txt` - 返回给客户端的最终响应

## 项目结构

```
opus-api/
├── cmd/server/          # 主程序入口
├── internal/
│   ├── converter/       # API 格式转换
│   ├── handler/         # HTTP 请求处理
│   ├── logger/          # 日志记录
│   ├── parser/          # 工具调用解析
│   ├── stream/          # 流式响应处理
│   ├── tokenizer/       # Token 计数
│   └── types/           # 类型定义
├── logs/                # 日志输出目录
└── test/                # 测试文件

```

## 开发

```bash
# 运行测试
go test ./...

# 构建二进制
go build -o server ./cmd/server
```

## License

MIT
