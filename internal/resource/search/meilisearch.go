// Package search 搜索中间件
//
//	@update 2024-09-17 11:47:23
package search

import (
	"fmt"

	"github.com/hcd233/Aris-blog/internal/config"
	"github.com/hcd233/Aris-blog/internal/logger"
	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"
)

const (
	primaryKey   = "id"
	userIndex    = "user"
	tagIndex     = "tag"
	articleIndex = "article"
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
		createTagIndex,
		createArticleIndex,
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
	tasks := []func() error{
		deleteUserIndex,
		deleteTagIndex,
		deleteArticleIndex,
	}
	for _, task := range tasks {
		err = task()
		if err != nil {
			return
		}
	}
	return
}

func deleteUserIndex() (err error) {
	info, err := ServiceManager.DeleteIndex(userIndex)
	if err != nil {
		logger.Logger.Error("[Delete Index] failed to delete index", zap.Error(err))
		return
	}

	logger.Logger.Info("[Delete Index] success to delete index", zap.String("Index", userIndex), zap.String("Status", string(info.Status)))
	return
}

func deleteTagIndex() (err error) {
	info, err := ServiceManager.DeleteIndex(tagIndex)
	if err != nil {
		logger.Logger.Error("[Delete Tag Index] failed to delete index", zap.Error(err))
		return
	}

	logger.Logger.Info("[Delete Tag Index] success to delete index", zap.String("Index", tagIndex), zap.String("Status", string(info.Status)))
	return
}

func deleteArticleIndex() (err error) {
	info, err := ServiceManager.DeleteIndex(articleIndex)
	if err != nil {
		logger.Logger.Error("[Delete Article Index] failed to delete index", zap.Error(err))
		return
	}

	logger.Logger.Info("[Delete Article Index] success to delete index", zap.String("Index", articleIndex), zap.String("Status", string(info.Status)))
	return
}

func createUserIndex() (err error) {
	info, err := ServiceManager.CreateIndex(&meilisearch.IndexConfig{
		Uid:        userIndex,
		PrimaryKey: primaryKey,
	})
	if err != nil {
		logger.Logger.Error("[Create Index] failed to create index", zap.Error(err))
		return
	}

	logger.Logger.Info("[Create Index] success to create index", zap.String("Index", userIndex), zap.String("Status", string(info.Status)))
	return
}

func createTagIndex() (err error) {
	info, err := ServiceManager.CreateIndex(&meilisearch.IndexConfig{
		Uid:        tagIndex,
		PrimaryKey: primaryKey,
	})
	if err != nil {
		logger.Logger.Error("[Create Tag Index] failed to create index", zap.Error(err))
		return
	}

	logger.Logger.Info("[Create Tag Index] success to create index", zap.String("Index", "tag"), zap.String("Status", string(info.Status)))
	return
}

func createArticleIndex() (err error) {
	info, err := ServiceManager.CreateIndex(&meilisearch.IndexConfig{
		Uid:        articleIndex,
		PrimaryKey: primaryKey,
	})
	if err != nil {
		logger.Logger.Error("[Create Article Index] failed to create index", zap.Error(err))
		return
	}

	logger.Logger.Info("[Create Article Index] success to create index", zap.String("Index", "article"), zap.String("Status", string(info.Status)))
	return
}
