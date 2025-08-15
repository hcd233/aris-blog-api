// Package constant 常量
package constant

const (
	CtxKeyUserID     = "userID"
	CtxKeyUserName   = "userName"
	CtxKeyPermission = "permission"
	CtxKeyBody       = "body"
	CtxKeyURI        = "uri"
	CtxKeyParam      = "param"
	CtxKeyTraceID    = "traceID"
	CtxKeyLimiter    = "limiter"

	// ListArticleVersionContentLength 分页查询文章版本中的内容长度限制
	//	update 2025-01-18 23:20:20
	ListArticleVersionContentLength = 100
)

// 推荐系统相关常量
const (
	// 协同过滤算法参数
	CFDefaultFactors        = 50    // 矩阵分解因子数量
	CFDefaultIterations     = 100   // 迭代次数
	CFDefaultLearningRate   = 0.01  // 学习率
	CFDefaultRegularization = 0.1   // 正则化参数
	CFMinRating            = 1.0    // 最小评分
	CFMaxRating            = 5.0    // 最大评分
	CFConvergenceThreshold = 0.0001 // 收敛阈值

	// 推荐参数
	DefaultRecommendationLimit = 10   // 默认推荐数量
	MaxRecommendationLimit     = 100  // 最大推荐数量
	MinSimilarityThreshold     = 0.1  // 最小相似度阈值
	UserProfileUpdateInterval  = 3600 // 用户画像更新间隔(秒)

	// 用户行为权重
	ViewWeight     = 1.0 // 浏览权重
	LikeWeight     = 2.0 // 点赞权重
	ShareWeight    = 3.0 // 分享权重
	CommentWeight  = 2.5 // 评论权重
	CollectWeight  = 4.0 // 收藏权重

	// 缓存相关
	UserProfileCachePrefix     = "user_profile:"
	RecommendationCachePrefix  = "recommendation:"
	SimilarityMatrixCacheKey   = "similarity_matrix"
	CacheExpiration           = 3600 // 缓存过期时间(秒)

	// 数据库相关
	BehaviorBatchSize = 1000 // 批量处理行为数据大小
	MinBehaviorCount  = 5    // 生成推荐所需的最少行为数量
)

// 用户行为类型
const (
	BehaviorTypeView    = "view"    // 浏览
	BehaviorTypeLike    = "like"    // 点赞
	BehaviorTypeShare   = "share"   // 分享
	BehaviorTypeComment = "comment" // 评论
	BehaviorTypeCollect = "collect" // 收藏
)

// 推荐类型
const (
	RecommendationTypeArticle = "article" // 文章推荐
	RecommendationTypeTag     = "tag"     // 标签推荐
	RecommendationTypeUser    = "user"    // 用户推荐
)
