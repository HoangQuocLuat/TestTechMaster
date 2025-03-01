package main

import (
	"fmt"
	"testThreeBe/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	app := config.NewIris(viperConfig)

	config.Bootstrap(&config.BootstrapConfig{
		DB:     db,
		Logger: log,
		App:    app,
		Viper:  viperConfig,
	})

	webPort := viperConfig.GetInt("WEB_PORT")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		app.Logger().Fatal(err)
	}
}
