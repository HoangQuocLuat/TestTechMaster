package route

import (
	http "testThreeBe/internal/delivery"
	"testThreeBe/internal/service"

	"github.com/kataras/iris/v12"
)

type RouteConfig struct {
	App             *iris.Application
	TransController *http.TransController
	CorsMiddleware  iris.Handler
	WebSocket       *service.WebSocketService
}

func (c *RouteConfig) Setup() {
	c.App.UseRouter(c.CorsMiddleware)
	api := c.App.Party("/api/v1")
	c.SetupChatRoute(api)
	c.SetupWebSocketRoute(api)
}

// Route WebSocket
func (c *RouteConfig) SetupWebSocketRoute(api iris.Party) {
	api.Get("/ws", func(ctx iris.Context) {
		w := ctx.ResponseWriter()
		r := ctx.Request()
		c.WebSocket.HandleConnection(w, r)
	})
}

// Route API
func (c *RouteConfig) SetupChatRoute(api iris.Party) {
	r := api.Party("/trans")
	r.Post("/", c.TransController.Trans)
}
