// Package main Main entry point for the Aris AI Go application.
//
//	@update 2024-10-30 08:50:53
package main

import (
	"github.com/hcd233/Aris-blog/cmd"
	_ "github.com/hcd233/Aris-blog/docs"
)

// @title           Aris-blog
// @version         1.0
// @description     Aris-blog API
// @host           9.134.115.68:8170
// @BasePath       /
func main() {
	cmd.Execute()
}
