package cron

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/resource/database"
	"github.com/hcd233/aris-blog-api/internal/resource/database/dao"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// QuotaCron 配额定时任务
//
//	@author centonhuang
//	@update 2025-09-30 16:04:21
type QuotaCron struct {
	cron    *cron.Cron
	db      *gorm.DB
	userDAO *dao.UserDAO
}

// NewQuotaCron 创建配额定时任务
//
//	@return Cron
//	@author centonhuang
//	@update 2025-09-30 16:04:27
func NewQuotaCron() Cron {
	return &QuotaCron{
		cron: cron.New(
			cron.WithLogger(newCronLoggerAdapter("QuotaCron", logger.Logger())),
		),
		db:      database.GetDBInstance(context.Background()),
		userDAO: dao.GetUserDAO(),
	}
}

// Start 启动定时任务
//
//	@receiver c *QuotaCron
//	@return error
//	@author centonhuang
//	@update 2025-09-30 16:03:59
func (c *QuotaCron) Start() error {
	// debug set 10 seconds
	// c.cron.AddFunc("*/10 * * * * *", c.deliverQuotas)
	entryID, err := c.cron.AddFunc("0 0 * * *", c.deliverQuotas)
	if err != nil {
		logger.Logger().Error("[QuotaCron] add func error", zap.Error(err))
		return err
	}

	logger.Logger().Info("[QuotaCron] add func success", zap.Int("entryID", int(entryID)))

	c.cron.Start()

	return nil
}

func (c *QuotaCron) deliverQuotas() {
	ctx := context.WithValue(context.Background(), constant.CtxKeyTraceID, uuid.New().String())
	logger := logger.WithCtx(ctx)

	param := &dao.PaginateParam{
		PageParam: &dao.PageParam{
			Page:     2,
			PageSize: -1,
		},
	}
	users, pageInfo, err := c.userDAO.Paginate(c.db, []string{"id", "permission", "llm_quota"}, []string{}, param)
	if err != nil {
		logger.Error("[QuotaCron] deliverQuotas get users error", zap.Error(err))
		return
	}
	permissionIDMapping := map[model.Permission][]uint{
		model.PermissionReader:  {},
		model.PermissionCreator: {},
		model.PermissionAdmin:   {},
	}

	for _, user := range *users {
		permissionIDMapping[user.Permission] = append(permissionIDMapping[user.Permission], user.ID)
	}
	logger.Info(
		"[QuotaCron] deliverQuotas stats",
		zap.Int64("total", pageInfo.Total),
		zap.Int("reader", len(permissionIDMapping[model.PermissionReader])),
		zap.Int("creator", len(permissionIDMapping[model.PermissionCreator])),
		zap.Int("admin", len(permissionIDMapping[model.PermissionAdmin])),
	)

	tx := c.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logger.Error("[QuotaCron] deliverQuotas panic", zap.Error(fmt.Errorf("panic occurred: %v", r)))
		} else if err != nil {
			tx.Rollback()
			logger.Error("[QuotaCron] deliverQuotas transaction error", zap.Error(err))
		} else {
			tx.Commit()
		}
	}()

	for permission, userIDs := range permissionIDMapping {
		quota := model.PermissionQuotaMapping[permission]
		err = c.userDAO.BatchDeliverLLMQuota(tx, userIDs, quota)
		if err != nil {
			return
		}
	}

	logger.Info("[QuotaCron] deliverQuotas success")
}
