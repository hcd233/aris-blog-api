// Package search 搜索中间件
//
//	@update 2024-09-17 11:47:23
package search

import (
	"fmt"

	"github.com/hcd233/Aris-AI-go/internal/config"
	"github.com/hcd233/Aris-AI-go/internal/logger"
	"github.com/hcd233/Aris-AI-go/internal/resource/database/model"
	"github.com/meilisearch/meilisearch-go"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

const (
	userIndex = "user"
)

// ServiceManager meilisearch 服务管理
//
//	@update 2024-09-18 12:29:08
var ServiceManager meilisearch.ServiceManager

func init() {
	ServiceManager = meilisearch.New(
		fmt.Sprintf("http://%s:%s", config.MeilisearchHost, config.MeilisearchPort),
		meilisearch.WithAPIKey(config.MeilisearchMasterKey),
	)
}

// CreateIndex 创建索引
//
//	@return err error
//	@author centonhuang
//	@update 2024-09-18 12:43:35
func CreateIndex() (err error) {
	tasks := []func() error{
		createUserIndex,
		updateUserIndex,
	}
	for _, task := range tasks {
		err = task()
		if err != nil {
			return
		}
	}
	return
}

// DeleteIndex 删除索引
//
//	@return err error
//	@author centonhuang
//	@update 2024-09-18 01:24:35
func DeleteIndex() (err error) {
	info, err := ServiceManager.DeleteIndex(userIndex)
	if err != nil {
		logger.Logger.Error("[Delete Index] failed to delete index", zap.Error(err))
		return
	}

	logger.Logger.Info("[Delete Index] success to delete index", zap.String("Index", userIndex), zap.String("Status", string(info.Status)))
	return
}

func createUserIndex() (err error) {
	info, err := ServiceManager.CreateIndex(&meilisearch.IndexConfig{
		Uid: userIndex,
	})
	if err != nil {
		logger.Logger.Error("[Create Index] failed to create index", zap.Error(err))
		return
	}

	logger.Logger.Info("[Create Index] success to create index", zap.String("Index", userIndex), zap.String("Status", string(info.Status)))
	return
}

func updateUserIndex() (err error) {
	users, err := model.QueryUsers(-1, -1)
	if err != nil {
		logger.Logger.Error("[Update Index] failed to query users", zap.Error(err))
		return
	}

	info, err := ServiceManager.Index(userIndex).AddDocuments(lo.Map(
		users,
		func(user *model.User, _ int) map[string]interface{} {
			return model.GetUserBasicInfo(user)
		},
	))
	if err != nil {
		logger.Logger.Error("[Update Index] failed to update index", zap.Error(err))
		return
	}

	logger.Logger.Info("[Update Index] success to update index", zap.String("Index", userIndex), zap.Int("Number", len(users)), zap.String("Status", string(info.Status)))
	return
}
