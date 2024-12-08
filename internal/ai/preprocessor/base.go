// Package preprocessor 预处理
//
//	@update 2024-12-08 14:49:58
package preprocessor

import "github.com/hcd233/Aris-blog/internal/ai/document"

// Preprocessor 预处理器
//
//	@author centonhuang
//	@update 2024-12-08 14:53:11
type Preprocessor[sourceT interface{}, documentT document.Document] interface {
	Process(source sourceT) (documents []documentT, err error)
	BatchProcess(sources []sourceT) (documents []documentT, err error)
	ProcessDocument(rawDocument documentT) (processedDocuments []documentT, err error)
	BatchProcessDocument(rawDocuments []documentT) (processedDocuments []documentT, err error)
}
