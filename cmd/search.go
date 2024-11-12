package cmd

import (
	"github.com/hcd233/Aris-blog/internal/resource/search"
	doc_dao "github.com/hcd233/Aris-blog/internal/resource/search/doc_dao"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "搜索引擎相关命令组",
	Long:  `提供一组用于管理和操作搜索引擎的命令，包括创建索引等功能。`,
}

var createIndexCmd = &cobra.Command{
	Use:   "create",
	Short: "创建索引",
	Long:  `创建一个新的索引，用于存储和搜索数据。`,
	Run: func(cmd *cobra.Command, args []string) {
		search.InitSearchEngine()
		daoArr := []doc_dao.DocDAO{
			doc_dao.GetUserDocDAO(),
			doc_dao.GetTagDocDAO(),
			doc_dao.GetArticleDocDAO(),
		}

		for _, dao := range daoArr {
			lo.Must0(dao.CreateIndex())
			lo.Must0(dao.SetFilterableAttributes())
		}
	},
}

var deleteIndexCmd = &cobra.Command{
	Use:   "delete",
	Short: "删除索引",
	Long:  `删除一个已有的索引，包括索引中的所有数据。`,
	Run: func(cmd *cobra.Command, args []string) {
		search.InitSearchEngine()
		daoArr := []doc_dao.DocDAO{
			doc_dao.GetUserDocDAO(),
			doc_dao.GetTagDocDAO(),
			doc_dao.GetArticleDocDAO(),
		}

		for _, dao := range daoArr {
			lo.Must0(dao.DeleteIndex())
		}
	},
}

func init() {
	searchCmd.AddCommand(createIndexCmd)
	searchCmd.AddCommand(deleteIndexCmd)
	rootCmd.AddCommand(searchCmd)
}
