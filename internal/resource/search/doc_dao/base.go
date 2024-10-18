// Package docdao 文档DAO接口s
package docdao

import (
	"encoding/json"
	"strconv"

	"github.com/hcd233/Aris-blog/internal/logger"
	"github.com/meilisearch/meilisearch-go"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

// DocDAO 文档DAO接口
//
//	@author centonhuang
//	@update 2024-10-18 01:38:58
type DocDAO interface {
	CreateIndex(client meilisearch.ServiceManager) error
	DeleteIndex(client meilisearch.ServiceManager) error
}

// BaseDocDAO 基础文档DAO
//
//	@author centonhuang
//	@update 2024-10-17 10:40:45
type BaseDocDAO[T interface{}] struct {
	IndexName string
}

// CreateIndex 创建索引
//
//	@param dao *BaseDocDAO[T]
//	@return CreateIndex
//	@author centonhuang
//	@update 2024-10-18 01:12:47
func (dao *BaseDocDAO[T]) CreateIndex(client meilisearch.ServiceManager) error {
	taskInfo, err := client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        dao.IndexName,
		PrimaryKey: "id",
	})
	if err != nil {
		logger.Logger.Error("[Create Index]",
			zap.String("indexName", dao.IndexName),
			zap.Error(err),
		)
		return err
	}

	logger.Logger.Info("[Create Index]",
		zap.String("taskType", string(taskInfo.Type)),
		zap.Int64("taskUID", taskInfo.TaskUID),
		zap.String("indexUID", taskInfo.IndexUID),
		zap.String("status", string(taskInfo.Status)),
	)
	return nil
}

// DeleteIndex 删除索引
//
//	@param dao *BaseDocDAO[T]
//	@return DeleteIndex
//	@author centonhuang
//	@update 2024-10-18 01:13:11
func (dao *BaseDocDAO[T]) DeleteIndex(client meilisearch.ServiceManager) error {
	taskInfo, err := client.DeleteIndex(dao.IndexName)
	if err != nil {
		logger.Logger.Error("[Delete Index]",
			zap.String("indexName", dao.IndexName),
			zap.Error(err),
		)
		return err
	}

	logger.Logger.Info("[Delete Index]",
		zap.String("taskType", string(taskInfo.Type)),
		zap.Int64("taskUID", taskInfo.TaskUID),
		zap.String("indexUID", taskInfo.IndexUID),
		zap.String("status", string(taskInfo.Status)),
	)
	return nil
}

// QueryDocument 查询文档
//
//	@param dao *BaseDocDAO[T]
//	@return QueryDocument
//	@author centonhuang
//	@update 2024-10-17 10:40:43
func (dao *BaseDocDAO[T]) QueryDocument(client meilisearch.ServiceManager, query string, limit int, offset int) ([]T, error) {
	searchRequest := &meilisearch.SearchRequest{
		Query:  query,
		Limit:  int64(limit),
		Offset: int64(offset),
	}
	searchResponse, err := client.Index(dao.IndexName).Search(query, searchRequest)
	if err != nil {
		logger.Logger.Error("[Query Document]",
			zap.String("indexName", dao.IndexName),
			zap.Error(err),
		)
		return nil, err
	}
	logger.Logger.Info("[Query Document]",
		zap.String("IndexUID", searchResponse.IndexUID),
		zap.Int64("TotalHits", searchResponse.TotalHits),
		zap.Int64("ProcessingTimeMs", searchResponse.ProcessingTimeMs),
	)
	docs := lo.Map(searchResponse.Hits, func(hit interface{}, _ int) T {
		doc := hit.(map[string]interface{})
		docBytes := lo.Must1(json.Marshal(doc))
		var data T
		json.Unmarshal(docBytes, &data)
		return data
	})
	return docs, nil
}

// AddDocument 添加文档
//
//	@param dao *BaseDocDAO[T]
//	@return AddDocument
//	@author centonhuang
//	@update 2024-10-17 10:40:59
func (dao *BaseDocDAO[T]) AddDocument(client meilisearch.ServiceManager, doc *T) error {
	taskInfo, err := client.Index(dao.IndexName).AddDocuments([]*T{doc})
	if err != nil {
		logger.Logger.Error("[Add Document]",
			zap.String("indexName", dao.IndexName),
			zap.Error(err),
		)
		return err
	}
	logger.Logger.Info("[Add Document]",
		zap.String("taskType", string(taskInfo.Type)),
		zap.Int64("taskUID", taskInfo.TaskUID),
		zap.String("indexUID", taskInfo.IndexUID),
		zap.String("status", string(taskInfo.Status)),
	)
	return nil
}

// UpdateDocument 更新文档
//
//	@param dao *BaseDocDAO[T]
//	@return UpdateDocument
//	@author centonhuang
//	@update 2024-10-17 10:41:04
func (dao *BaseDocDAO[T]) UpdateDocument(client meilisearch.ServiceManager, doc *T) error {
	taskInfo, err := client.Index(dao.IndexName).UpdateDocuments([]*T{doc})
	if err != nil {
		logger.Logger.Error("[Update Document]",
			zap.String("indexName", dao.IndexName),
			zap.Error(err),
		)
		return err
	}
	logger.Logger.Info("[Update Document]",
		zap.String("taskType", string(taskInfo.Type)),
		zap.Int64("taskUID", taskInfo.TaskUID),
		zap.String("indexUID", taskInfo.IndexUID),
		zap.String("status", string(taskInfo.Status)),
	)
	return nil
}

// BatchUpdateDocuments 批量更新文档
//
//	@param dao *BaseDocDAO[T]
//	@return BatchUpdateDocuments
//	@author centonhuang
//	@update 2024-10-18 04:12:10
func (dao *BaseDocDAO[T]) BatchUpdateDocuments(client meilisearch.ServiceManager, docs []*T) error {
	if len(docs) == 0 {
		logger.Logger.Warn("[Batch Update Document]", zap.String("indexName", dao.IndexName), zap.String("message", "No document to update"))
		return nil
	}
	taskInfo, err := client.Index(dao.IndexName).UpdateDocuments(docs)
	if err != nil {
		logger.Logger.Error("[Update Document]",
			zap.String("indexName", dao.IndexName),
			zap.Error(err),
		)
		return err
	}
	logger.Logger.Info("[Batch Update Document]",
		zap.Int("docNum", len(docs)),
		zap.String("taskType", string(taskInfo.Type)),
		zap.Int64("taskUID", taskInfo.TaskUID),
		zap.String("indexUID", taskInfo.IndexUID),
		zap.String("status", string(taskInfo.Status)),
	)
	return nil
}

// DeleteDocument 删除文档
//
//	@param dao *BaseDocDAO[T]
//	@return DeleteDocument
//	@author centonhuang
//	@update 2024-10-17 10:41:10
func (dao *BaseDocDAO[T]) DeleteDocument(client meilisearch.ServiceManager, id uint) error {
	taskInfo, err := client.Index(dao.IndexName).DeleteDocument(strconv.FormatUint(uint64(id), 10))
	if err != nil {
		logger.Logger.Error("[Delete Document]",
			zap.String("indexName", dao.IndexName),
			zap.Error(err),
		)
		return err
	}
	logger.Logger.Info("[Delete Document]",
		zap.String("taskType", string(taskInfo.Type)),
		zap.Int64("taskUID", taskInfo.TaskUID),
		zap.String("indexUID", taskInfo.IndexUID),
		zap.String("status", string(taskInfo.Status)),
	)

	return nil
}
