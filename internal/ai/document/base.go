// Package document 文档类
//
//	@update 2024-12-07 14:11:15
package document

// Document 文档接口
//
//	@author centonhuang
//	@update 2024-12-07 14:11:32
type Document interface {
	Map() map[string]interface{}
	String() string
}
