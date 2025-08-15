package recommendation

import (
	"encoding/json"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/protocol"
)

// TagPreference 标签偏好
type TagPreference struct {
	Tag    string  `json:"tag"`
	Weight float64 `json:"weight"`
	Count  int     `json:"count"`
}

// BehaviorPattern 行为模式
type BehaviorPattern struct {
	Type           string    `json:"type"`
	Count          int       `json:"count"`
	AvgScore       float64   `json:"avgScore"`
	LastBehavior   time.Time `json:"lastBehavior"`
	DailyPattern   []int     `json:"dailyPattern"`   // 24小时分布
	WeeklyPattern  []int     `json:"weeklyPattern"`  // 一周分布
}

// UserProfileBuilder 用户画像构建器
type UserProfileBuilder struct {
	behaviors map[uint][]BehaviorData // userID -> behaviors
	mutex     sync.RWMutex
}

// BehaviorData 行为数据
type BehaviorData struct {
	ItemID       uint                   `json:"itemId"`
	ItemType     string                 `json:"itemType"`
	BehaviorType string                 `json:"behaviorType"`
	Score        float64                `json:"score"`
	Weight       float64                `json:"weight"`
	Tags         []string               `json:"tags"`
	Context      map[string]interface{} `json:"context"`
	Timestamp    time.Time              `json:"timestamp"`
}

// NewUserProfileBuilder 创建用户画像构建器
func NewUserProfileBuilder() *UserProfileBuilder {
	return &UserProfileBuilder{
		behaviors: make(map[uint][]BehaviorData),
	}
}

// AddBehavior 添加用户行为数据
func (upb *UserProfileBuilder) AddBehavior(userID uint, behavior BehaviorData) {
	upb.mutex.Lock()
	defer upb.mutex.Unlock()

	if upb.behaviors[userID] == nil {
		upb.behaviors[userID] = make([]BehaviorData, 0)
	}

	upb.behaviors[userID] = append(upb.behaviors[userID], behavior)
}

// BuildProfile 构建用户画像
func (upb *UserProfileBuilder) BuildProfile(userID uint) *protocol.UserProfile {
	upb.mutex.RLock()
	defer upb.mutex.RUnlock()

	behaviors, exists := upb.behaviors[userID]
	if !exists || len(behaviors) < constant.MinBehaviorCount {
		return &protocol.UserProfile{
			UserID:        userID,
			Preferences:   make(map[string]float64),
			Interests:     []string{},
			BehaviorStats: make(map[string]int),
			LastUpdated:   time.Now().Unix(),
			Metadata:      make(map[string]interface{}),
		}
	}

	// 计算标签偏好
	preferences := upb.calculateTagPreferences(behaviors)

	// 提取兴趣类别
	interests := upb.extractInterests(preferences, 10) // 取前10个兴趣

	// 统计行为模式
	behaviorStats := upb.calculateBehaviorStats(behaviors)

	// 生成元数据
	metadata := upb.generateMetadata(behaviors)

	return &protocol.UserProfile{
		UserID:        userID,
		Preferences:   preferences,
		Interests:     interests,
		BehaviorStats: behaviorStats,
		LastUpdated:   time.Now().Unix(),
		Metadata:      metadata,
	}
}

// calculateTagPreferences 计算标签偏好
func (upb *UserProfileBuilder) calculateTagPreferences(behaviors []BehaviorData) map[string]float64 {
	tagCounts := make(map[string]float64)
	totalWeight := 0.0

	for _, behavior := range behaviors {
		// 获取行为权重
		behaviorWeight := upb.getBehaviorWeight(behavior.BehaviorType)
		
		// 时间衰减因子
		timeDecay := upb.calculateTimeDecay(behavior.Timestamp)
		
		// 最终权重
		finalWeight := behaviorWeight * timeDecay * behavior.Weight

		for _, tag := range behavior.Tags {
			tag = strings.TrimSpace(strings.ToLower(tag))
			if tag != "" {
				tagCounts[tag] += finalWeight
				totalWeight += finalWeight
			}
		}
	}

	// 归一化
	preferences := make(map[string]float64)
	if totalWeight > 0 {
		for tag, count := range tagCounts {
			preferences[tag] = count / totalWeight
		}
	}

	return preferences
}

