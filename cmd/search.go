package cmd

import (
	"github.com/hcd233/Aris-blog/internal/resource/search"
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
		lo.Must0(search.CreateIndex())
	},
}

var deleteIndexCmd = &cobra.Command{
	Use:   "delete",
	Short: "删除索引",
	Long:  `删除一个已有的索引，包括索引中的所有数据。`,
	Run: func(cmd *cobra.Command, args []string) {
		lo.Must0(search.DeleteIndex())
	},
}

func init() {
	searchCmd.AddCommand(createIndexCmd)
	searchCmd.AddCommand(deleteIndexCmd)
	rootCmd.AddCommand(searchCmd)

	createIndexCmd.Flags().StringP("uid", "", "", "索引的唯一标识符")
	createIndexCmd.Flags().StringP("primaryKey", "", "", "索引的主键字段")

	createIndexCmd.MarkFlagRequired("uid")
}
