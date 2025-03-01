package http

import (
	model "testThreeBe/internal/models"
	"testThreeBe/internal/service"
	"testThreeBe/internal/usecase"

	"github.com/kataras/iris/v12"
)

type TransController struct {
	UseCase   *usecase.TransUseCase
	WebSocket *service.WebSocketService
}

func NewTransController(useCase *usecase.TransUseCase, ws *service.WebSocketService) *TransController {
	return &TransController{UseCase: useCase, WebSocket: ws}
}

func (c *TransController) Trans(ctx iris.Context) {
	var request model.ChatRequest
	if err := ctx.ReadJSON(&request); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid request"})
		return
	}
	_, err := c.UseCase.Trans(ctx, &request)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}
}
