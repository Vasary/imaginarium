package main

import (
	"context"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	e := echo.New()
	e.HideBanner = true
	e.Debug = true

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})

	exporter := echo.New()
	exporter.HideBanner = true
	prom := prometheus.NewPrometheus("echo", nil)

	e.Use(prom.HandlerFunc)
	prom.SetMetricsPath(exporter)

	go func() {
		if err := exporter.Start(":9360"); err != nil && err != http.ErrServerClosed {
			exporter.Logger.Fatal("Shutting down the server")
		}
	}()

	go func() {
		if err := e.Start(":80"); err != nil && err != http.ErrServerClosed {
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
