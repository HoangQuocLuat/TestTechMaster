package repository

import (
	"testThreeBe/internal/entity"

	"github.com/sirupsen/logrus"
)

type WordRepository struct {
	Repository[entity.Word]
	Log *logrus.Logger
}

func NewWordRepository(log *logrus.Logger) *WordRepository {
	return &WordRepository{
		Log: log,
	}
}
