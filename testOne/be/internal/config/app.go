package config

import (
	http "testOneBe/internal/delivery"
	"testOneBe/internal/delivery/middleware"
	"testOneBe/internal/delivery/route"
	"testOneBe/internal/usecase"

	"github.com/kataras/iris/v12"
	"github.com/spf13/viper"
)

type BootstrapConfig struct {
	App   *iris.Application
	Viper *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	// Setup use cases
	chatUseCase := usecase.NewChatUseCase(config.Viper)

	// Setup controllers
	chatController := http.NewChatController(chatUseCase)

	// Setup middlewares
	corsMiddleware := middleware.NewCORSMiddleware()

	// Setup routes
	routeConfig := route.RouteConfig{
		App:            config.App,
		ChatController: chatController,
		CorsMiddleware: corsMiddleware,
	}
	routeConfig.Setup()

	config.App.Logger().Info("Iris application bootstrap completed")
}
