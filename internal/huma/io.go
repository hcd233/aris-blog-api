package humadto

// Output 通用响应包装（Huma 风格），用于 huma 操作返回值
//
//	author centonhuang
//	update 2025-10-28 00:00:00
type Output[T any] struct {
	Body T `json:"body"`
}

// BodyInput 通用请求体包装（Huma 风格），用于 huma 操作输入值
// 说明：在 Huma 中，输入可包含 Body/Path/Query/Header 等多个部分，这里仅封装 Body，
// 其余部分通过对应的 Path/Query DTO 搭配使用
//
//	author centonhuang
//	update 2025-10-28 00:00:00
type BodyInput[T any] struct {
	Body T
}
