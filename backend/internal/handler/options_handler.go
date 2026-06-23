package handler

import (
	"net/http"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/logging"

	"github.com/gin-gonic/gin"
)

const defaultOptionPage = 1
const defaultOptionPerPage = 500

func normalizeOptionRequest(page, perPage int) (int, int) {
	if page < 1 {
		page = defaultOptionPage
	}
	if perPage < 1 {
		perPage = defaultOptionPerPage
	}
	if perPage > defaultOptionPerPage {
		perPage = defaultOptionPerPage
	}
	return page, perPage
}

func (h *Handler) ListProjectOptions(c *gin.Context) {
	var req dto.ProjectOptionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid query params"))
		return
	}
	req.Page, req.PerPage = normalizeOptionRequest(req.Page, req.PerPage)

	items, total, err := h.svc.Project.Options(req, true)
	if err != nil {
		logging.Logger.Error("failed in ListProjectOptions", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessPaginated(items, req.Page, req.PerPage, total))
}

func (h *Handler) AdminListProjectOptions(c *gin.Context) {
	var req dto.ProjectOptionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid query params"))
		return
	}
	req.Page, req.PerPage = normalizeOptionRequest(req.Page, req.PerPage)

	items, total, err := h.svc.Project.Options(req, false)
	if err != nil {
		logging.Logger.Error("failed in AdminListProjectOptions", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessPaginated(items, req.Page, req.PerPage, total))
}

func (h *Handler) AdminListCaseOptions(c *gin.Context) {
	var req dto.CaseOptionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid query params"))
		return
	}
	req.Page, req.PerPage = normalizeOptionRequest(req.Page, req.PerPage)

	items, total, err := h.svc.Case.Options(req)
	if err != nil {
		logging.Logger.Error("failed in AdminListCaseOptions", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessPaginated(items, req.Page, req.PerPage, total))
}

func (h *Handler) AdminListTestimonialOptions(c *gin.Context) {
	var req dto.TestimonialOptionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid query params"))
		return
	}
	req.Page, req.PerPage = normalizeOptionRequest(req.Page, req.PerPage)

	items, total, err := h.svc.Testimonial.Options(req)
	if err != nil {
		logging.Logger.Error("failed in AdminListTestimonialOptions", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessPaginated(items, req.Page, req.PerPage, total))
}

func (h *Handler) AdminListPageOptions(c *gin.Context) {
	var req dto.PageOptionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid query params"))
		return
	}
	req.Page, req.PerPage = normalizeOptionRequest(req.Page, req.PerPage)

	items, total, err := h.svc.Page.Options(req)
	if err != nil {
		logging.Logger.Error("failed in AdminListPageOptions", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessPaginated(items, req.Page, req.PerPage, total))
}
