package main

import (
	"embed"
	"fmt"
	"io/fs"
	"komiko/config"
	"komiko/db"
	"komiko/router"
	"net/http"
	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
)

//go:embed dist/*
var embeddedFiles embed.FS

func main() {
	db := db.Connect()

	r := gin.Default()

	router.RegisterRoutes(r, db)

	distFS, _ := fs.Sub(embeddedFiles, "dist")
	assetsFS, _ := fs.Sub(distFS, "assets")
	r.StaticFS("/assets", http.FS(assetsFS))
	r.StaticFileFS("/favicon.ico", "favicon.ico", http.FS(distFS))
	r.NoRoute(func(c *gin.Context) {
		data, err := embeddedFiles.ReadFile("dist/index.html")
		if err != nil {
			c.String(500, "index.html not found")
			return
		}
		c.Data(200, "text/html; charset=utf-8", data)
	})

	cfg := config.GetConfig()
	if err := r.Run(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
