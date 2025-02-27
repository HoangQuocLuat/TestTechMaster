package route

import (
	http "testOneBe/internal/delivery"

	"github.com/kataras/iris/v12"
)

type RouteConfig struct {
	App            *iris.Application
	ChatController *http.ChatController
	CorsMiddleware iris.Handler
}

func (c *RouteConfig) Setup() {
	c.App.UseRouter(c.CorsMiddleware)
	api := c.App.Party("/api/v1")
	c.SetupChatRoute(api)
}

func (c *RouteConfig) SetupChatRoute(api iris.Party) {
	chat := api.Party("/chat")
	chat.Post("/", c.ChatController.Chat)
}
