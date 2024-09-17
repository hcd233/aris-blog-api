package search

import (
	"github.com/hcd233/Aris-AI-go/internal/logger"
	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"
)

// QueryUserIndex 查询用户索引
//
//	@param query string
//	@return []map[string]interface{}
//	@return error
//	@author centonhuang
//	@update 2024-09-18 12:52:51
func QueryUserIndex(query string, limit int64, offset int64) ([]interface{}, error) {
	response, err := ServiceManager.Index(userIndex).Search(query, &meilisearch.SearchRequest{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		logger.Logger.Error("[Query User Index] failed to query user index", zap.Error(err))
		return nil, err
	}

	logger.Logger.Info("[Query User Index] success to query user index", zap.String("Query", query), zap.Int64("Limit", limit), zap.Int64("Offset", offset))
	return response.Hits, nil
}
