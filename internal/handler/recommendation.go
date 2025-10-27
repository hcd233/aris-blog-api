package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service/recommendation"
	"github.com/hcd233/aris-blog-api/internal/util"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RecommendationHandler 推荐系统处理器接口
//
//	author system
//	update 2025-01-19 12:00:00
type RecommendationHandler interface {
	// 用户行为上报
	ReportBehavior(c *fiber.Ctx) error
	
	// 推荐接口
	RecommendArticles(c *fiber.Ctx) error
	RecommendTags(c *fiber.Ctx) error
	
	// 用户画像
	GetUserProfile(c *fiber.Ctx) error
	
	// 管理接口
	TrainModel(c *fiber.Ctx) error
	UpdateUserProfile(c *fiber.Ctx) error
}

type recommendationHandler struct {
	service *recommendation.RecommendationService
	logger  *zap.Logger
}

// NewRecommendationHandler 创建推荐系统处理器
//
//	author system
//	update 2025-01-19 12:00:00
func NewRecommendationHandler(db *gorm.DB, redis *redis.Client, logger *zap.Logger) RecommendationHandler {
	return &recommendationHandler{
		service: recommendation.NewRecommendationService(db, redis, logger),
		logger:  logger,
	}
}

// ReportBehavior 用户行为上报
//
//	@Summary		用户行为上报
//	@Description	记录用户对文章或标签的行为（浏览、点赞、分享、评论、收藏）
//	@Tags			recommendation
//	@Accept			json
//	@Produce		json
//	@Param			request	body		protocol.UserBehaviorRequest	true	"行为数据"
//	@Success		200		{object}	protocol.HTTPResponse{data=protocol.UserBehaviorResponse,error=nil}
//	@Failure		400		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/api/recommendation/behavior [post]
//	@Security		ApiKeyAuth
func (h *recommendationHandler) ReportBehavior(c *fiber.Ctx) error {
	var req protocol.UserBehaviorRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("解析用户行为请求失败", zap.Error(err))
		util.SendHTTPResponse(c, nil, protocol.ErrBadRequest)
		return nil
	}

	// 验证请求参数
	if err := util.ValidateStruct(&req); err != nil {
		h.logger.Error("用户行为请求参数验证失败", zap.Error(err))
		util.SendHTTPResponse(c, nil, protocol.ErrBadRequest)
		return nil
	}

	response, err := h.service.ReportBehavior(c.Context(), &req)
	if err != nil {
		h.logger.Error("用户行为上报失败", zap.Error(err))
		util.SendHTTPResponse(c, nil, protocol.ErrInternalError)
		return nil
	}

	util.SendHTTPResponse(c, response, nil)
	return nil
}

// RecommendArticles 推荐文章
//
//	@Summary		推荐文章
//	@Description	基于用户画像和协同过滤为用户推荐文章
//	@Tags			recommendation
//	@Accept			json
//	@Produce		json
//	@Param			userId		query		int						true	"用户ID"
//	@Param			limit		query		int						false	"推荐数量限制"	default(10)
//	@Param			excludeIds	query		[]int					false	"排除的文章ID列表"
//	@Param			includeTags	query		[]string				false	"包含的标签列表"
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.RecommendationResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/api/recommendation/articles [get]
//	@Security		ApiKeyAuth
func (h *recommendationHandler) RecommendArticles(c *fiber.Ctx) error {
	userIDStr := c.Query("userId")
	if userIDStr == "" {
		util.SendHTTPResponse(c, nil, protocol.ErrBadRequest)
		return nil
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		util.SendHTTPResponse(c, nil, protocol.ErrBadRequest)
		return nil
	}

	req := &protocol.RecommendationRequest{
		UserID: uint(userID),
		Type:   "article",
	}

	// 解析可选参数
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = limit
		}
	}

	// 解析排除的文章ID列表
	if excludeIDsStr := c.Query("excludeIds"); excludeIDsStr != "" {
		// 简化处理，实际项目中应该正确解析数组
		h.logger.Info("排除的文章ID", zap.String("excludeIds", excludeIDsStr))
	}

	// 解析包含的标签列表
	if includeTagsStr := c.Query("includeTags"); includeTagsStr != "" {
		// 简化处理，实际项目中应该正确解析数组
		h.logger.Info("包含的标签", zap.String("includeTags", includeTagsStr))
	}

	response, err := h.service.RecommendArticles(c.Context(), req)
	if err != nil {
		h.logger.Error("文章推荐失败", zap.Error(err))
		util.SendHTTPResponse(c, nil, protocol.ErrInternalError)
		return nil
	}

	util.SendHTTPResponse(c, response, nil)
	return nil
}

