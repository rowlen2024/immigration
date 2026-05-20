package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/logging"
	"mygo-immigration/backend/internal/model"

	"github.com/gin-gonic/gin"
)

const uploadDir = "./uploads"

func (h *Handler) UploadMedia(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "file is required"))
		return
	}

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "failed to create upload directory"))
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	safeName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	savePath := filepath.Join(uploadDir, safeName)

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "failed to open uploaded file"))
		return
	}
	defer src.Close()

	dst, err := os.Create(savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "failed to create file"))
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "failed to save file"))
		return
	}

	mediaModel := &model.Media{
		Filename:     safeName,
		OriginalName: file.Filename,
		URL:          "/uploads/" + safeName,
		MimeType:     file.Header.Get("Content-Type"),
		SizeBytes:    uint64(file.Size),
	}
	media, err := h.svc.Media.Upload(mediaModel)
	if err != nil {
		logging.Logger.Error("failed in UploadMedia", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}

	c.JSON(http.StatusCreated, dto.Success(media))
}

func (h *Handler) ListMedia(c *gin.Context) {
	page, perPage := parsePagination(c)
	search := c.Query("search")

	mediaList, total, err := h.svc.Media.List(page, perPage, search)
	if err != nil {
		logging.Logger.Error("failed in ListMedia", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessPaginated(mediaList, page, perPage, total))
}

func (h *Handler) DeleteMedia(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid media id"))
		return
	}

	if err := h.svc.Media.Delete(id); err != nil {
		logging.Logger.Error("failed in DeleteMedia", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(nil))
}
