// Package search 搜索中间件
//
//	@update 2024-09-17 11:47:23
package search

import (
	"fmt"

	"github.com/hcd233/Aris-blog/internal/config"
	docdao "github.com/hcd233/Aris-blog/internal/resource/search/doc_dao"
	"github.com/meilisearch/meilisearch-go"
	"github.com/samber/lo"
)

const (
	primaryKey = "id"
)

// serviceManager meilisearch 服务管理
//
//	@update 2024-09-18 12:29:08
var serviceManager meilisearch.ServiceManager

// GetSearchEngine 获取搜索引擎实例
//
//	@return meilisearch.ServiceManager
//	@author centonhuang
//	@update 2024-10-17 09:47:38
func GetSearchEngine() meilisearch.ServiceManager {
	return serviceManager
}

// InitSearchEngine 初始化搜索引擎
//
//	@author centonhuang
//	@update 2024-09-22 10:05:14
func InitSearchEngine() {
	serviceManager = meilisearch.New(
		fmt.Sprintf("http://%s:%s", config.MeilisearchHost, config.MeilisearchPort),
		meilisearch.WithAPIKey(config.MeilisearchMasterKey),
	)
	lo.Must1(serviceManager.Health())
}

// CreateIndex 创建索引
//
//	@return err error
//	@author centonhuang
//	@update 2024-09-18 12:43:35
func CreateIndex() (err error) {
	client := GetSearchEngine()

	daoArr := []docdao.DocDAO{
		docdao.GetUserDocDAO(),
		docdao.GetTagDocDAO(),
		docdao.GetArticleDocDAO(),
	}

	for _, dao := range daoArr {
		lo.Must0(dao.CreateIndex(client))
	}
	return
}

// DeleteIndex 删除索引
//
//	@return err error
//	@author centonhuang
//	@update 2024-09-18 01:24:35
func DeleteIndex() (err error) {
	client := GetSearchEngine()

	daoArr := []docdao.DocDAO{
		docdao.GetUserDocDAO(),
		docdao.GetTagDocDAO(),
		docdao.GetArticleDocDAO(),
	}

	for _, dao := range daoArr {
		lo.Must0(dao.DeleteIndex(client))
	}
	return
}
