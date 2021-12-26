package rest

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/sys/unix"
	"net/http"
)

type Status struct {
	Health bool `json:"health"`
}

func Health(c echo.Context) error {
	if writable(StoragePath) {
		return c.JSON(http.StatusOK, Status{true})
	}

	return c.JSON(http.StatusInternalServerError, Status{false})
}

func writable(path string) bool {
	return unix.Access(path, unix.W_OK) == nil
}
