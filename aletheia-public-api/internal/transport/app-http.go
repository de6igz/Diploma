// GENERATED BY 'T'ransport 'G'enerator. DO NOT EDIT.
package transport

import (
	"aletheia-public-api/interfaces"

	"github.com/gofiber/fiber/v2"
)

type httpApp struct {
	errorHandler     ErrorHandler
	maxBatchSize     int
	maxParallelBatch int
	svc              *serverApp
	base             interfaces.App
}

func NewApp(svcApp interfaces.App) (srv *httpApp) {

	srv = &httpApp{
		base: svcApp,
		svc:  newServerApp(svcApp),
	}
	return
}

func (http *httpApp) Service() *serverApp {
	return http.svc
}

func (http *httpApp) WithLog() *httpApp {
	http.svc.WithLog()
	return http
}

func (http *httpApp) WithMetrics() *httpApp {
	http.svc.WithMetrics()
	return http
}

func (http *httpApp) WithErrorHandler(handler ErrorHandler) *httpApp {
	http.errorHandler = handler
	return http
}

func (http *httpApp) SetRoutes(route *fiber.App) {
	route.Get("/v1/me", http.serveGetMe)
}
