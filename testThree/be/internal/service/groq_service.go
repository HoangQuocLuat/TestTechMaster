package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	model "testThreeBe/internal/models"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type GroqService struct {
	Viper  *viper.Viper
	Logger *logrus.Logger
}

func NewGroqService(viper *viper.Viper, logger *logrus.Logger) *GroqService {
	return &GroqService{
		Viper:  viper,
		Logger: logger,
	}
}

func (g *GroqService) CallGroqAPIGetContent(message string) (string, error) {
	g.Logger.Info("Calling Groq API...")

	payload := map[string]interface{}{
		"model": "deepseek-r1-distill-llama-70b",
		"messages": []map[string]string{
			{"role": "user", "content": message},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		g.Logger.WithError(err).Error("Failed to marshal request payload")
		return "", err
	}

	url := g.Viper.GetString("GROQ_URL_CHAT")
	apiKey := g.Viper.GetString("GROQ_API_KEY")

	g.Logger.WithFields(logrus.Fields{
		"url": url,
	}).Info("Sending request to Groq API")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		g.Logger.WithError(err).Error("Failed to create HTTP request")
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		g.Logger.WithError(err).Error("HTTP request to Groq API failed")
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		g.Logger.WithError(err).Error("Failed to read response body")
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		g.Logger.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
			"response":    string(respBody),
		}).Error("Groq API returned an error")
		return "", errors.New("Groq API error: " + string(respBody))
	}

	var groqRes model.ChatResponse
	if err := json.Unmarshal(respBody, &groqRes); err != nil {
		g.Logger.WithError(err).Error("Failed to unmarshal Groq API response")
		return "", err
	}

	if len(groqRes.Choices) == 0 {
		g.Logger.Warn("Groq API returned an empty response")
		return "", errors.New("Groq API returned an empty response")
	}

	content := groqRes.Choices[0].Message.Content
	cleanedContent := removeThinkTags(content)

	g.Logger.Info("Successfully received response from Groq API")

	return cleanedContent, nil
}

func removeThinkTags(content string) string {
	re := regexp.MustCompile(`(?s)<think>.*?</think>\s*`)
	return re.ReplaceAllString(content, "")
}
