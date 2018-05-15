package main

import (
	"jwt-auth/api"
	"jwt-auth/service"
	"log"
	"net/http"

	"io/ioutil"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const projectID = "golang-jwt-auth"

func main() {
	serviceAccount, err := ioutil.ReadFile("service-account.json")
	if err != nil {
		log.Fatal(err)
	}
	err = api.Init(api.Config{
		ServiceAccountJSON: serviceAccount,
		ProjectID:          projectID,
	})
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(
		middleware.Recover(),
		middleware.Secure(),
		middleware.Logger(),
		middleware.Gzip(),
		middleware.BodyLimit("2M"),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{
				"http://localhost:8080",
			},
			AllowHeaders: []string{
				echo.HeaderOrigin,
				echo.HeaderContentLength,
				echo.HeaderAcceptEncoding,
				echo.HeaderContentType,
				echo.HeaderAuthorization,
			},
			AllowMethods: []string{
				echo.GET,
				echo.POST,
			},
			MaxAge: 3600,
		}),
	)

	// Health check
	e.GET("/_ah/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Register services
	service.Auth(e.Group("/auth"))

	e.Start(":9000")
}
