package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	model "testOneBe/internal/models"

	"github.com/gomarkdown/markdown"
	"github.com/spf13/viper"
)

type ChatUseCase struct {
	APIKey  string
	URLChat string
}

func NewChatUseCase(config *viper.Viper) *ChatUseCase {
	return &ChatUseCase{
		APIKey:  config.GetString("GROQ_API_KEY"),
		URLChat: config.GetString("GROQ_URL_CHAT"),
	}
}

func (c *ChatUseCase) Chat(ctx context.Context, request *model.ChatRequest) (string, error) {
	// Chuẩn bị payload
	payload := map[string]interface{}{
		"model": "deepseek-r1-distill-llama-70b",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": request.Message,
			},
		},
	}

	body, _ := json.Marshal(payload)

	// Gửi request đến Groq API
	req, _ := http.NewRequest("POST", c.URLChat, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Đọc response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Kiểm tra status code
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Groq API error: " + string(respBody))
	}

	// Parse JSON response
	var groqRes model.ChatResponse
	if err := json.Unmarshal(respBody, &groqRes); err != nil {
		return "", err
	}

	// Lấy nội dung tin nhắn
	if len(groqRes.Choices) == 0 {
		return "", errors.New("Groq API returned empty response")
	}
	responseText := groqRes.Choices[0].Message.Content

	// Convert Markdown to HTML
	htmlOutput := markdown.ToHTML([]byte(responseText), nil, nil)

	return string(htmlOutput), nil
}
