// Description: Main entry point for the Aris AI Go application.
package main

import (
	"fmt"
	"net/http"

	"github.com/hcd233/Aris-AI-go/internal/config"
	"github.com/hcd233/Aris-AI-go/internal/logger"
	"github.com/hcd233/Aris-AI-go/internal/resource/database"
	"github.com/hcd233/Aris-AI-go/internal/resource/database/model"
	"github.com/hcd233/Aris-AI-go/internal/router"
	"github.com/samber/lo"
)

func main() {
	logger.InitLogger()
	logger.Logger.Info("[logger] init logger successfully")
	database.InitDatabase()
	lo.Must0(database.DB.AutoMigrate(model.Models...))
	logger.Logger.Info("[database] init database successfully")
	router.InitRouter(router.Router)

	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", config.Port),
		Handler:        router.Router,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}
	s.ListenAndServe()
}
