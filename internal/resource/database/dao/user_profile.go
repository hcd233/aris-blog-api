package dao

import (
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	"gorm.io/gorm"
)

// UserProfileDAO 用户画像DAO
//
//	author system
//	update 2025-01-19 12:00:00
type UserProfileDAO struct {
	baseDAO[model.UserProfile]
}

// CreateProfile 创建用户画像
func (dao *UserProfileDAO) CreateProfile(db *gorm.DB, profile *model.UserProfile) error {
	return dao.Create(db, profile)
}

// UpdateProfile 更新用户画像
func (dao *UserProfileDAO) UpdateProfile(db *gorm.DB, userID uint, profile *model.UserProfile) error {
	updates := map[string]interface{}{
		"preferences":    profile.Preferences,
		"interests":      profile.Interests,
		"behavior_stats": profile.BehaviorStats,
		"metadata":       profile.Metadata,
		"last_updated":   profile.LastUpdated,
		"version":        gorm.Expr("version + 1"),
	}
	return db.Model(&model.UserProfile{}).Where("user_id = ?", userID).Updates(updates).Error
}

// GetProfileByUserID 通过用户ID获取用户画像
func (dao *UserProfileDAO) GetProfileByUserID(db *gorm.DB, userID uint) (*model.UserProfile, error) {
	var profile model.UserProfile
	err := db.Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// UpsertProfile 创建或更新用户画像
func (dao *UserProfileDAO) UpsertProfile(db *gorm.DB, profile *model.UserProfile) error {
	// 尝试查找现有记录
	var existingProfile model.UserProfile
	err := db.Where("user_id = ?", profile.UserID).First(&existingProfile).Error
	
	if err == gorm.ErrRecordNotFound {
		// 不存在则创建
		return dao.CreateProfile(db, profile)
	} else if err != nil {
		return err
	} else {
		// 存在则更新
		return dao.UpdateProfile(db, profile.UserID, profile)
	}
}

// BatchGetProfiles 批量获取用户画像
func (dao *UserProfileDAO) BatchGetProfiles(db *gorm.DB, userIDs []uint) ([]model.UserProfile, error) {
	var profiles []model.UserProfile
	err := db.Where("user_id IN ?", userIDs).Find(&profiles).Error
	return profiles, err
}

// GetProfilesForUpdate 获取需要更新的用户画像
func (dao *UserProfileDAO) GetProfilesForUpdate(db *gorm.DB, beforeTimestamp int64, limit int) ([]model.UserProfile, error) {
	var profiles []model.UserProfile
	query := db.Where("last_updated < ?", beforeTimestamp).Order("last_updated ASC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Find(&profiles).Error
	return profiles, err
}

// GetProfileStatistics 获取用户画像统计信息
func (dao *UserProfileDAO) GetProfileStatistics(db *gorm.DB) (map[string]interface{}, error) {
	var stats map[string]interface{} = make(map[string]interface{})
	
	// 总用户画像数
	var totalProfiles int64
	err := db.Model(&model.UserProfile{}).Count(&totalProfiles).Error
	if err != nil {
		return nil, err
	}
	stats["total_profiles"] = totalProfiles
	
	// 最近更新的画像数（24小时内）
	var recentUpdates int64
	recentTimestamp := int64(86400) // 24小时前的时间戳（简化计算）
	err = db.Model(&model.UserProfile{}).Where("last_updated > ?", recentTimestamp).Count(&recentUpdates).Error
	if err != nil {
		return nil, err
	}
	stats["recent_updates"] = recentUpdates
	
	// 平均版本号
	var avgVersion float64
	err = db.Model(&model.UserProfile{}).Select("AVG(version)").Scan(&avgVersion).Error
	if err != nil {
		return nil, err
	}
	stats["avg_version"] = avgVersion
	
	return stats, nil
}

// DeleteProfile 删除用户画像
func (dao *UserProfileDAO) DeleteProfile(db *gorm.DB, userID uint) error {
	return db.Where("user_id = ?", userID).Delete(&model.UserProfile{}).Error
}

// GetOutdatedProfiles 获取过期的用户画像
func (dao *UserProfileDAO) GetOutdatedProfiles(db *gorm.DB, beforeTimestamp int64, limit int) ([]uint, error) {
	var userIDs []uint
	query := db.Model(&model.UserProfile{}).
		Select("user_id").
		Where("last_updated < ?", beforeTimestamp)
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Pluck("user_id", &userIDs).Error
	return userIDs, err
}

// GetProfilesByVersion 按版本号获取用户画像
func (dao *UserProfileDAO) GetProfilesByVersion(db *gorm.DB, minVersion int, limit int) ([]model.UserProfile, error) {
	var profiles []model.UserProfile
	query := db.Where("version >= ?", minVersion).Order("version DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Find(&profiles).Error
	return profiles, err
}