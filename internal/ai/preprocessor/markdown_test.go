package preprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const source = `# Aris-url-gen

[ [English](README.md) | 简体中文 ]

## 介绍

一个高性能的短链接生成服务，基于Go语言开发。项目名称来源于游戏《碧蓝档案》中的角色Aris，如下图所示。

---

<p align="center">
  <img src="https://raw.githubusercontent.com/hcd233/Aris-AI/master/assets/110531412.jpg" width="50%">
  <br>Aris: Blue Archive 中的角色
</p>

---

## 功能特性

- 生成短链接：将长URL转换为短链接
- 支持自定义过期时间
- 双向缓存：使用Redis实现高性能缓存
- RESTful API：提供标准的HTTP接口
- 数据持久化：使用MySQL存储URL映射关系

## 技术栈

- **Web框架**: [Fiber](https://github.com/gofiber/fiber)
- **ORM**: [GORM](https://gorm.io)
- **缓存**: Redis
- **数据库**: MySQL
- **日志**: [Zap](https://github.com/uber-go/zap)

## API接口

### 1. 生成短链接

` + "```http" + `
POST /v1/shortURL
Content-Type: application/json

{
    "originalURL": "https://example.com/very/long/url",
    "expireDays": 7  // 可选，过期天数
}
` + "```" + `

### 2. 访问短链接

` + "```http" + `
GET /v1/s/{shortURL}
` + "```" + `

## 项目结构

` + "```" + `
.
├── cmd/                # 命令行入口
├── internal/          
│   ├── api/           # API相关代码
│   │   ├── dao/       # 数据访问层
│   │   ├── dto/       # 数据传输对象
│   │   ├── handler/   # 请求处理器
│   │   └── service/   # 业务逻辑层
│   ├── config/        # 配置
│   ├── logger/        # 日志
│   ├── resource/      # 资源
│   └── util/          # 工具函数
└── main.go            # 主入口
` + "```" + `

## 安装部署

### 前置要求

- Go 1.20+
- MySQL 8.0+
- Redis 6.0+

### 本地开发

1. 克隆仓库

` + "```bash" + `
git clone https://github.com/hcd233/Aris-url-gen.git
cd Aris-url-gen
` + "```" + `

2. 安装依赖

` + "```bash" + `
go mod download
` + "```" + `

3. 配置环境变量

参考 ` + "`api.env.template`" + ` 配置相关环境变量

4. 运行服务

` + "```bash" + `
go run main.go
` + "```" + `

## 部署方式

### Docker 部署

1. 创建必要的数据卷:

` + "```bash" + `
docker volume create mysql-data
docker volume create redis-data
` + "```" + `

2. 使用 docker-compose 部署:

` + "```bash" + `
# 开发环境部署
docker compose -f docker/debug-docker-compose.yml up -d

# 生产环境部署
docker compose -f docker/docker-compose.yml up -d
` + "```" + `

### Kubernetes 部署

1. 创建命名空间和配置:

` + "```bash" + `
kubectl apply -f kubernetes/namespace.yml
kubectl apply -f kubernetes/configmaps.yml
` + "```" + `

2. 创建 secrets (需要先配置 secrets.yml):

` + "```bash" + `
cp kubernetes/secrets.yml.template kubernetes/secrets.yml
# 编辑 secrets.yml 填入实际的密钥值
kubectl apply -f kubernetes/secrets.yml
` + "```" + `

3. 创建存储和部署服务:

` + "```bash" + `
kubectl apply -f kubernetes/pvc.yml
kubectl apply -f kubernetes/deployments.yml
kubectl apply -f kubernetes/services.yml
` + "```" + `

### Helm 部署

1. 配置 values.yaml:

` + "```bash" + `
cp helm/aris-url-gen/values.yaml.template helm/aris-url-gen/values.yaml
# 编辑 values.yaml 填入实际的配置值
` + "```" + `

2. 使用 Helm 安装:

` + "```bash" + `
helm install aris-url-gen helm/aris-url-gen
` + "```" + `

3. 升级或卸载:

` + "```bash" + `
# 升级
helm upgrade aris-url-gen helm/aris-url-gen

# 卸载
helm uninstall aris-url-gen
` + "```" + `

## 许可证

本项目采用 Apache License 2.0 许可证。详见 [LICENSE](LICENSE) 文件。`

