package cron

import (
	"fmt"

	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/resource/database"
	"github.com/hcd233/aris-blog-api/internal/resource/database/dao"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Cron interface {
	Start()
}

type QuotaCron struct {
	cron    *cron.Cron
	db      *gorm.DB
	userDAO *dao.UserDAO
}

func NewQuotaCron() Cron {
	return &QuotaCron{
		cron: cron.New(
			cron.WithLogger(newCronLoggerAdapter("QuotaCron", logger.Logger())),
		),
		db:      database.GetDBInstance(),
		userDAO: dao.GetUserDAO(),
	}
}

func (c *QuotaCron) Start() {
	// debug set 10 seconds
	// c.cron.AddFunc("every 10s", c.deliverQuotas)
	c.cron.AddFunc("daily", c.deliverQuotas)
	c.cron.Start()
}

func (c *QuotaCron) deliverQuotas() {
	logger := logger.Logger()
	users, pageInfo, err := c.userDAO.Paginate(c.db, []string{"id", "permission", "llm_quota"}, []string{}, 2, -1)
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
