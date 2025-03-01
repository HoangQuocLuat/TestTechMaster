package repository

import (
	"testThreeBe/internal/entity"

	"github.com/sirupsen/logrus"
)

type DialogRepository struct {
	Repository[entity.Dialog]
	Log *logrus.Logger
}

func NewDialogRepository(log *logrus.Logger) *DialogRepository {
	return &DialogRepository{
		Log: log,
	}
}
