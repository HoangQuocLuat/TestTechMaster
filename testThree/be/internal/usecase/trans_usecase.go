package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testThreeBe/internal/entity"
	model "testThreeBe/internal/models"
	"testThreeBe/internal/repository"
	"testThreeBe/internal/service"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type TransUseCase struct {
	DB                   *gorm.DB
	Viper                *viper.Viper
	Log                  *logrus.Logger
	GroqService          *service.GroqService
	WebSocket            *service.WebSocketService
	DialogRepository     *repository.DialogRepository
	WordRepository       *repository.WordRepository
	WordDialogRepository *repository.WordDialogRepository
}

func NewTransUseCase(db *gorm.DB, viper *viper.Viper, log *logrus.Logger, groqService *service.GroqService,
	dialogRepository *repository.DialogRepository, wordRepository *repository.WordRepository,
	worDialogRepository *repository.WordDialogRepository, wbSocket *service.WebSocketService) *TransUseCase {
	return &TransUseCase{
		DB:                   db,
		Viper:                viper,
		Log:                  log,
		GroqService:          groqService,
		WebSocket:            wbSocket,
		DialogRepository:     dialogRepository,
		WordRepository:       wordRepository,
		WordDialogRepository: worDialogRepository,
	}
}

func (c *TransUseCase) Trans(ctx context.Context, request *model.ChatRequest) (string, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Gọi API tạo hội thoại
	content, err := c.GroqService.CallGroqAPIGetContent(request.Message)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	// 2. Lưu hội thoại vào DB
	dialog := &entity.Dialog{
		Lang:    "vi",
		Content: content,
	}
	if err := c.DialogRepository.Create(tx, dialog); err != nil {
		c.Log.WithError(err).Error("failed to create dialog")
		tx.Rollback()
		return "", err
	}
	// Gửi thông báo trạng thái qua WebSocket
	type WSMessage struct {
		Content string `json:"content"`
	}

	msg, _ := json.Marshal(WSMessage{
		Content: "Đang tiến hành trích xuất và dịch...",
	})
	c.WebSocket.BroadcastMessage(msg)

	// 3. Trích xuất từ khóa
	keywords, err := c.extractKeywords(content)
	if err != nil {
		return "", err
	}

	// 4. Dịch từ
	_, translatedWords, err := c.translateWords(keywords)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	// 5. Lưu từ và mối quan hệ
	for _, word := range translatedWords {
		wordEntity := &entity.Word{
			Lang:      "en",
			Content:   word["vi"],
			Translate: word["en"],
		}
		c.Log.Infof("Translated word: %+v", translatedWords)
		c.Log.Infof("Extracted keywords: %v", wordEntity)
		if err := c.WordRepository.Create(tx, wordEntity); err != nil {
			c.Log.WithError(err).Error("failed to create word")
			tx.Rollback()
			return "", err
		}

		wordDialog := &entity.WordDialog{
			DialogID: dialog.ID,
			WordID:   wordEntity.ID,
		}
		if err := c.WordDialogRepository.Create(tx, wordDialog); err != nil {
			c.Log.WithError(err).Error("failed to create word-dialog relationship")
			tx.Rollback()
			return "", err
		}
	}

	//Commit transaction khi tất cả thành công
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("failed to commit transaction")
		tx.Rollback()
		return "", err
	}

	msg2, _ := json.Marshal(WSMessage{
		Content: "Hoàn thành...",
	})
	c.WebSocket.BroadcastMessage(msg2)

	return "Success", nil
}

// trich xuat tu vung
func (c *TransUseCase) extractKeywords(dialogContent string) ([]string, error) {
	prompt := fmt.Sprintf(`Từ hội thoại "%s", hãy lọc ra danh sách từ quan trọng, 
bỏ qua danh từ tên riêng cần học. Không cần giải thích, xuất kết quả ra dạng JSON trong thẻ "words".`, dialogContent)

	result, err := c.GroqService.CallGroqAPIGetContent(prompt)
	if err != nil {
		return nil, err
	}

	// Loại bỏ phần ```json và ``` nếu có
	cleanedResult := strings.TrimSpace(result)
	cleanedResult = strings.TrimPrefix(cleanedResult, "```json")
	cleanedResult = strings.TrimSuffix(cleanedResult, "```")

	var words struct {
		Words []string `json:"words"`
	}
	if err := json.Unmarshal([]byte(cleanedResult), &words); err != nil {
		c.Log.WithError(err).Error("Failed to parse JSON response from Groq API")
		return nil, err
	}

	return words.Words, nil

}

// dich tieng anh
func (c *TransUseCase) translateWords(words []string) (string, []map[string]string, error) {
	wordsJSON, err := json.Marshal(words)
	if err != nil {
		c.Log.WithError(err).Error("Failed to convert words to JSON")
		return "", nil, err
	}

	prompt := fmt.Sprintf(`Dịch từng từ trong danh sách words %s sang tiếng Anh rồi trả JSON 
				gồm mảng trong đó mỗi phần tử sẽ gồm từ tiếng Việt và từ 
				tiếng Anh tương đương. Không cần giải thích.`, wordsJSON)

	result, err := c.GroqService.CallGroqAPIGetContent(prompt)
	if err != nil {
		return "", nil, err
	}

	// Loại bỏ phần ```json và ``` nếu có
	cleanedResult := strings.TrimSpace(result)
	cleanedResult = strings.TrimPrefix(cleanedResult, "```json")
	cleanedResult = strings.TrimSuffix(cleanedResult, "```")

	var translations []map[string]string
	if err := json.Unmarshal([]byte(cleanedResult), &translations); err != nil {
		c.Log.WithError(err).Error("Failed to parse JSON response")
		return "", nil, err
	}

	c.Log.Infof("Translated words: %v", translations)
	return result, translations, nil
}
