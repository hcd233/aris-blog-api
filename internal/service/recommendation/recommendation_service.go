package recommendation

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/database/dao"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RecommendationService 推荐服务
type RecommendationService struct {
	db                  *gorm.DB
	redis               *redis.Client
	logger              *zap.Logger
	collaborativeFilter *CollaborativeFilter
	userProfileBuilder  *UserProfileBuilder
	userBehaviorDAO     *dao.UserBehaviorDAO
	userProfileDAO      *dao.UserProfileDAO
	recommendationLogDAO *dao.RecommendationLogDAO
	articleDAO          *dao.ArticleDAO
	tagDAO              *dao.TagDAO
	mutex               sync.RWMutex
}

// NewRecommendationService 创建推荐服务实例
func NewRecommendationService(db *gorm.DB, redis *redis.Client, logger *zap.Logger) *RecommendationService {
	return &RecommendationService{
		db:                  db,
		redis:               redis,
		logger:              logger,
		collaborativeFilter: NewCollaborativeFilter(),
		userProfileBuilder:  NewUserProfileBuilder(),
		userBehaviorDAO:     dao.GetUserBehaviorDAO(),
		userProfileDAO:      dao.GetUserProfileDAO(),
		recommendationLogDAO: dao.GetRecommendationLogDAO(),
		articleDAO:          dao.GetArticleDAO(),
		tagDAO:              dao.GetTagDAO(),
	}
}

// ReportBehavior 用户行为上报
func (rs *RecommendationService) ReportBehavior(ctx context.Context, req *protocol.UserBehaviorRequest) (*protocol.UserBehaviorResponse, error) {
	// 验证请求参数
	if req.UserID == 0 || req.ItemID == 0 {
		return &protocol.UserBehaviorResponse{
			Success: false,
			Message: "用户ID和物品ID不能为空",
		}, nil
	}

	// 构建行为记录
	behavior := &model.UserBehavior{
		UserID:       req.UserID,
		ItemID:       req.ItemID,
		ItemType:     req.ItemType,
		BehaviorType: req.BehaviorType,
		Score:        req.Score,
		Weight:       rs.getBehaviorWeight(req.BehaviorType),
		Timestamp:    time.Now().Unix(),
	}

	// 序列化上下文信息
	if req.Context != nil {
		contextBytes, err := json.Marshal(req.Context)
		if err != nil {
			rs.logger.Error("序列化上下文信息失败", zap.Error(err))
		} else {
			behavior.Context = string(contextBytes)
		}
	}

	// 保存到数据库
	if err := rs.userBehaviorDAO.CreateBehavior(rs.db, behavior); err != nil {
		rs.logger.Error("保存用户行为失败", zap.Error(err))
		return &protocol.UserBehaviorResponse{
			Success: false,
			Message: "保存行为数据失败",
		}, err
	}

	// 异步更新协同过滤模型和用户画像
	go rs.updateModelsAsync(req.UserID, behavior)

	return &protocol.UserBehaviorResponse{
		Success: true,
		Message: "行为上报成功",
	}, nil
}

