package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"imaginarium/config"
	"imaginarium/rest"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var cfg *config.Config

func init() {
	var configPath string

	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")
	flag.Parse()

	cfg = config.NewConfig(configPath)

	rest.StoragePath = cfg.Storage.Path
	rest.FileSize = cfg.Server.Uploader.MaxSize

	rest.SupportedTypes = make(map[string]bool)
	for _, v := range cfg.Server.Uploader.Allow {
		rest.SupportedTypes[v] = true
	}
}

func main() {
	e := echo.New()
	e.HideBanner = true
	e.Debug = true

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	e.GET("/health", rest.Health)
	e.POST("/upload", rest.Upload)
	e.GET("/", func(c echo.Context) error {
		return c.HTML(200, "Image storage 0.0.1")
	})
	e.GET("/:name", rest.Download)

	exporter := echo.New()
	exporter.HideBanner = true
	prom := prometheus.NewPrometheus("imaginarium", nil)

	e.Use(prom.HandlerFunc)
	prom.SetMetricsPath(exporter)

	go func() {
		if err := exporter.Start(fmt.Sprintf(":%s", cfg.Exporter.Port)); err != nil && err != http.ErrServerClosed {
			exporter.Logger.Fatal("Shutting down the server")
		}
	}()

	go func() {
		if err := e.Start(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	if err := exporter.Shutdown(ctx); err != nil {
		exporter.Logger.Fatal(err)
	}
}
