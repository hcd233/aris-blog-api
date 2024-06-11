// Description: Main entry point for the Aris AI Go application.
package main

import (
	"fmt"
	"net/http"

	"github.com/hcd233/Aris-AI-go/internal/config"
	"github.com/hcd233/Aris-AI-go/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	router.StartupRouter(r)
	config.InitEnvironment()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", config.Port),
		Handler:        r,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}
	s.ListenAndServe()
}