// RecommendArticles 推荐文章
func (rs *RecommendationService) RecommendArticles(ctx context.Context, req *protocol.RecommendationRequest) (*protocol.RecommendationResponse, error) {
	limit := req.Limit
	if limit <= 0 || limit > constant.MaxRecommendationLimit {
		limit = constant.DefaultRecommendationLimit
	}

	// 尝试从缓存获取推荐结果
	cacheKey := fmt.Sprintf("%s%d:article:%d", constant.RecommendationCachePrefix, req.UserID, limit)
	cachedResult, err := rs.getRecommendationFromCache(ctx, cacheKey)
	if err == nil && cachedResult != nil {
		return cachedResult, nil
	}

	// 获取用户画像
	userProfile, err := rs.getUserProfile(req.UserID)
	if err != nil {
		rs.logger.Error("获取用户画像失败", zap.Error(err))
		// 降级到热门推荐
		return rs.getPopularArticles(ctx, req)
	}

	var recommendations []protocol.RecommendationItem
	var algorithm string

	// 基于协同过滤的推荐
	cfRecommendations := rs.collaborativeFilter.Recommend(int(req.UserID), rs.uintSliceToIntSlice(req.ExcludeIDs), limit*2)
	if len(cfRecommendations) > 0 {
		cfItems := rs.buildArticleRecommendations(cfRecommendations, "协同过滤", limit/2)
		recommendations = append(recommendations, cfItems...)
		algorithm = "collaborative_filtering"
	}

	// 基于内容的推荐（基于用户画像）
	contentItems, err := rs.getContentBasedRecommendations(userProfile, req.ExcludeIDs, req.IncludeTags, limit-len(recommendations))
	if err == nil {
		recommendations = append(recommendations, contentItems...)
		if algorithm == "" {
			algorithm = "content_based"
		} else {
			algorithm = "hybrid"
		}
	}

	// 如果推荐数量不足，补充热门文章
	if len(recommendations) < limit {
		popularItems, err := rs.getPopularArticleItems(limit - len(recommendations))
		if err == nil {
			recommendations = append(recommendations, popularItems...)
		}
	}

	// 去重并限制数量
	recommendations = rs.deduplicateRecommendations(recommendations)
	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	response := &protocol.RecommendationResponse{
		Items:     recommendations,
		Total:     len(recommendations),
		Algorithm: algorithm,
		Timestamp: time.Now().Unix(),
	}

	// 缓存推荐结果
	rs.cacheRecommendationResult(ctx, cacheKey, response)

	// 记录推荐日志
	go rs.logRecommendation(req.UserID, "article", algorithm, recommendations)

	return response, nil
}

// RecommendTags 推荐标签
func (rs *RecommendationService) RecommendTags(ctx context.Context, req *protocol.RecommendationRequest) (*protocol.RecommendationResponse, error) {
	limit := req.Limit
	if limit <= 0 || limit > constant.MaxRecommendationLimit {
		limit = constant.DefaultRecommendationLimit
	}

	// 获取用户画像
	userProfile, err := rs.getUserProfile(req.UserID)
	if err != nil {
		rs.logger.Error("获取用户画像失败", zap.Error(err))
		// 降级到热门标签推荐
		return rs.getPopularTags(ctx, req)
	}

	var recommendations []protocol.RecommendationItem

	// 基于用户画像推荐标签
	excludeTagNames := make([]string, 0)
	recommendedTags := rs.userProfileBuilder.GetRecommendedTags(userProfile, excludeTagNames, limit*2)

	for i, tagName := range recommendedTags {
		if i >= limit {
			break
		}

		score := userProfile.Preferences[tagName]
		recommendations = append(recommendations, protocol.RecommendationItem{
			ID:    uint(i + 1), // 这里应该是真实的标签ID，简化处理
			Type:  "tag",
			Score: score,
			Reason: "基于用户偏好推荐",
			Title: tagName,
			Tags:  []string{tagName},
		})
	}

	response := &protocol.RecommendationResponse{
		Items:     recommendations,
		Total:     len(recommendations),
		Algorithm: "user_profile_based",
		Timestamp: time.Now().Unix(),
	}

	// 记录推荐日志
	go rs.logRecommendation(req.UserID, "tag", "user_profile_based", recommendations)

	return response, nil
}

// GetUserProfile 获取用户画像
func (rs *RecommendationService) GetUserProfile(ctx context.Context, req *protocol.UserProfileRequest) (*protocol.UserProfileResponse, error) {
	profile, err := rs.getUserProfile(req.UserID)
	if err != nil {
		return &protocol.UserProfileResponse{Profile: nil}, err
	}

	return &protocol.UserProfileResponse{Profile: profile}, nil
}

