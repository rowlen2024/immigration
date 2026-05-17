package handler

import (
	"net/http"

	"mygo-immigration/backend/internal/dto"
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
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(items))
}

func (h *Handler) ListAllTestimonials(c *gin.Context) {
	items, err := h.svc.Testimonial.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
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
	created, err := h.svc.Testimonial.Create(projectID, &item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
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
	var item model.Testimonial
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}
	updated, err := h.svc.Testimonial.Update(tid, &item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
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
	if err := h.svc.Testimonial.HardDelete(tid); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(nil))
}
