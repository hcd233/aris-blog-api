package dao

import (
	"time"

	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	"gorm.io/gorm"
)

// RecommendationLogDAO 推荐日志DAO
//
//	author system
//	update 2025-01-19 12:00:00
type RecommendationLogDAO struct {
	baseDAO[model.RecommendationLog]
}

// CreateLog 创建推荐日志
func (dao *RecommendationLogDAO) CreateLog(db *gorm.DB, log *model.RecommendationLog) error {
	return dao.Create(db, log)
}

// BatchCreateLogs 批量创建推荐日志
func (dao *RecommendationLogDAO) BatchCreateLogs(db *gorm.DB, logs []model.RecommendationLog) error {
	if len(logs) == 0 {
		return nil
	}
	return db.CreateInBatches(logs, 100).Error
}

// GetUserRecommendationLogs 获取用户推荐日志
func (dao *RecommendationLogDAO) GetUserRecommendationLogs(db *gorm.DB, userID uint, requestType string, limit int) ([]model.RecommendationLog, error) {
	var logs []model.RecommendationLog
	query := db.Where("user_id = ?", userID)
	
	if requestType != "" {
		query = query.Where("request_type = ?", requestType)
	}
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Order("timestamp DESC").Find(&logs).Error
	return logs, err
}

// UpdateClickedItems 更新点击的物品
func (dao *RecommendationLogDAO) UpdateClickedItems(db *gorm.DB, logID uint, clickedIDs string) error {
	return db.Model(&model.RecommendationLog{}).
		Where("id = ?", logID).
		Update("clicked_ids", clickedIDs).Error
}

// GetRecommendationStats 获取推荐统计信息
func (dao *RecommendationLogDAO) GetRecommendationStats(db *gorm.DB, startTime, endTime time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	timeRange := []interface{}{startTime.Unix(), endTime.Unix()}
	
	// 总推荐数
	var totalRecommendations int64
	err := db.Model(&model.RecommendationLog{}).
		Where("timestamp BETWEEN ? AND ?", timeRange[0], timeRange[1]).
		Count(&totalRecommendations).Error
	if err != nil {
		return nil, err
	}
	stats["total_recommendations"] = totalRecommendations
	
	// 有点击的推荐数
	var clickedRecommendations int64
	err = db.Model(&model.RecommendationLog{}).
		Where("timestamp BETWEEN ? AND ? AND clicked_ids != ''", timeRange[0], timeRange[1]).
		Count(&clickedRecommendations).Error
	if err != nil {
		return nil, err
	}
	stats["clicked_recommendations"] = clickedRecommendations
	
	// 点击率
	var clickThroughRate float64
	if totalRecommendations > 0 {
		clickThroughRate = float64(clickedRecommendations) / float64(totalRecommendations)
	}
	stats["click_through_rate"] = clickThroughRate
	
	// 按算法统计
	var algorithmStats []struct {
		Algorithm string `json:"algorithm"`
		Count     int64  `json:"count"`
	}
	err = db.Model(&model.RecommendationLog{}).
		Select("algorithm, COUNT(*) as count").
		Where("timestamp BETWEEN ? AND ?", timeRange[0], timeRange[1]).
		Group("algorithm").
		Find(&algorithmStats).Error
	if err != nil {
		return nil, err
	}
	
	algorithmMap := make(map[string]int64)
	for _, stat := range algorithmStats {
		algorithmMap[stat.Algorithm] = stat.Count
	}
	stats["algorithm_stats"] = algorithmMap
	
	// 按请求类型统计
	var typeStats []struct {
		RequestType string `json:"request_type"`
		Count       int64  `json:"count"`
	}
	err = db.Model(&model.RecommendationLog{}).
		Select("request_type, COUNT(*) as count").
		Where("timestamp BETWEEN ? AND ?", timeRange[0], timeRange[1]).
		Group("request_type").
		Find(&typeStats).Error
	if err != nil {
		return nil, err
	}
	
	typeMap := make(map[string]int64)
	for _, stat := range typeStats {
		typeMap[stat.RequestType] = stat.Count
	}
	stats["type_stats"] = typeMap
	
	return stats, nil
}