// TrainModel 训练推荐模型
func (rs *RecommendationService) TrainModel(ctx context.Context) error {
	rs.logger.Info("开始训练推荐模型")

	// 获取用户行为数据
	behaviors, err := rs.userBehaviorDAO.GetRecentBehaviors(rs.db, 24*30, 0) // 获取最近30天的行为
	if err != nil {
		return fmt.Errorf("获取行为数据失败: %w", err)
	}

	// 构建协同过滤训练数据
	for _, behavior := range behaviors {
		rating := behavior.Score
		if rating == 0 {
			rating = behavior.Weight // 如果没有显式评分，使用权重作为评分
		}
		rs.collaborativeFilter.AddRating(int(behavior.UserID), int(behavior.ItemID), rating)
	}

	// 训练协同过滤模型
	if err := rs.collaborativeFilter.Train(); err != nil {
		return fmt.Errorf("训练协同过滤模型失败: %w", err)
	}

	rs.logger.Info("推荐模型训练完成")
	return nil
}

// UpdateUserProfile 更新用户画像
func (rs *RecommendationService) UpdateUserProfile(ctx context.Context, userID uint) error {
	// 获取用户行为数据
	behaviors, err := rs.userBehaviorDAO.GetUserBehaviors(rs.db, userID, "", 1000)
	if err != nil {
		return fmt.Errorf("获取用户行为数据失败: %w", err)
	}

	// 转换为用户画像构建器需要的格式
	var behaviorData []BehaviorData
	for _, behavior := range behaviors {
		tags := rs.getItemTags(behavior.ItemID, behavior.ItemType)
		
		var context map[string]interface{}
		if behavior.Context != "" {
			json.Unmarshal([]byte(behavior.Context), &context)
		}

		behaviorData = append(behaviorData, BehaviorData{
			ItemID:       behavior.ItemID,
			ItemType:     behavior.ItemType,
			BehaviorType: behavior.BehaviorType,
			Score:        behavior.Score,
			Weight:       behavior.Weight,
			Tags:         tags,
			Context:      context,
			Timestamp:    time.Unix(behavior.Timestamp, 0),
		})
	}

	// 构建用户画像
	profile := rs.userProfileBuilder.UpdateProfileFromBehaviors(userID, behaviorData)

	// 序列化并保存到数据库
	serializedData, err := rs.userProfileBuilder.SerializeProfile(profile)
	if err != nil {
		return fmt.Errorf("序列化用户画像失败: %w", err)
	}

	dbProfile := &model.UserProfile{
		UserID:        userID,
		Preferences:   serializedData["preferences"],
		Interests:     serializedData["interests"],
		BehaviorStats: serializedData["behavior_stats"],
		Metadata:      serializedData["metadata"],
		LastUpdated:   profile.LastUpdated,
	}

	if err := rs.userProfileDAO.UpsertProfile(rs.db, dbProfile); err != nil {
		return fmt.Errorf("保存用户画像失败: %w", err)
	}

	// 缓存用户画像
	rs.cacheUserProfile(ctx, userID, profile)

	return nil
}

// 辅助方法

func (rs *RecommendationService) getBehaviorWeight(behaviorType string) float64 {
	switch behaviorType {
	case constant.BehaviorTypeView:
		return constant.ViewWeight
	case constant.BehaviorTypeLike:
		return constant.LikeWeight
	case constant.BehaviorTypeShare:
		return constant.ShareWeight
	case constant.BehaviorTypeComment:
		return constant.CommentWeight
	case constant.BehaviorTypeCollect:
		return constant.CollectWeight
	default:
		return 1.0
	}
}

func (rs *RecommendationService) updateModelsAsync(userID uint, behavior *model.UserBehavior) {
	// 更新协同过滤模型
	rating := behavior.Score
	if rating == 0 {
		rating = behavior.Weight
	}
	rs.collaborativeFilter.AddRating(int(userID), int(behavior.ItemID), rating)

	// 更新用户画像（异步）
	go func() {
		if err := rs.UpdateUserProfile(context.Background(), userID); err != nil {
			rs.logger.Error("异步更新用户画像失败", zap.Error(err))
		}
	}()
}