// getBehaviorWeight 获取行为权重
func (upb *UserProfileBuilder) getBehaviorWeight(behaviorType string) float64 {
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

// calculateTimeDecay 计算时间衰减因子
func (upb *UserProfileBuilder) calculateTimeDecay(timestamp time.Time) float64 {
	hours := time.Since(timestamp).Hours()
	
	// 使用指数衰减，半衰期为30天
	halfLife := 30.0 * 24 // 30天的小时数
	decay := math.Exp(-math.Ln2 * hours / halfLife)
	
	// 限制最小值
	if decay < 0.1 {
		decay = 0.1
	}
	
	return decay
}

// extractInterests 提取兴趣类别
func (upb *UserProfileBuilder) extractInterests(preferences map[string]float64, limit int) []string {
	type tagWeight struct {
		tag    string
		weight float64
	}

	var tagWeights []tagWeight
	for tag, weight := range preferences {
		tagWeights = append(tagWeights, tagWeight{tag: tag, weight: weight})
	}

	// 按权重排序
	sort.Slice(tagWeights, func(i, j int) bool {
		return tagWeights[i].weight > tagWeights[j].weight
	})

	// 提取前N个标签
	var interests []string
	for i, tw := range tagWeights {
		if i >= limit {
			break
		}
		interests = append(interests, tw.tag)
	}

	return interests
}

// calculateBehaviorStats 计算行为统计
func (upb *UserProfileBuilder) calculateBehaviorStats(behaviors []BehaviorData) map[string]int {
	stats := make(map[string]int)
	patterns := make(map[string]*BehaviorPattern)

	for _, behavior := range behaviors {
		// 基础统计
		stats[behavior.BehaviorType]++
		stats["total"]++

		// 按物品类型统计
		stats[behavior.ItemType+"_"+behavior.BehaviorType]++

		// 行为模式分析
		if patterns[behavior.BehaviorType] == nil {
			patterns[behavior.BehaviorType] = &BehaviorPattern{
				Type:          behavior.BehaviorType,
				Count:         0,
				AvgScore:      0,
				DailyPattern:  make([]int, 24),
				WeeklyPattern: make([]int, 7),
			}
		}

		pattern := patterns[behavior.BehaviorType]
		pattern.Count++
		pattern.AvgScore = (pattern.AvgScore*float64(pattern.Count-1) + behavior.Score) / float64(pattern.Count)
		pattern.LastBehavior = behavior.Timestamp

		// 时间模式分析
		hour := behavior.Timestamp.Hour()
		weekday := int(behavior.Timestamp.Weekday())
		pattern.DailyPattern[hour]++
		pattern.WeeklyPattern[weekday]++
	}

	// 添加模式统计
	for behaviorType, pattern := range patterns {
		// 找出活跃时段
		maxHour := 0
		maxWeekday := 0
		for i, count := range pattern.DailyPattern {
			if count > pattern.DailyPattern[maxHour] {
				maxHour = i
			}
		}
		for i, count := range pattern.WeeklyPattern {
			if count > pattern.WeeklyPattern[maxWeekday] {
				maxWeekday = i
			}
		}

		stats[behaviorType+"_active_hour"] = maxHour
		stats[behaviorType+"_active_weekday"] = maxWeekday
		stats[behaviorType+"_avg_score"] = int(pattern.AvgScore * 100) // 保留两位小数
	}

	return stats
}

// generateMetadata 生成元数据
func (upb *UserProfileBuilder) generateMetadata(behaviors []BehaviorData) map[string]interface{} {
	metadata := make(map[string]interface{})

	if len(behaviors) == 0 {
		return metadata
	}

	// 时间范围
	firstBehavior := behaviors[0].Timestamp
	lastBehavior := behaviors[0].Timestamp
	for _, behavior := range behaviors {
		if behavior.Timestamp.Before(firstBehavior) {
			firstBehavior = behavior.Timestamp
		}
		if behavior.Timestamp.After(lastBehavior) {
			lastBehavior = behavior.Timestamp
		}
	}

	metadata["first_behavior"] = firstBehavior.Unix()
	metadata["last_behavior"] = lastBehavior.Unix()
	metadata["behavior_span_days"] = int(lastBehavior.Sub(firstBehavior).Hours() / 24)

	// 活跃度指标
	activeDays := upb.calculateActiveDays(behaviors)
	metadata["active_days"] = activeDays
	metadata["avg_behaviors_per_day"] = float64(len(behaviors)) / float64(activeDays)

	// 多样性指标
	uniqueItems := make(map[uint]bool)
	uniqueTags := make(map[string]bool)
	for _, behavior := range behaviors {
		uniqueItems[behavior.ItemID] = true
		for _, tag := range behavior.Tags {
			uniqueTags[tag] = true
		}
	}
	metadata["unique_items"] = len(uniqueItems)
	metadata["unique_tags"] = len(uniqueTags)
	metadata["diversity_score"] = float64(len(uniqueItems)) / float64(len(behaviors))

	return metadata
}

// calculateActiveDays 计算活跃天数
func (upb *UserProfileBuilder) calculateActiveDays(behaviors []BehaviorData) int {
	activeDays := make(map[string]bool)
	for _, behavior := range behaviors {
		day := behavior.Timestamp.Format("2006-01-02")
		activeDays[day] = true
	}
	return len(activeDays)
}

// GetUserSimilarity 计算用户相似度
func (upb *UserProfileBuilder) GetUserSimilarity(profile1, profile2 *protocol.UserProfile) float64 {
	if profile1 == nil || profile2 == nil {
		return 0.0
	}

	// 计算偏好相似度（余弦相似度）
	var dotProduct, norm1, norm2 float64
	
	allTags := make(map[string]bool)
	for tag := range profile1.Preferences {
		allTags[tag] = true
	}
	for tag := range profile2.Preferences {
		allTags[tag] = true
	}

	for tag := range allTags {
		pref1 := profile1.Preferences[tag]
		pref2 := profile2.Preferences[tag]
		
		dotProduct += pref1 * pref2
		norm1 += pref1 * pref1
		norm2 += pref2 * pref2
	}

	if norm1 == 0 || norm2 == 0 {
		return 0.0
	}

	similarity := dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2))
	
	// 考虑兴趣重叠度
	commonInterests := 0
	for _, interest1 := range profile1.Interests {
		for _, interest2 := range profile2.Interests {
			if interest1 == interest2 {
				commonInterests++
				break
			}
		}
	}
	
	totalInterests := len(profile1.Interests) + len(profile2.Interests) - commonInterests
	interestSimilarity := 0.0
	if totalInterests > 0 {
		interestSimilarity = float64(commonInterests) / float64(totalInterests)
	}

	// 综合相似度（偏好权重70%，兴趣权重30%）
	return similarity*0.7 + interestSimilarity*0.3
}

