// Package search 搜索中间件
//
//	@update 2024-09-17 11:47:23
package search

import (
	"fmt"

	"github.com/hcd233/Aris-blog/internal/config"
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
	_ = lo.Must1(serviceManager.Health())
}
