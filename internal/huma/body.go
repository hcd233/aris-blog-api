package humadto

import "github.com/hcd233/aris-blog-api/internal/resource/database/model"

// ========== 通用/用户 ==========

type RefreshTokenBody struct {
	RefreshToken string `json:"refreshToken"`
}

type UpdateUserBody struct {
	UserName string `json:"userName"`
}

// ========== 文章 ==========

type CreateArticleBody struct {
	Title      string   `json:"title"`
	Slug       string   `json:"slug"`
	Tags       []string `json:"tags"`
	CategoryID uint     `json:"categoryID"`
}

type UpdateArticleBody struct {
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	CategoryID uint   `json:"categoryID"`
}

type CreateArticleVersionBody struct {
	Content string `json:"content"`
}

type UpdateArticleStatusBody struct {
	Status model.ArticleStatus `json:"status"`
}

type CreateArticleCommentBody struct {
	ArticleID uint   `json:"articleID"`
	ReplyTo   uint   `json:"replyTo"`
	Content   string `json:"content"`
}

// ========== 点赞 / 交互 ==========

type LikeBody struct {
	Undo bool `json:"undo"`
}

type LikeArticleBody struct {
	LikeBody
	ArticleID uint `json:"articleID"`
}

type LikeCommentBody struct {
	LikeBody
	CommentID uint `json:"commentID"`
}

type LikeTagBody struct {
	LikeBody
	TagID uint `json:"tagID"`
}

type LogUserViewArticleBody struct {
	ArticleID uint `json:"articleID"`
	Progress  int8 `json:"progress"`
}

// ========== 标签 ==========

type CreateTagBody struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type UpdateTagBody struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// ========== 分类 ==========

type CreateCategoryBody struct {
	ParentID uint   `json:"parentID"`
	Name     string `json:"name"`
}

type UpdateCategoryBody struct {
	Name     string `json:"name"`
	ParentID uint   `json:"parentID"`
}

// ========== 提示词 & AI ==========

type Template struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CreatePromptBody struct {
	Templates []Template `json:"templates"`
}

type AIAPPRequestBody struct {
	Temperature float32 `json:"temperature"`
}

type GenerateContentCompletionBody struct {
	AIAPPRequestBody
	Context     string `json:"context"`
	Instruction string `json:"instruction"`
	Reference   string `json:"reference"`
}

type GenerateArticleSummaryBody struct {
	AIAPPRequestBody
	ArticleID   uint   `json:"articleID"`
	Instruction string `json:"instruction"`
}

type GenerateArticleQABody struct {
	AIAPPRequestBody
	ArticleID uint   `json:"articleID"`
	Question  string `json:"question"`
}

type GenerateTermExplainationBody struct {
	AIAPPRequestBody
	ArticleID uint   `json:"articleID"`
	Term      string `json:"term"`
	Position  uint   `json:"position"`
}

// ========== Huma 输入包装 ==========

type (
	RefreshTokenInput              struct{ Body RefreshTokenBody }
	UpdateUserInput                struct{ Body UpdateUserBody }
	CreateArticleInput             struct{ Body CreateArticleBody }
	UpdateArticleInput             struct{ Body UpdateArticleBody }
	CreateArticleVersionInput      struct{ Body CreateArticleVersionBody }
	UpdateArticleStatusInput       struct{ Body UpdateArticleStatusBody }
	CreateArticleCommentInput      struct{ Body CreateArticleCommentBody }
	LikeArticleInput               struct{ Body LikeArticleBody }
	LikeCommentInput               struct{ Body LikeCommentBody }
	LikeTagInput                   struct{ Body LikeTagBody }
	LogUserViewArticleInput        struct{ Body LogUserViewArticleBody }
	CreateTagInput                 struct{ Body CreateTagBody }
	UpdateTagInput                 struct{ Body UpdateTagBody }
	CreateCategoryInput            struct{ Body CreateCategoryBody }
	UpdateCategoryInput            struct{ Body UpdateCategoryBody }
	CreatePromptInput              struct{ Body CreatePromptBody }
	GenerateContentCompletionInput struct{ Body GenerateContentCompletionBody }
	GenerateArticleSummaryInput    struct{ Body GenerateArticleSummaryBody }
	GenerateArticleQAInput         struct{ Body GenerateArticleQABody }
	GenerateTermExplainationInput  struct{ Body GenerateTermExplainationBody }
)
