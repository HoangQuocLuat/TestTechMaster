package config

import (
	http "testThreeBe/internal/delivery"
	"testThreeBe/internal/delivery/middleware"
	"testThreeBe/internal/delivery/route"
	"testThreeBe/internal/repository"
	"testThreeBe/internal/service"
	"testThreeBe/internal/usecase"

	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB     *gorm.DB
	Logger *logrus.Logger
	App    *iris.Application
	Viper  *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {

	// Setup repository
	dialogRepository := repository.NewDialogRepository(config.Logger)
	wordRepository := repository.NewWordRepository(config.Logger)
	wordDialogRepository := repository.NewWordDialogRepository(config.Logger)

	// Setup service groq
	groqService := service.NewGroqService(config.Viper, config.Logger)
	// Setup service websocket
	webSocketService := service.NewWebSocketService()

	// Setup use cases
	transUseCase := usecase.NewTransUseCase(config.DB, config.Viper, config.Logger,
		groqService, dialogRepository, wordRepository, wordDialogRepository, webSocketService)

	// Setup controllers
	transController := http.NewTransController(transUseCase, webSocketService)

	// Setup middlewares
	corsMiddleware := middleware.NewCORSMiddleware()

	// Setup routes
	routeConfig := route.RouteConfig{
		App:             config.App,
		TransController: transController,
		CorsMiddleware:  corsMiddleware,
		WebSocket:       webSocketService,
	}
	routeConfig.Setup()

	config.App.Logger().Info("Iris application bootstrap completed")
}
