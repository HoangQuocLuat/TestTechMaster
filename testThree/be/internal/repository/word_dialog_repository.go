package repository

import (
	"testThreeBe/internal/entity"

	"github.com/sirupsen/logrus"
)

type WordDialogRepository struct {
	Repository[entity.WordDialog]
	Log *logrus.Logger
}

func NewWordDialogRepository(log *logrus.Logger) *WordDialogRepository {
	return &WordDialogRepository{
		Log: log,
	}
}
