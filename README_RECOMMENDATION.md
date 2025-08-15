# 推荐系统模块

## 概述

本推荐系统是一个高性能、可扩展的个性化推荐服务，基于协同过滤算法和用户画像分析，为用户提供精准的文章和标签推荐。

## 核心特性

### 🤖 智能推荐算法
- **协同过滤算法**：基于矩阵分解的协同过滤，支持SVD、NMF等先进算法
- **用户画像构建**：基于用户行为数据动态构建用户兴趣画像
- **混合推荐策略**：结合协同过滤和基于内容的推荐，提供更准确的推荐结果
- **冷启动处理**：新用户和新物品的冷启动推荐策略

### 📊 用户行为分析
- **多维度行为追踪**：支持浏览、点赞、分享、评论、收藏等行为类型
- **时间衰减权重**：考虑行为的时效性，近期行为权重更高
- **行为模式分析**：分析用户的时间偏好和活跃模式
- **异常行为检测**：识别和过滤异常用户行为

### 🎯 个性化推荐
- **实时推荐**：基于用户当前行为实时调整推荐结果
- **多样性控制**：确保推荐结果的多样性，避免信息茧房
- **可解释推荐**：提供推荐理由，增强用户信任度
- **A/B测试支持**：支持多种推荐策略的A/B测试

### ⚡ 高性能架构
- **缓存优化**：Redis缓存热点数据，显著提升响应速度
- **异步处理**：用户行为异步处理，不影响用户体验
- **批量训练**：定时批量训练模型，保证推荐质量
- **数据压缩**：高效的数据存储和传输格式

## 系统架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   用户行为上报   │    │   推荐请求处理   │    │   用户画像查询   │
│   ReportBehavior│    │ RecommendArticles│    │  GetUserProfile │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Recommendation Service                       │
├─────────────────┬─────────────────┬─────────────────┬───────────┤
│ Collaborative   │ User Profile    │ Content Based   │ Popular   │
│ Filter          │ Builder         │ Recommender     │ Items     │
└─────────────────┴─────────────────┴─────────────────┴───────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Data Access Layer                        │
├─────────────────┬─────────────────┬─────────────────────────────┤
│ UserBehaviorDAO │ UserProfileDAO  │ RecommendationLogDAO        │
└─────────────────┴─────────────────┴─────────────────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   PostgreSQL    │    │      Redis      │    │   定时任务      │
│   持久化存储     │    │   缓存热点数据   │    │   模型训练      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## API 接口

### 用户行为上报

```http
POST /v1/recommendation/behavior
Authorization: Bearer <token>
Content-Type: application/json

{
  "userId": 1,
  "itemId": 123,
  "itemType": "article",
  "behaviorType": "view",
  "score": 4.5,
  "context": {
    "source": "homepage",
    "timestamp": 1642867200
  }
}
```

**响应**：
```json
{
  "data": {
    "success": true,
    "message": "行为上报成功"
  },
  "error": null
}
```

### 推荐文章

```http
GET /v1/recommendation/articles?userId=1&limit=10&excludeIds=1,2,3
Authorization: Bearer <token>
```

**响应**：
```json
{
  "data": {
    "items": [
      {
        "id": 456,
        "type": "article",
        "score": 0.92,
        "reason": "协同过滤推荐",
        "title": "推荐的文章标题",
        "tags": ["技术", "算法"]
      }
    ],
    "total": 10,
    "algorithm": "hybrid",
    "timestamp": 1642867200
  },
  "error": null
}
```

### 推荐标签

```http
GET /v1/recommendation/tags?userId=1&limit=5
Authorization: Bearer <token>
```

**响应**：
```json
{
  "data": {
    "items": [
      {
        "id": 1,
        "type": "tag",
        "score": 0.85,
        "reason": "基于用户偏好推荐",
        "title": "机器学习"
      }
    ],
    "total": 5,
    "algorithm": "user_profile_based",
    "timestamp": 1642867200
  },
  "error": null
}
```

