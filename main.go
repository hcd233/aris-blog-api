// Package main Main entry point for the Aris AI Go application.
//
//	@update 2024-10-30 08:50:53
package main

import (
	"github.com/hcd233/aris-blog-api/cmd"
)

// @title           Aris-blog
// @version         1.0
// @description     Aris-blog API

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host           9.134.115.68:8170
// @BasePath       /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cmd.Execute()
}
