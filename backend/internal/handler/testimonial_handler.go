package handler

import (
	"net/http"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/logging"
	"mygo-immigration/backend/internal/model"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListProjectTestimonials(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	items, err := h.svc.Testimonial.ListByProject(projectID)
	if err != nil {
		logging.Logger.Error("failed in ListProjectTestimonials", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(items))
}

func (h *Handler) ListAllTestimonials(c *gin.Context) {
	items, err := h.svc.Testimonial.ListAll()
	if err != nil {
		logging.Logger.Error("failed in ListAllTestimonials", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(items))
}

func (h *Handler) AdminListTestimonials(c *gin.Context) {
	var items []model.Testimonial
	var err error

	if pidStr := c.Query("project_id"); pidStr != "" {
		pid, parseErr := parseIDParam(c, "project_id")
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project_id"))
			return
		}
		items, err = h.svc.Testimonial.ListByProject(pid)
	} else {
		items, err = h.svc.Testimonial.ListAll()
	}

	if err != nil {
		logging.Logger.Error("failed in AdminListTestimonials", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(items))
}

func (h *Handler) CreateProjectTestimonial(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	var item model.Testimonial
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}
	item.ID = 0 // ensure DB auto-increment
	created, err := h.svc.Testimonial.Create(projectID, &item)
	if err != nil {
		logging.Logger.Warn("business error in CreateProjectTestimonial", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.Success(created))
}

func (h *Handler) UpdateProjectTestimonial(c *gin.Context) {
	_, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	tid, err := parseIDParam(c, "tid")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid testimonial id"))
		return
	}
	var req dto.UpdateTestimonialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}
	updated, err := h.svc.Testimonial.Update(tid, req)
	if err != nil {
		logging.Logger.Warn("business error in UpdateProjectTestimonial", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(updated))
}

func (h *Handler) DeleteProjectTestimonial(c *gin.Context) {
	_, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	tid, err := parseIDParam(c, "tid")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid testimonial id"))
		return
	}
	if err := h.svc.Testimonial.Delete(tid); err != nil {
		logging.Logger.Warn("business error in DeleteProjectTestimonial", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(nil))
}
