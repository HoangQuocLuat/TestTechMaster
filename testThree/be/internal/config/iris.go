package config

import (
	"github.com/kataras/iris/v12"
	"github.com/spf13/viper"
)

func NewIris(config *viper.Viper) *iris.Application {
	app := iris.New()
	app.Configure(iris.WithConfiguration(iris.Configuration{
		VHost: config.GetString("app.name"),
	}))
	return app
}