// UpdateProfileFromBehaviors 从行为数据更新用户画像
func (upb *UserProfileBuilder) UpdateProfileFromBehaviors(userID uint, behaviors []BehaviorData) *protocol.UserProfile {
	upb.mutex.Lock()
	defer upb.mutex.Unlock()

	// 替换行为数据
	upb.behaviors[userID] = behaviors

	// 重新构建画像
	return upb.BuildProfile(userID)
}

// GetRecommendedTags 基于用户画像推荐标签
func (upb *UserProfileBuilder) GetRecommendedTags(profile *protocol.UserProfile, excludeTags []string, limit int) []string {
	if profile == nil || len(profile.Preferences) == 0 {
		return []string{}
	}

	excludeSet := make(map[string]bool)
	for _, tag := range excludeTags {
		excludeSet[strings.ToLower(tag)] = true
	}

	type tagScore struct {
		tag   string
		score float64
	}

	var candidates []tagScore
	for tag, weight := range profile.Preferences {
		if !excludeSet[strings.ToLower(tag)] {
			candidates = append(candidates, tagScore{tag: tag, score: weight})
		}
	}

	// 按分数排序
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].score > candidates[j].score
	})

	var recommendedTags []string
	for i, candidate := range candidates {
		if i >= limit {
			break
		}
		recommendedTags = append(recommendedTags, candidate.tag)
	}

	return recommendedTags
}

// SerializeProfile 序列化用户画像
func (upb *UserProfileBuilder) SerializeProfile(profile *protocol.UserProfile) (map[string]string, error) {
	result := make(map[string]string)

	preferencesBytes, err := json.Marshal(profile.Preferences)
	if err != nil {
		return nil, err
	}
	result["preferences"] = string(preferencesBytes)

	interestsBytes, err := json.Marshal(profile.Interests)
	if err != nil {
		return nil, err
	}
	result["interests"] = string(interestsBytes)

	behaviorStatsBytes, err := json.Marshal(profile.BehaviorStats)
	if err != nil {
		return nil, err
	}
	result["behavior_stats"] = string(behaviorStatsBytes)

	metadataBytes, err := json.Marshal(profile.Metadata)
	if err != nil {
		return nil, err
	}
	result["metadata"] = string(metadataBytes)

	return result, nil
}

// DeserializeProfile 反序列化用户画像
func (upb *UserProfileBuilder) DeserializeProfile(userID uint, data map[string]string, lastUpdated int64) (*protocol.UserProfile, error) {
	profile := &protocol.UserProfile{
		UserID:      userID,
		LastUpdated: lastUpdated,
	}

	if preferencesStr, exists := data["preferences"]; exists {
		if err := json.Unmarshal([]byte(preferencesStr), &profile.Preferences); err != nil {
			return nil, err
		}
	}

	if interestsStr, exists := data["interests"]; exists {
		if err := json.Unmarshal([]byte(interestsStr), &profile.Interests); err != nil {
			return nil, err
		}
	}

	if behaviorStatsStr, exists := data["behavior_stats"]; exists {
		if err := json.Unmarshal([]byte(behaviorStatsStr), &profile.BehaviorStats); err != nil {
			return nil, err
		}
	}

	if metadataStr, exists := data["metadata"]; exists {
		if err := json.Unmarshal([]byte(metadataStr), &profile.Metadata); err != nil {
			return nil, err
		}
	}

	return profile, nil
}