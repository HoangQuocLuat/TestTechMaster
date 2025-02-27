package main

import (
	"fmt"
	"testOneBe/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	app := config.NewIris(viperConfig)

	config.Bootstrap(&config.BootstrapConfig{
		App:   app,
		Viper: viperConfig,
	})

	webPort := viperConfig.GetInt("WEB_PORT")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		app.Logger().Fatal(err)
	}
}
