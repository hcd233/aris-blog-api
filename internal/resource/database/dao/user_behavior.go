package dao

import (
	"time"

	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	"gorm.io/gorm"
)

// UserBehaviorDAO 用户行为DAO
//
//	author system
//	update 2025-01-19 12:00:00
type UserBehaviorDAO struct {
	baseDAO[model.UserBehavior]
}

// CreateBehavior 创建用户行为记录
func (dao *UserBehaviorDAO) CreateBehavior(db *gorm.DB, behavior *model.UserBehavior) error {
	return dao.Create(db, behavior)
}

// BatchCreateBehaviors 批量创建用户行为记录
func (dao *UserBehaviorDAO) BatchCreateBehaviors(db *gorm.DB, behaviors []model.UserBehavior) error {
	if len(behaviors) == 0 {
		return nil
	}
	return db.CreateInBatches(behaviors, 100).Error
}

// GetUserBehaviors 获取用户行为记录
func (dao *UserBehaviorDAO) GetUserBehaviors(db *gorm.DB, userID uint, itemType string, limit int) ([]model.UserBehavior, error) {
	var behaviors []model.UserBehavior
	query := db.Where("user_id = ?", userID)
	
	if itemType != "" {
		query = query.Where("item_type = ?", itemType)
	}
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Order("timestamp DESC").Find(&behaviors).Error
	return behaviors, err
}

// GetUserBehaviorsByTimeRange 按时间范围获取用户行为
func (dao *UserBehaviorDAO) GetUserBehaviorsByTimeRange(db *gorm.DB, userID uint, startTime, endTime time.Time) ([]model.UserBehavior, error) {
	var behaviors []model.UserBehavior
	err := db.Where("user_id = ? AND timestamp BETWEEN ? AND ?", userID, startTime.Unix(), endTime.Unix()).
		Order("timestamp DESC").
		Find(&behaviors).Error
	return behaviors, err
}

// GetBehaviorsByItem 获取特定物品的行为记录
func (dao *UserBehaviorDAO) GetBehaviorsByItem(db *gorm.DB, itemID uint, itemType string, behaviorType string) ([]model.UserBehavior, error) {
	var behaviors []model.UserBehavior
	query := db.Where("item_id = ? AND item_type = ?", itemID, itemType)
	
	if behaviorType != "" {
		query = query.Where("behavior_type = ?", behaviorType)
	}
	
	err := query.Order("timestamp DESC").Find(&behaviors).Error
	return behaviors, err
}

// GetBehaviorStats 获取行为统计
func (dao *UserBehaviorDAO) GetBehaviorStats(db *gorm.DB, userID uint) (map[string]interface{}, error) {
	var results []struct {
		BehaviorType string `json:"behavior_type"`
		ItemType     string `json:"item_type"`
		Count        int64  `json:"count"`
		AvgScore     float64 `json:"avg_score"`
	}
	
	err := db.Model(&model.UserBehavior{}).
		Select("behavior_type, item_type, COUNT(*) as count, AVG(score) as avg_score").
		Where("user_id = ?", userID).
		Group("behavior_type, item_type").
		Find(&results).Error
	
	if err != nil {
		return nil, err
	}
	
	stats := make(map[string]interface{})
	for _, result := range results {
		key := result.ItemType + "_" + result.BehaviorType
		stats[key+"_count"] = result.Count
		stats[key+"_avg_score"] = result.AvgScore
	}
	
	return stats, nil
}

// GetRecentBehaviors 获取最近的行为记录
func (dao *UserBehaviorDAO) GetRecentBehaviors(db *gorm.DB, hours int, limit int) ([]model.UserBehavior, error) {
	var behaviors []model.UserBehavior
	since := time.Now().Add(-time.Duration(hours) * time.Hour).Unix()
	
	query := db.Where("timestamp >= ?", since).Order("timestamp DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Find(&behaviors).Error
	return behaviors, err
}

// DeleteOldBehaviors 删除旧的行为记录
func (dao *UserBehaviorDAO) DeleteOldBehaviors(db *gorm.DB, beforeTime time.Time) error {
	return db.Where("timestamp < ?", beforeTime.Unix()).Delete(&model.UserBehavior{}).Error
}

// GetUserItemInteractions 获取用户对特定物品的交互历史
func (dao *UserBehaviorDAO) GetUserItemInteractions(db *gorm.DB, userID uint, itemID uint, itemType string) ([]model.UserBehavior, error) {
	var behaviors []model.UserBehavior
	err := db.Where("user_id = ? AND item_id = ? AND item_type = ?", userID, itemID, itemType).
		Order("timestamp DESC").
		Find(&behaviors).Error
	return behaviors, err
}

// GetPopularItems 获取热门物品
func (dao *UserBehaviorDAO) GetPopularItems(db *gorm.DB, itemType string, behaviorType string, hours int, limit int) ([]struct {
	ItemID uint  `json:"item_id"`
	Count  int64 `json:"count"`
}, error) {
	var results []struct {
		ItemID uint  `json:"item_id"`
		Count  int64 `json:"count"`
	}
	
	query := db.Model(&model.UserBehavior{}).
		Select("item_id, COUNT(*) as count").
		Where("item_type = ?", itemType)
	
	if behaviorType != "" {
		query = query.Where("behavior_type = ?", behaviorType)
	}
	
	if hours > 0 {
		since := time.Now().Add(-time.Duration(hours) * time.Hour).Unix()
		query = query.Where("timestamp >= ?", since)
	}
	
	query = query.Group("item_id").Order("count DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Find(&results).Error
	return results, err
}

// GetUserSimilarityData 获取用户相似度计算所需数据
func (dao *UserBehaviorDAO) GetUserSimilarityData(db *gorm.DB, userIDs []uint, itemType string) (map[uint]map[uint]float64, error) {
	var behaviors []model.UserBehavior
	err := db.Where("user_id IN ? AND item_type = ?", userIDs, itemType).
		Select("user_id, item_id, score").
		Find(&behaviors).Error
	
	if err != nil {
		return nil, err
	}
	
	// 构建用户-物品评分矩阵
	matrix := make(map[uint]map[uint]float64)
	for _, behavior := range behaviors {
		if matrix[behavior.UserID] == nil {
			matrix[behavior.UserID] = make(map[uint]float64)
		}
		// 如果同一用户对同一物品有多个行为，取最高分数
		if existingScore, exists := matrix[behavior.UserID][behavior.ItemID]; !exists || behavior.Score > existingScore {
			matrix[behavior.UserID][behavior.ItemID] = behavior.Score
		}
	}
	
	return matrix, nil
}

// GetActiveUsers 获取活跃用户
func (dao *UserBehaviorDAO) GetActiveUsers(db *gorm.DB, hours int, minBehaviors int) ([]uint, error) {
	var userIDs []uint
	since := time.Now().Add(-time.Duration(hours) * time.Hour).Unix()
	
	err := db.Model(&model.UserBehavior{}).
		Select("user_id").
		Where("timestamp >= ?", since).
		Group("user_id").
		Having("COUNT(*) >= ?", minBehaviors).
		Pluck("user_id", &userIDs).Error
	
	return userIDs, err
}