package handler

import (
	"encoding/json"
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
	"mygo-immigration/backend/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
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
	dst.Close()
	src.Close()

	// Compress the file if > 5MB
	newPath, newSize, err := h.svc.Media.CompressIfLarge(savePath, file.Size)
	if err != nil {
		logging.Logger.Warn("compressIfLarge failed, using original", "error", err)
	} else {
		savePath = newPath
	}

	// Generate JPEG variants according to upload context
	baseName := safeName[:len(safeName)-len(filepath.Ext(safeName))]
	ctx := c.Query("context")
	if ctx == "" {
		ctx = c.PostForm("context")
	}
	if ctx != "" && !service.IsValidContext(ctx) {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid context: "+ctx))
		return
	}
	if ctx == "" {
		ctx = "general"
	}
	variants, err := service.GenerateVariantsForContext(savePath, baseName, ctx)
	if err != nil {
		logging.Logger.Warn("generateVariants failed", "error", err)
		variants = nil
	}

	variantsJSON, _ := json.Marshal(variants)

	mimeType := file.Header.Get("Content-Type")
	if strings.HasSuffix(savePath, ".jpg") {
		mimeType = "image/jpeg"
	}

	fileSize := uint64(file.Size)
	if newSize > 0 {
		fileSize = uint64(newSize)
	}

	mediaModel := &model.Media{
		Filename:     filepath.Base(savePath),
		OriginalName: file.Filename,
		URL:          "/uploads/" + filepath.Base(savePath),
		MimeType:     mimeType,
		SizeBytes:    fileSize,
		Variants:     datatypes.JSON(variantsJSON),
	}
	media, err := h.svc.Media.Upload(mediaModel)
	if err != nil {
		logging.Logger.Warn("business error in UploadMedia", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.Success(media))
}

func (h *Handler) ListMedia(c *gin.Context) {
	var req dto.MediaListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid query params"))
		return
	}

	mediaList, total, err := h.svc.Media.List(req)
	if err != nil {
		logging.Logger.Error("failed in ListMedia", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}

	if req.Page > 0 && req.PerPage > 0 {
		c.JSON(http.StatusOK, dto.SuccessPaginated(mediaList, req.Page, req.PerPage, total))
	} else {
		c.JSON(http.StatusOK, dto.Success(mediaList))
	}
}

func (h *Handler) DeleteMedia(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid media id"))
		return
	}

	if err := h.svc.Media.Delete(id); err != nil {
		logging.Logger.Warn("business error in DeleteMedia", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(nil))
}

// FindUnusedMedia returns all media records not referenced by any content.
func (h *Handler) FindUnusedMedia(c *gin.Context) {
	unused, err := h.svc.Media.FindUnused()
	if err != nil {
		logging.Logger.Error("failed in FindUnusedMedia", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(unused))
}

// CleanupUnusedMedia deletes the specified unused media records and their files.
func (h *Handler) CleanupUnusedMedia(c *gin.Context) {
	var req struct {
		IDs []uint64 `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, dto.Error(400, "请提供要清理的媒体ID列表"))
		return
	}

	deleted, failed, err := h.svc.Media.CleanupUnused(req.IDs)
	if err != nil {
		logging.Logger.Error("failed in CleanupUnusedMedia", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(map[string]interface{}{
		"deleted": deleted,
		"failed":  failed,
	}))
}