// GetTopPerformingAlgorithms 获取表现最好的算法
func (dao *RecommendationLogDAO) GetTopPerformingAlgorithms(db *gorm.DB, startTime, endTime time.Time, limit int) ([]struct {
	Algorithm        string  `json:"algorithm"`
	TotalCount       int64   `json:"total_count"`
	ClickedCount     int64   `json:"clicked_count"`
	ClickThroughRate float64 `json:"click_through_rate"`
}, error) {
	var results []struct {
		Algorithm        string  `json:"algorithm"`
		TotalCount       int64   `json:"total_count"`
		ClickedCount     int64   `json:"clicked_count"`
		ClickThroughRate float64 `json:"click_through_rate"`
	}
	
	query := `
		SELECT 
			algorithm,
			COUNT(*) as total_count,
			SUM(CASE WHEN clicked_ids != '' THEN 1 ELSE 0 END) as clicked_count,
			CASE 
				WHEN COUNT(*) > 0 THEN CAST(SUM(CASE WHEN clicked_ids != '' THEN 1 ELSE 0 END) AS FLOAT) / COUNT(*)
				ELSE 0 
			END as click_through_rate
		FROM recommendation_logs 
		WHERE timestamp BETWEEN ? AND ?
		GROUP BY algorithm
		ORDER BY click_through_rate DESC
	`
	
	if limit > 0 {
		query += " LIMIT ?"
		err := db.Raw(query, startTime.Unix(), endTime.Unix(), limit).Scan(&results).Error
		return results, err
	}
	
	err := db.Raw(query, startTime.Unix(), endTime.Unix()).Scan(&results).Error
	return results, err
}

// GetRecentLogs 获取最近的推荐日志
func (dao *RecommendationLogDAO) GetRecentLogs(db *gorm.DB, hours int, limit int) ([]model.RecommendationLog, error) {
	var logs []model.RecommendationLog
	since := time.Now().Add(-time.Duration(hours) * time.Hour).Unix()
	
	query := db.Where("timestamp >= ?", since).Order("timestamp DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Find(&logs).Error
	return logs, err
}

// DeleteOldLogs 删除旧的推荐日志
func (dao *RecommendationLogDAO) DeleteOldLogs(db *gorm.DB, beforeTime time.Time) error {
	return db.Where("timestamp < ?", beforeTime.Unix()).Delete(&model.RecommendationLog{}).Error
}

// GetLogsBySessionID 根据会话ID获取推荐日志
func (dao *RecommendationLogDAO) GetLogsBySessionID(db *gorm.DB, sessionID string) ([]model.RecommendationLog, error) {
	var logs []model.RecommendationLog
	err := db.Where("session_id = ?", sessionID).Order("timestamp DESC").Find(&logs).Error
	return logs, err
}

// GetUserAlgorithmPreference 获取用户算法偏好统计
func (dao *RecommendationLogDAO) GetUserAlgorithmPreference(db *gorm.DB, userID uint, days int) (map[string]float64, error) {
	since := time.Now().AddDate(0, 0, -days).Unix()
	
	var results []struct {
		Algorithm        string  `json:"algorithm"`
		TotalCount       int64   `json:"total_count"`
		ClickedCount     int64   `json:"clicked_count"`
		ClickThroughRate float64 `json:"click_through_rate"`
	}
	
	query := `
		SELECT 
			algorithm,
			COUNT(*) as total_count,
			SUM(CASE WHEN clicked_ids != '' THEN 1 ELSE 0 END) as clicked_count,
			CASE 
				WHEN COUNT(*) > 0 THEN CAST(SUM(CASE WHEN clicked_ids != '' THEN 1 ELSE 0 END) AS FLOAT) / COUNT(*)
				ELSE 0 
			END as click_through_rate
		FROM recommendation_logs 
		WHERE user_id = ? AND timestamp >= ?
		GROUP BY algorithm
		ORDER BY click_through_rate DESC
	`
	
	err := db.Raw(query, userID, since).Scan(&results).Error
	if err != nil {
		return nil, err
	}
	
	preferences := make(map[string]float64)
	for _, result := range results {
		preferences[result.Algorithm] = result.ClickThroughRate
	}
	
	return preferences, nil
}

// GetPopularRecommendedItems 获取热门推荐物品
func (dao *RecommendationLogDAO) GetPopularRecommendedItems(db *gorm.DB, requestType string, hours int, limit int) ([]uint, error) {
	since := time.Now().Add(-time.Duration(hours) * time.Hour).Unix()
	
	// 这里需要解析 item_ids JSON 字段，简化实现假设存储格式为逗号分隔的ID
	// 实际项目中应该使用合适的JSON解析方法
	query := `
		SELECT item_id, COUNT(*) as count
		FROM (
			SELECT unnest(string_to_array(item_ids, ','))::int as item_id
			FROM recommendation_logs 
			WHERE request_type = ? AND timestamp >= ?
		) as items
		GROUP BY item_id
		ORDER BY count DESC
	`
	
	if limit > 0 {
		query += " LIMIT ?"
	}
	
	var results []struct {
		ItemID uint `json:"item_id"`
		Count  int  `json:"count"`
	}
	
	var err error
	if limit > 0 {
		err = db.Raw(query, requestType, since, limit).Scan(&results).Error
	} else {
		err = db.Raw(query, requestType, since).Scan(&results).Error
	}
	
	if err != nil {
		return nil, err
	}
	
	var itemIDs []uint
	for _, result := range results {
		itemIDs = append(itemIDs, result.ItemID)
	}
	
	return itemIDs, nil
}