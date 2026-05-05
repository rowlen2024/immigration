package handler

import (
	"mygo-immigration/backend/internal/config"
	"mygo-immigration/backend/internal/repository"
	"mygo-immigration/backend/internal/service"

	"gorm.io/gorm"
)

type Handler struct {
	svc *service.Service
}

func New(db *gorm.DB, cfg *config.Config) *Handler {
	repo := repository.New(db)
	svc := service.New(repo, cfg)
	return &Handler{svc: svc}
}
