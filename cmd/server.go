package cmd

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-AI-go/internal/config"
	"github.com/hcd233/Aris-AI-go/internal/middleware"
	"github.com/hcd233/Aris-AI-go/internal/router"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "服务器命令组",
	Long:  `包含服务器相关操作的命令组`,
}

var startServerCmd = &cobra.Command{
	Use:   "start",
	Short: "启动API服务器",
	Long:  `启动并运行API服务器，监听指定的主机和端口`,
	Run: func(cmd *cobra.Command, args []string) {
		host, port := lo.Must1(cmd.Flags().GetString("host")), lo.Must1(cmd.Flags().GetString("port"))

		r := gin.New()
		r.Use(middleware.LoggerMiddleware(), gin.Recovery())
		router.InitRouter(r)

		s := &http.Server{
			Addr:           fmt.Sprintf("%s:%s", host, port),
			Handler:        r,
			ReadTimeout:    config.ReadTimeout,
			WriteTimeout:   config.WriteTimeout,
			MaxHeaderBytes: config.MaxHeaderBytes,
		}
		lo.Must0(s.ListenAndServe())
	},
}

func init() {
	serverCmd.AddCommand(startServerCmd)
	rootCmd.AddCommand(serverCmd)

	startServerCmd.Flags().StringP("host", "", "localhost", "监听的主机")
	startServerCmd.Flags().StringP("port", "p", "8080", "监听的端口")
}
