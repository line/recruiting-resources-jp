// Package main starts application
package main

import (
	"todo-example/server/cmd/apiserver/app"
	"todo-example/server/pkg/config"
)

func main() {
	c := config.GetConf()
	server := &app.App{}
	server.Initialize(c)
	server.Run()
}
