package cmd

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/resource/storage"
	objdao "github.com/hcd233/aris-blog-api/internal/resource/storage/obj_dao"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var objectCmd = &cobra.Command{
	Use:   "object",
	Short: "对象存储相关命令组",
	Long:  `提供一组用于管理和操作对象存储的命令，包括创建桶、创建目录、上传文件等。`,
}

var bucketCmd = &cobra.Command{
	Use:   "bucket",
	Short: "桶相关命令组",
	Long:  `提供一组用于管理和操作桶的命令，包括创建桶、删除桶等。`,
}

var createBucketCmd = &cobra.Command{
	Use:   "create",
	Short: "创建桶",
	Long:  `创建桶。`,
	Run: func(_ *cobra.Command, _ []string) {
		ctx := context.Background()
		logger := logger.Logger()
		storage.InitObjectStorage()

		imageObjDAO := objdao.GetImageObjDAO()
		lo.Must0(imageObjDAO.CreateBucket(ctx))

		logger.Info("[Object Storage] Bucket created",
			zap.String("bucket", imageObjDAO.GetBucketName(ctx)))

		thumbnailObjDAO := objdao.GetThumbnailObjDAO()
		lo.Must0(thumbnailObjDAO.CreateBucket(ctx))

		logger.Info("[Object Storage] Bucket created",
			zap.String("bucket", thumbnailObjDAO.GetBucketName(ctx)))
	},
}

func init() {
	bucketCmd.AddCommand(createBucketCmd)
	objectCmd.AddCommand(bucketCmd)
	rootCmd.AddCommand(objectCmd)
}