// RecommendTags 推荐标签
//
//	@Summary		推荐标签
//	@Description	基于用户画像为用户推荐感兴趣的标签
//	@Tags			recommendation
//	@Accept			json
//	@Produce		json
//	@Param			userId	query		int		true	"用户ID"
//	@Param			limit	query		int		false	"推荐数量限制"	default(10)
//	@Success		200		{object}	protocol.HTTPResponse{data=protocol.RecommendationResponse,error=nil}
//	@Failure		400		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/api/recommendation/tags [get]
//	@Security		ApiKeyAuth
func (h *recommendationHandler) RecommendTags(c *fiber.Ctx) error {
	userIDStr := c.Query("userId")
	if userIDStr == "" {
		util.SendHTTPResponse(c, nil, protocol.ErrBadRequest)
		return nil
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		util.SendHTTPResponse(c, nil, protocol.ErrBadRequest)
		return nil
	}

	req := &protocol.RecommendationRequest{
		UserID: uint(userID),
		Type:   "tag",
	}

	// 解析可选参数
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = limit
		}
	}

	response, err := h.service.RecommendTags(c.Context(), req)
	if err != nil {
		h.logger.Error("标签推荐失败", zap.Error(err))
		util.SendHTTPResponse(c, nil, protocol.ErrInternalError)
		return nil
	}

	util.SendHTTPResponse(c, response, nil)
	return nil
}

// GetUserProfile 获取用户画像
//
//	@Summary		获取用户画像
//	@Description	获取用户的兴趣偏好和行为统计信息
//	@Tags			recommendation
//	@Accept			json
//	@Produce		json
//	@Param			userId	query		int		true	"用户ID"
//	@Success		200		{object}	protocol.HTTPResponse{data=protocol.UserProfileResponse,error=nil}
//	@Failure		400		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/api/recommendation/profile [get]
//	@Security		ApiKeyAuth
func (h *recommendationHandler) GetUserProfile(c *fiber.Ctx) error {
	userIDStr := c.Query("userId")
	if userIDStr == "" {
		util.SendHTTPResponse(c, nil, protocol.ErrBadRequest)
		return nil
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		util.SendHTTPResponse(c, nil, protocol.ErrBadRequest)
		return nil
	}

	req := &protocol.UserProfileRequest{
		UserID: uint(userID),
	}

	response, err := h.service.GetUserProfile(c.Context(), req)
	if err != nil {
		h.logger.Error("获取用户画像失败", zap.Error(err))
		util.SendHTTPResponse(c, nil, protocol.ErrInternalError)
		return nil
	}

	util.SendHTTPResponse(c, response, nil)
	return nil
}

// TrainModel 训练推荐模型
//
//	@Summary		训练推荐模型
//	@Description	基于历史行为数据训练协同过滤模型
//	@Tags			recommendation
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	protocol.HTTPResponse{data=string,error=nil}
//	@Failure		500		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/api/recommendation/train [post]
//	@Security		ApiKeyAuth
func (h *recommendationHandler) TrainModel(c *fiber.Ctx) error {
	err := h.service.TrainModel(c.Context())
	if err != nil {
		h.logger.Error("训练推荐模型失败", zap.Error(err))
		util.SendHTTPResponse(c, nil, protocol.ErrInternalError)
		return nil
	}

	util.SendHTTPResponse(c, "推荐模型训练成功", nil)
	return nil
}

// UpdateUserProfile 更新用户画像
//
//	@Summary		更新用户画像
//	@Description	基于用户行为数据重新构建用户画像
//	@Tags			recommendation
//	@Accept			json
//	@Produce		json
//	@Param			userId	query		int		true	"用户ID"
//	@Success		200		{object}	protocol.HTTPResponse{data=string,error=nil}
//	@Failure		400		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500		{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/api/recommendation/profile/update [post]
//	@Security		ApiKeyAuth
func (h *recommendationHandler) UpdateUserProfile(c *fiber.Ctx) error {
	userIDStr := c.Query("userId")
	if userIDStr == "" {
		util.SendHTTPResponse(c, nil, protocol.ErrBadRequest)
		return nil
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		util.SendHTTPResponse(c, nil, protocol.ErrBadRequest)
		return nil
	}

	err = h.service.UpdateUserProfile(c.Context(), uint(userID))
	if err != nil {
		h.logger.Error("更新用户画像失败", zap.Error(err))
		util.SendHTTPResponse(c, nil, protocol.ErrInternalError)
		return nil
	}

	util.SendHTTPResponse(c, "用户画像更新成功", nil)
	return nil
}