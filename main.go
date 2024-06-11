// Description: Main entry point for the Aris AI Go application.
package main

import (
	"net/http"
	"time"

	"Aris-AI-go/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	router.StartupRouter(r)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
