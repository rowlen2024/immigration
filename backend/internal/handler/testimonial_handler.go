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

	items, _, err := h.svc.Testimonial.List(dto.TestimonialListRequest{
		ProjectID: &projectID,
	})
	if err != nil {
		logging.Logger.Error("failed in ListProjectTestimonials", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, dto.Success(items))
}

func (h *Handler) ListAllTestimonials(c *gin.Context) {
	var req dto.TestimonialListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid query params"))
		return
	}

	items, total, err := h.svc.Testimonial.List(req)
	if err != nil {
		logging.Logger.Error("failed in ListAllTestimonials", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}

	if req.Page > 0 && req.PerPage > 0 {
		c.JSON(http.StatusOK, dto.SuccessPaginated(items, req.Page, req.PerPage, total))
	} else {
		c.JSON(http.StatusOK, dto.Success(items))
	}
}

func (h *Handler) AdminListTestimonials(c *gin.Context) {
	var req dto.TestimonialListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid query params"))
		return
	}

	items, total, err := h.svc.Testimonial.List(req)
	if err != nil {
		logging.Logger.Error("failed in AdminListTestimonials", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}

	if req.Page > 0 && req.PerPage > 0 {
		c.JSON(http.StatusOK, dto.SuccessPaginated(items, req.Page, req.PerPage, total))
	} else {
		c.JSON(http.StatusOK, dto.Success(items))
	}
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