func (rs *RecommendationService) getUserProfile(userID uint) (*protocol.UserProfile, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("%s%d", constant.UserProfileCachePrefix, userID)
	cachedProfile, err := rs.getUserProfileFromCache(context.Background(), cacheKey)
	if err == nil && cachedProfile != nil {
		return cachedProfile, nil
	}

	// 从数据库获取
	dbProfile, err := rs.userProfileDAO.GetProfileByUserID(rs.db, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 用户画像不存在，创建新的
			return rs.createNewUserProfile(userID)
		}
		return nil, err
	}

	// 反序列化
	data := map[string]string{
		"preferences":    dbProfile.Preferences,
		"interests":      dbProfile.Interests,
		"behavior_stats": dbProfile.BehaviorStats,
		"metadata":       dbProfile.Metadata,
	}

	profile, err := rs.userProfileBuilder.DeserializeProfile(userID, data, dbProfile.LastUpdated)
	if err != nil {
		return nil, err
	}

	// 缓存结果
	rs.cacheUserProfile(context.Background(), userID, profile)

	return profile, nil
}

func (rs *RecommendationService) createNewUserProfile(userID uint) (*protocol.UserProfile, error) {
	// 创建空的用户画像
	profile := &protocol.UserProfile{
		UserID:        userID,
		Preferences:   make(map[string]float64),
		Interests:     []string{},
		BehaviorStats: make(map[string]int),
		LastUpdated:   time.Now().Unix(),
		Metadata:      make(map[string]interface{}),
	}

	// 尝试更新用户画像
	go func() {
		if err := rs.UpdateUserProfile(context.Background(), userID); err != nil {
			rs.logger.Error("创建用户画像时更新失败", zap.Error(err))
		}
	}()

	return profile, nil
}

func (rs *RecommendationService) getItemTags(itemID uint, itemType string) []string {
	// 根据物品类型获取标签
	// 这里简化实现，实际项目中应该查询数据库
	return []string{"default"}
}

func (rs *RecommendationService) uintSliceToIntSlice(uintSlice []uint) []int {
	intSlice := make([]int, len(uintSlice))
	for i, v := range uintSlice {
		intSlice[i] = int(v)
	}
	return intSlice
}

func (rs *RecommendationService) buildArticleRecommendations(cfResults []RecommendationResult, reason string, limit int) []protocol.RecommendationItem {
	var items []protocol.RecommendationItem
	for i, result := range cfResults {
		if i >= limit {
			break
		}
		
		items = append(items, protocol.RecommendationItem{
			ID:     uint(result.ItemID),
			Type:   "article",
			Score:  result.Score,
			Reason: reason,
		})
	}
	return items
}

func (rs *RecommendationService) getContentBasedRecommendations(profile *protocol.UserProfile, excludeIDs []uint, includeTags []string, limit int) ([]protocol.RecommendationItem, error) {
	// 基于用户画像的内容推荐
	var items []protocol.RecommendationItem
	
	// 根据用户偏好标签推荐相关文章
	for tag, weight := range profile.Preferences {
		if len(items) >= limit {
			break
		}
		
		// 这里应该查询数据库获取包含该标签的文章
		// 简化实现
		items = append(items, protocol.RecommendationItem{
			ID:     uint(len(items) + 1),
			Type:   "article",
			Score:  weight,
			Reason: fmt.Sprintf("基于兴趣标签: %s", tag),
			Tags:   []string{tag},
		})
	}
	
	return items, nil
}

func (rs *RecommendationService) getPopularArticles(ctx context.Context, req *protocol.RecommendationRequest) (*protocol.RecommendationResponse, error) {
	// 热门文章推荐
	limit := req.Limit
	if limit <= 0 {
		limit = constant.DefaultRecommendationLimit
	}

	popularItems, err := rs.getPopularArticleItems(limit)
	if err != nil {
		return nil, err
	}

	return &protocol.RecommendationResponse{
		Items:     popularItems,
		Total:     len(popularItems),
		Algorithm: "popular",
		Timestamp: time.Now().Unix(),
	}, nil
}

