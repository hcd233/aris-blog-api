package search

import (
	"github.com/hcd233/Aris-blog/internal/logger"
	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"
)

// QueryUserFromIndex 查询用户索引
//
//	@param query string
//	@return []map[string]interface{}
//	@return error
//	@author centonhuang
//	@update 2024-09-18 12:52:51
func QueryUserFromIndex(query string, limit int, offset int) ([]interface{}, error) {
	response, err := ServiceManager.Index(userIndex).Search(query, &meilisearch.SearchRequest{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		logger.Logger.Error("[Query User Index] failed to query user index", zap.Error(err))
		return nil, err
	}

	return response.Hits, nil
}

// AddUserIntoIndex 添加用户索引
//
//	@param user map[string]interface{}
//	@return error
//	@author centonhuang
//	@update 2024-09-18 01:37:36
func AddUserIntoIndex(user map[string]interface{}) error {
	info, err := ServiceManager.Index(userIndex).AddDocuments([]map[string]interface{}{user})
	if err != nil {
		logger.Logger.Error("[Add User Index] failed to add user index", zap.Error(err))
		return err
	}

	logger.Logger.Info("[Add User Index] success to add user index", zap.Any("User", user), zap.String("Status", string(info.Status)))
	return nil
}

// UpdateUserInIndex 更新用户索引
//
//	@param user map[string]interface{}
//	@return error
//	@author centonhuang
//	@update 2024-09-18 01:41:04
func UpdateUserInIndex(user map[string]interface{}) error {
	info, err := ServiceManager.Index(userIndex).UpdateDocuments([]map[string]interface{}{user})
	if err != nil {
		logger.Logger.Error("[Update User Index] failed to update user index", zap.Error(err))
		return err
	}

	logger.Logger.Info("[Update User Index] success to update user index", zap.Any("User", user), zap.String("Status", string(info.Status)))
	return nil
}
