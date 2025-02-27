package http

import (
	model "testOneBe/internal/models"
	"testOneBe/internal/usecase"

	"github.com/kataras/iris/v12"
)

type ChatController struct {
	UseCase *usecase.ChatUseCase
}

func NewChatController(useCase *usecase.ChatUseCase) *ChatController {
	return &ChatController{UseCase: useCase}
}

func (c *ChatController) Chat(ctx iris.Context) {
	var request model.ChatRequest
	if err := ctx.ReadJSON(&request); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid request"})
		return
	}

	response, err := c.UseCase.Chat(ctx.Request().Context(), &request)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.JSON(response)
}
