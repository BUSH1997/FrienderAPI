package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

var (
	router = echo.New()
)

func main() {
	fmt.Println("Hello world")

	router.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	router.Logger.Fatal(router.Start(":8080"))
}