func (rs *RecommendationService) getPopularArticleItems(limit int) ([]protocol.RecommendationItem, error) {
	popularData, err := rs.userBehaviorDAO.GetPopularItems(rs.db, "article", "", 24*7, limit)
	if err != nil {
		return nil, err
	}

	var items []protocol.RecommendationItem
	for _, item := range popularData {
		items = append(items, protocol.RecommendationItem{
			ID:     item.ItemID,
			Type:   "article",
			Score:  float64(item.Count),
			Reason: "热门推荐",
		})
	}

	return items, nil
}

func (rs *RecommendationService) getPopularTags(ctx context.Context, req *protocol.RecommendationRequest) (*protocol.RecommendationResponse, error) {
	// 热门标签推荐 - 简化实现
	var items []protocol.RecommendationItem
	defaultTags := []string{"技术", "生活", "娱乐", "学习", "工作"}
	
	for i, tag := range defaultTags {
		if i >= req.Limit {
			break
		}
		items = append(items, protocol.RecommendationItem{
			ID:     uint(i + 1),
			Type:   "tag",
			Score:  float64(len(defaultTags) - i),
			Reason: "热门标签",
			Title:  tag,
		})
	}

	return &protocol.RecommendationResponse{
		Items:     items,
		Total:     len(items),
		Algorithm: "popular_tags",
		Timestamp: time.Now().Unix(),
	}, nil
}

func (rs *RecommendationService) deduplicateRecommendations(items []protocol.RecommendationItem) []protocol.RecommendationItem {
	seen := make(map[uint]bool)
	var result []protocol.RecommendationItem

	for _, item := range items {
		if !seen[item.ID] {
			seen[item.ID] = true
			result = append(result, item)
		}
	}

	return result
}

func (rs *RecommendationService) logRecommendation(userID uint, requestType, algorithm string, items []protocol.RecommendationItem) {
	var itemIDs []string
	var scores []string

	for _, item := range items {
		itemIDs = append(itemIDs, strconv.Itoa(int(item.ID)))
		scores = append(scores, fmt.Sprintf("%.3f", item.Score))
	}

	log := &model.RecommendationLog{
		UserID:      userID,
		RequestType: requestType,
		Algorithm:   algorithm,
		ItemIDs:     fmt.Sprintf("[%s]", fmt.Sprintf(`"%s"`, itemIDs)),
		Scores:      fmt.Sprintf("[%s]", fmt.Sprintf("%s", scores)),
		Timestamp:   time.Now().Unix(),
	}

	if err := rs.recommendationLogDAO.CreateLog(rs.db, log); err != nil {
		rs.logger.Error("记录推荐日志失败", zap.Error(err))
	}
}

// 缓存相关方法

func (rs *RecommendationService) getRecommendationFromCache(ctx context.Context, key string) (*protocol.RecommendationResponse, error) {
	data, err := rs.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var response protocol.RecommendationResponse
	if err := json.Unmarshal([]byte(data), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (rs *RecommendationService) cacheRecommendationResult(ctx context.Context, key string, response *protocol.RecommendationResponse) {
	data, err := json.Marshal(response)
	if err != nil {
		rs.logger.Error("序列化推荐结果失败", zap.Error(err))
		return
	}

	if err := rs.redis.Set(ctx, key, data, time.Duration(constant.CacheExpiration)*time.Second).Err(); err != nil {
		rs.logger.Error("缓存推荐结果失败", zap.Error(err))
	}
}

func (rs *RecommendationService) getUserProfileFromCache(ctx context.Context, key string) (*protocol.UserProfile, error) {
	data, err := rs.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var profile protocol.UserProfile
	if err := json.Unmarshal([]byte(data), &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (rs *RecommendationService) cacheUserProfile(ctx context.Context, userID uint, profile *protocol.UserProfile) {
	key := fmt.Sprintf("%s%d", constant.UserProfileCachePrefix, userID)
	data, err := json.Marshal(profile)
	if err != nil {
		rs.logger.Error("序列化用户画像失败", zap.Error(err))
		return
	}

	if err := rs.redis.Set(ctx, key, data, time.Duration(constant.CacheExpiration)*time.Second).Err(); err != nil {
		rs.logger.Error("缓存用户画像失败", zap.Error(err))
	}
}