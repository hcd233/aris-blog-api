# 推荐系统使用示例

## 基本使用流程

### 1. 用户行为上报

```bash
# 用户浏览文章
curl -X POST "http://localhost:8080/v1/recommendation/behavior" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 1,
    "itemId": 123,
    "itemType": "article",
    "behaviorType": "view",
    "score": 1.0,
    "context": {
      "source": "homepage",
      "duration": 120
    }
  }'

# 用户点赞文章
curl -X POST "http://localhost:8080/v1/recommendation/behavior" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 1,
    "itemId": 123,
    "itemType": "article",
    "behaviorType": "like",
    "score": 2.0
  }'

# 用户收藏文章
curl -X POST "http://localhost:8080/v1/recommendation/behavior" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 1,
    "itemId": 123,
    "itemType": "article",
    "behaviorType": "collect",
    "score": 4.0
  }'
```

### 2. 获取文章推荐

```bash
# 获取推荐文章
curl -X GET "http://localhost:8080/v1/recommendation/articles?userId=1&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# 排除特定文章的推荐
curl -X GET "http://localhost:8080/v1/recommendation/articles?userId=1&limit=10&excludeIds=1,2,3" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 3. 获取标签推荐

```bash
# 获取推荐标签
curl -X GET "http://localhost:8080/v1/recommendation/tags?userId=1&limit=5" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 4. 查看用户画像

```bash
# 获取用户画像
curl -X GET "http://localhost:8080/v1/recommendation/profile?userId=1" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 5. 管理员操作

```bash
# 训练推荐模型（需要管理员权限）
curl -X POST "http://localhost:8080/v1/recommendation/admin/train" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"

# 更新用户画像（需要管理员权限）
curl -X POST "http://localhost:8080/v1/recommendation/admin/profile/update?userId=1" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

## 典型业务场景

### 场景1：新用户首次访问

```javascript
// 前端代码示例
async function handleNewUserVisit(userId) {
  // 1. 先获取热门内容作为冷启动推荐
  const hotArticles = await fetch('/v1/recommendation/articles?userId=' + userId + '&limit=5');
  
  // 2. 用户开始浏览时上报行为
  await reportBehavior(userId, articleId, 'view');
  
  // 3. 用户有互动行为时继续上报
  if (userLiked) {
    await reportBehavior(userId, articleId, 'like');
  }
}

async function reportBehavior(userId, itemId, behaviorType) {
  return fetch('/v1/recommendation/behavior', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + getToken()
    },
    body: JSON.stringify({
      userId: userId,
      itemId: itemId,
      itemType: 'article',
      behaviorType: behaviorType,
      context: {
        timestamp: Date.now(),
        source: 'web'
      }
    })
  });
}
```

### 场景2：个性化首页推荐

```javascript
// 首页个性化推荐
async function loadPersonalizedHomepage(userId) {
  try {
    // 获取个性化推荐文章
    const response = await fetch(`/v1/recommendation/articles?userId=${userId}&limit=20`);
    const data = await response.json();
    
    if (data.error) {
      // 推荐失败时降级到热门内容
      return loadHotArticles();
    }
    
    // 渲染推荐内容
    renderArticles(data.data.items);
    
    // 记录推荐展示
    logRecommendationImpression(data.data.algorithm, data.data.items);
    
  } catch (error) {
    console.error('获取推荐失败:', error);
    // 降级处理
    return loadHotArticles();
  }
}
```

### 场景3：实时推荐更新

```javascript
// 基于用户实时行为更新推荐
class RecommendationEngine {
  constructor(userId) {
    this.userId = userId;
    this.behaviorQueue = [];
    this.recommendationCache = new Map();
  }
  
  // 上报用户行为
  async reportBehavior(itemId, behaviorType, context = {}) {
    const behavior = {
      userId: this.userId,
      itemId: itemId,
      itemType: 'article',
      behaviorType: behaviorType,
      context: {
        ...context,
        timestamp: Date.now()
      }
    };
    
    // 立即上报
    await this.sendBehavior(behavior);
    
    // 如果是重要行为，刷新推荐缓存
    if (['like', 'share', 'collect'].includes(behaviorType)) {
      this.clearRecommendationCache();
    }
  }
  
  // 获取推荐内容
  async getRecommendations(type = 'article', limit = 10) {
    const cacheKey = `${type}-${limit}`;
    
    // 检查缓存
    if (this.recommendationCache.has(cacheKey)) {
      const cached = this.recommendationCache.get(cacheKey);
      if (Date.now() - cached.timestamp < 300000) { // 5分钟缓存
        return cached.data;
      }
    }
    
    // 获取新推荐
    const url = `/v1/recommendation/${type}s?userId=${this.userId}&limit=${limit}`;
    const response = await fetch(url, {
      headers: { 'Authorization': 'Bearer ' + getToken() }
    });
    
    const data = await response.json();
    
    // 缓存结果
    this.recommendationCache.set(cacheKey, {
      data: data.data,
      timestamp: Date.now()
    });
    
    return data.data;
  }
  
  clearRecommendationCache() {
    this.recommendationCache.clear();
  }
}
```