### 获取用户画像

```http
GET /v1/recommendation/profile?userId=1
Authorization: Bearer <token>
```

**响应**：
```json
{
  "data": {
    "profile": {
      "userId": 1,
      "preferences": {
        "机器学习": 0.85,
        "深度学习": 0.72,
        "算法": 0.68
      },
      "interests": ["机器学习", "深度学习", "算法"],
      "behaviorStats": {
        "view": 150,
        "like": 45,
        "share": 12
      },
      "lastUpdated": 1642867200
    }
  },
  "error": null
}
```

## 算法配置

所有算法参数都在 `internal/constant/constant.go` 中定义，可以根据业务需求调整：

```go
const (
    // 协同过滤算法参数
    CFDefaultFactors        = 50    // 矩阵分解因子数量
    CFDefaultIterations     = 100   // 迭代次数
    CFDefaultLearningRate   = 0.01  // 学习率
    CFDefaultRegularization = 0.1   // 正则化参数
    
    // 用户行为权重
    ViewWeight     = 1.0 // 浏览权重
    LikeWeight     = 2.0 // 点赞权重
    ShareWeight    = 3.0 // 分享权重
    CommentWeight  = 2.5 // 评论权重
    CollectWeight  = 4.0 // 收藏权重
    
    // 推荐参数
    DefaultRecommendationLimit = 10   // 默认推荐数量
    MaxRecommendationLimit     = 100  // 最大推荐数量
    MinSimilarityThreshold     = 0.1  // 最小相似度阈值
)
```

## 部署和运行

### 环境要求

- Go 1.21+
- PostgreSQL 12+
- Redis 6.0+

### 数据库初始化

系统会自动创建以下表：
- `user_behaviors` - 用户行为记录
- `user_profiles` - 用户画像数据
- `recommendation_logs` - 推荐日志

### 定时任务

系统包含以下定时任务：
- **每小时**：训练推荐模型
- **每10分钟**：更新活跃用户画像
- **每天凌晨2点**：清理90天前的历史数据

### 性能监控

- 推荐响应时间监控
- 缓存命中率统计
- 算法效果评估
- 用户行为分析报告

## 扩展开发

### 添加新的推荐算法

1. 在 `internal/service/recommendation/` 下创建新的算法文件
2. 实现 `Recommender` 接口
3. 在 `RecommendationService` 中集成新算法

### 自定义用户画像特征

1. 修改 `UserProfile` 结构体
2. 在 `UserProfileBuilder` 中添加新的特征提取逻辑
3. 更新相关的序列化/反序列化方法

### 添加新的行为类型

1. 在 `constant.go` 中定义新的行为类型常量
2. 设置相应的行为权重
3. 更新验证逻辑和处理流程

## 最佳实践

### 数据质量
- 定期清理异常行为数据
- 实现行为数据的去重逻辑
- 监控数据质量指标

### 算法优化
- 根据业务特点调整算法参数
- 定期评估推荐效果
- 实施A/B测试验证改进效果

### 性能优化
- 合理设置缓存过期时间
- 优化数据库查询性能
- 使用批量处理减少数据库压力

## 故障排查

### 常见问题

1. **推荐结果为空**
   - 检查用户是否有足够的行为数据
   - 验证协同过滤模型是否已训练
   - 确认物品数据是否正常

2. **推荐响应慢**
   - 检查Redis缓存是否正常
   - 优化数据库查询
   - 考虑降级到热门推荐

3. **用户画像更新失败**
   - 检查行为数据格式
   - 验证数据库连接
   - 查看日志错误信息

### 日志监控

系统使用zap记录详细日志，关键监控点：
- 推荐请求处理时间
- 算法训练成功/失败状态
- 缓存命中率
- 异常行为检测结果

## 许可证

本推荐系统模块遵循项目整体许可证。