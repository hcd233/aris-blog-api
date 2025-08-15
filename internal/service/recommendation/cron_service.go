package recommendation

import (
	"context"
	"time"

	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/resource/database/dao"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CronService 推荐系统定时任务服务
type CronService struct {
	cron                *cron.Cron
	recommendationService *RecommendationService
	userBehaviorDAO     *dao.UserBehaviorDAO
	logger              *zap.Logger
}

// NewCronService 创建定时任务服务
func NewCronService(db *gorm.DB, redis *redis.Client, logger *zap.Logger) *CronService {
	return &CronService{
		cron:                cron.New(cron.WithSeconds()),
		recommendationService: NewRecommendationService(db, redis, logger),
		userBehaviorDAO:     dao.GetUserBehaviorDAO(),
		logger:              logger,
	}
}

// Start 启动定时任务
func (cs *CronService) Start() error {
	cs.logger.Info("启动推荐系统定时任务")

	// 每小时训练一次推荐模型
	_, err := cs.cron.AddFunc("0 0 * * * *", cs.trainModelJob)
	if err != nil {
		return err
	}

	// 每10分钟更新一批用户画像
	_, err = cs.cron.AddFunc("0 */10 * * * *", cs.updateUserProfilesJob)
	if err != nil {
		return err
	}

	// 每天凌晨2点清理旧数据
	_, err = cs.cron.AddFunc("0 0 2 * * *", cs.cleanupOldDataJob)
	if err != nil {
		return err
	}

	cs.cron.Start()
	return nil
}

// Stop 停止定时任务
func (cs *CronService) Stop() {
	cs.logger.Info("停止推荐系统定时任务")
	cs.cron.Stop()
}

// trainModelJob 训练推荐模型任务
func (cs *CronService) trainModelJob() {
	cs.logger.Info("开始执行推荐模型训练任务")
	
	ctx := context.Background()
	if err := cs.recommendationService.TrainModel(ctx); err != nil {
		cs.logger.Error("推荐模型训练失败", zap.Error(err))
	} else {
		cs.logger.Info("推荐模型训练完成")
	}
}

// updateUserProfilesJob 更新用户画像任务
func (cs *CronService) updateUserProfilesJob() {
	cs.logger.Info("开始执行用户画像更新任务")
	
	ctx := context.Background()
	
	// 获取最近活跃的用户
	activeUsers, err := cs.userBehaviorDAO.GetActiveUsers(
		cs.recommendationService.db, 
		constant.UserProfileUpdateInterval/3600, // 转换为小时
		constant.MinBehaviorCount,
	)
	if err != nil {
		cs.logger.Error("获取活跃用户失败", zap.Error(err))
		return
	}

	// 批量更新用户画像，限制并发数
	const maxConcurrency = 10
	semaphore := make(chan struct{}, maxConcurrency)
	
	for _, userID := range activeUsers {
		semaphore <- struct{}{} // 获取信号量
		go func(uid uint) {
			defer func() { <-semaphore }() // 释放信号量
			
			if err := cs.recommendationService.UpdateUserProfile(ctx, uid); err != nil {
				cs.logger.Error("更新用户画像失败", 
					zap.Uint("userID", uid), 
					zap.Error(err))
			}
		}(userID)
	}

	// 等待所有任务完成
	for i := 0; i < maxConcurrency; i++ {
		semaphore <- struct{}{}
	}

	cs.logger.Info("用户画像更新任务完成", zap.Int("userCount", len(activeUsers)))
}

// cleanupOldDataJob 清理旧数据任务
func (cs *CronService) cleanupOldDataJob() {
	cs.logger.Info("开始执行数据清理任务")
	
	// 清理90天前的行为数据
	cleanupTime := time.Now().AddDate(0, 0, -90)
	if err := cs.userBehaviorDAO.DeleteOldBehaviors(cs.recommendationService.db, cleanupTime); err != nil {
		cs.logger.Error("清理旧行为数据失败", zap.Error(err))
	} else {
		cs.logger.Info("旧行为数据清理完成")
	}

	// 清理90天前的推荐日志
	recommendationLogDAO := dao.GetRecommendationLogDAO()
	if err := recommendationLogDAO.DeleteOldLogs(cs.recommendationService.db, cleanupTime); err != nil {
		cs.logger.Error("清理旧推荐日志失败", zap.Error(err))
	} else {
		cs.logger.Info("旧推荐日志清理完成")
	}

	cs.logger.Info("数据清理任务完成")
}