func TestMarkdownPreprocessor_Process(t *testing.T) {
	tests := []struct {
		name             string
		source           string
		chunkSize        uint
		chunkOverlap     uint
		minHeaderLevel   uint
		allowedTokens    []string
		disallowedTokens []string
		wantChunks       int
		wantHeaders      []string
		wantErr          bool
	}{
		{
			name:             "基本分割测试",
			source:           source,
			chunkSize:        1024,
			chunkOverlap:     256,
			minHeaderLevel:   2,
			allowedTokens:    nil,
			disallowedTokens: nil,
			wantChunks:       10,
			wantHeaders:      nil,
			wantErr:          false,
		},
		{
			name:             "空文档测试",
			source:           "",
			chunkSize:        1024,
			chunkOverlap:     256,
			minHeaderLevel:   1,
			allowedTokens:    nil,
			disallowedTokens: nil,
			wantChunks:       0,
			wantHeaders:      nil,
			wantErr:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewMarkdownPreprocessor(
				tt.chunkSize,
				tt.chunkOverlap,
				tt.allowedTokens,
				tt.disallowedTokens,
			)
			docs, err := p.Process(tt.source)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantChunks, len(docs), "chunk数量不匹配")

			if tt.wantChunks > 0 {
				headers := make([]string, len(docs))
				for i, doc := range docs {
					if h, ok := doc.Metadata["header"]; ok {
						headers[i] = h.(string)
					} else {
						headers[i] = ""
					}
				}
				assert.Equal(t, tt.wantHeaders, headers, "header不匹配")

				// 验证metadata中的其他字段
				for i, doc := range docs {
					assert.Equal(t, i, doc.Metadata["index"], "index不匹配")
					assert.Equal(t, len(docs), doc.Metadata["total"], "total不匹配")
				}
			}
		})
	}
}

func TestMarkdownPreprocessor_Constructor(t *testing.T) {
	tests := []struct {
		name             string
		chunkSize        uint
		chunkOverlap     uint
		minHeaderLevel   uint
		allowedTokens    []string
		disallowedTokens []string
		wantChunkSize    uint
		wantOverlap      uint
		wantHeaderLevel  uint
	}{
		{
			name:             "使用默认值",
			chunkSize:        0,
			chunkOverlap:     0,
			minHeaderLevel:   0,
			allowedTokens:    nil,
			disallowedTokens: nil,
			wantChunkSize:    defaultChunkSize,
			wantOverlap:      defaultChunkOverlap,
			wantHeaderLevel:  defaultMinHeaderLevel,
		},
		{
			name:             "自定义值",
			chunkSize:        2048,
			chunkOverlap:     512,
			minHeaderLevel:   2,
			allowedTokens:    []string{"<special>"},
			disallowedTokens: []string{"<unk>"},
			wantChunkSize:    2048,
			wantOverlap:      512,
			wantHeaderLevel:  2,
		},
		{
			name:             "无效的标题等级",
			chunkSize:        1024,
			chunkOverlap:     256,
			minHeaderLevel:   7,
			allowedTokens:    nil,
			disallowedTokens: nil,
			wantChunkSize:    1024,
			wantOverlap:      256,
			wantHeaderLevel:  defaultMinHeaderLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewMarkdownPreprocessor(
				tt.chunkSize,
				tt.chunkOverlap,
				tt.allowedTokens,
				tt.disallowedTokens,
			)

			assert.Equal(t, tt.wantChunkSize, p.chunkSize, "chunkSize不匹配")
			assert.Equal(t, tt.wantOverlap, p.chunkOverlap, "chunkOverlap不匹配")
			assert.Equal(t, tt.wantHeaderLevel, p.minHeaderLevel, "minHeaderLevel不匹配")
			assert.Equal(t, tt.allowedTokens, p.allowedSpecialTokens, "allowedTokens不匹配")
			assert.Equal(t, tt.disallowedTokens, p.disallowedSpecialTokens, "disallowedTokens不匹配")
		})
	}
}
