package search

import (
	"strconv"

	"github.com/hcd233/Aris-blog/internal/logger"
	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"
)

// QueryTagFromIndex 查询标签索引
//
//	@param query string
//	@return []map[string]interface{}
//	@return error
//	@author centonhuang
//	@update 2024-09-18 12:52:51
func QueryTagFromIndex(query string, limit int, offset int) ([]interface{}, error) {
	response, err := ServiceManager.Index(tagIndex).Search(query, &meilisearch.SearchRequest{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		logger.Logger.Error("[Query Tag Index] failed to query tag index", zap.Error(err))
		return nil, err
	}

	return response.Hits, nil
}

// AddTagIntoIndex 添加标签索引
//
//	@param tag map[string]interface{}
//	@return error
//	@author centonhuang
//	@update 2024-09-18 01:37:36
func AddTagIntoIndex(tag map[string]interface{}) error {
	info, err := ServiceManager.Index(tagIndex).AddDocuments([]map[string]interface{}{tag})
	if err != nil {
		logger.Logger.Error("[Add Tag Index] failed to add tag index", zap.Error(err))
		return err
	}

	logger.Logger.Info("[Add Tag Index] success to add tag index", zap.Any("Tag", tag), zap.String("Status", string(info.Status)))
	return nil
}

// UpdateTagInIndex 更新标签索引
//
//	@param tag map[string]interface{}
//	@return error
//	@author centonhuang
//	@update 2024-09-22 05:33:18
func UpdateTagInIndex(tag map[string]interface{}) error {
	info, err := ServiceManager.Index(tagIndex).UpdateDocuments([]map[string]interface{}{tag})
	if err != nil {
		logger.Logger.Error("[Update Tag Index] failed to update tag index", zap.Error(err))
		return err
	}

	logger.Logger.Info("[Update Tag Index] success to update tag index", zap.Any("Tag", tag), zap.String("Status", string(info.Status)))
	return nil
}

// DeleteTagFromIndex 删除标签索引
//
//	@param tagID uint
//	@return error
//	@author centonhuang
//	@update 2024-09-22 05:44:51
func DeleteTagFromIndex(tagID uint) error {
	info, err := ServiceManager.Index(tagIndex).DeleteDocument(strconv.Itoa(int(tagID)))
	if err != nil {
		logger.Logger.Error("[Delete Tag From Index] failed to delete tag from index", zap.Error(err))
		return err
	}

	logger.Logger.Info("[Delete Tag From Index] success to delete tag from index", zap.Any("TagID", tagID), zap.String("Status", string(info.Status)))
	return nil
}