## 性能优化建议

### 1. 批量上报行为数据

```javascript
// 批量上报优化
class BehaviorBatcher {
  constructor(userId, batchSize = 10, flushInterval = 5000) {
    this.userId = userId;
    this.batchSize = batchSize;
    this.flushInterval = flushInterval;
    this.behaviors = [];
    this.timer = null;
  }
  
  addBehavior(behavior) {
    this.behaviors.push({
      ...behavior,
      userId: this.userId,
      timestamp: Date.now()
    });
    
    if (this.behaviors.length >= this.batchSize) {
      this.flush();
    } else if (!this.timer) {
      this.timer = setTimeout(() => this.flush(), this.flushInterval);
    }
  }
  
  async flush() {
    if (this.behaviors.length === 0) return;
    
    const behaviorsToSend = [...this.behaviors];
    this.behaviors = [];
    
    if (this.timer) {
      clearTimeout(this.timer);
      this.timer = null;
    }
    
    // 批量发送（需要后端支持）
    try {
      await fetch('/v1/recommendation/behaviors/batch', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + getToken()
        },
        body: JSON.stringify({ behaviors: behaviorsToSend })
      });
    } catch (error) {
      console.error('批量上报行为失败:', error);
      // 失败时可以重试或记录
    }
  }
}
```

### 2. 推荐结果预加载

```javascript
// 推荐预加载策略
class RecommendationPreloader {
  constructor(userId) {
    this.userId = userId;
    this.preloadedData = new Map();
  }
  
  // 预加载推荐数据
  async preloadRecommendations() {
    const tasks = [
      this.preloadArticles(),
      this.preloadTags(),
      this.preloadUserProfile()
    ];
    
    await Promise.allSettled(tasks);
  }
  
  async preloadArticles() {
    try {
      const response = await fetch(`/v1/recommendation/articles?userId=${this.userId}&limit=50`);
      const data = await response.json();
      this.preloadedData.set('articles', {
        data: data.data,
        timestamp: Date.now()
      });
    } catch (error) {
      console.error('预加载文章推荐失败:', error);
    }
  }
  
  getPreloadedArticles(limit = 10) {
    const cached = this.preloadedData.get('articles');
    if (!cached || Date.now() - cached.timestamp > 600000) { // 10分钟过期
      return null;
    }
    
    return {
      ...cached.data,
      items: cached.data.items.slice(0, limit)
    };
  }
}
```

## 监控和分析

### 推荐效果监控

```javascript
// 推荐效果跟踪
class RecommendationTracker {
  constructor() {
    this.metrics = {
      impressions: 0,
      clicks: 0,
      conversions: 0
    };
  }
  
  // 记录推荐展示
  trackImpression(recommendations) {
    this.metrics.impressions += recommendations.length;
    
    // 发送监控数据
    this.sendMetrics('impression', {
      count: recommendations.length,
      algorithm: recommendations.algorithm,
      timestamp: Date.now()
    });
  }
  
  // 记录推荐点击
  trackClick(itemId, position) {
    this.metrics.clicks++;
    
    this.sendMetrics('click', {
      itemId: itemId,
      position: position,
      timestamp: Date.now()
    });
  }
  
  // 计算点击率
  getClickThroughRate() {
    return this.metrics.impressions > 0 
      ? this.metrics.clicks / this.metrics.impressions 
      : 0;
  }
}
```

这些示例展示了如何在实际项目中集成和使用推荐系统，包括基本的API调用、性能优化策略和效果监控方法。