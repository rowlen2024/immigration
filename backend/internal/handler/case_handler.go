package handler

import (
	"net/http"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/logging"
	"mygo-immigration/backend/internal/model"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetCase(c *gin.Context) {
	slug := c.Param("slug")
	cs, err := h.svc.Case.GetBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Error(404, "case not found"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(cs))
}

func (h *Handler) ListCases(c *gin.Context) {
	var req dto.CaseListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid query params"))
		return
	}

	cases, total, err := h.svc.Case.List(req)
	if err != nil {
		logging.Logger.Error("failed in ListCases", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}

	if req.Page > 0 && req.PerPage > 0 {
		c.JSON(http.StatusOK, dto.SuccessPaginated(cases, req.Page, req.PerPage, total))
	} else {
		c.JSON(http.StatusOK, dto.Success(cases))
	}
}

func (h *Handler) AdminListCases(c *gin.Context) {
	var req dto.CaseListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid query params"))
		return
	}

	cases, total, err := h.svc.Case.List(req)
	if err != nil {
		logging.Logger.Error("failed in AdminListCases", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}

	if req.Page > 0 && req.PerPage > 0 {
		c.JSON(http.StatusOK, dto.SuccessPaginated(cases, req.Page, req.PerPage, total))
	} else {
		c.JSON(http.StatusOK, dto.Success(cases))
	}
}

func (h *Handler) CreateCase(c *gin.Context) {
	var caseModel model.Case
	if err := c.ShouldBindJSON(&caseModel); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	caseModel.ID = 0 // ensure DB auto-increment is used
	created, err := h.svc.Case.Create(&caseModel)
	if err != nil {
		logging.Logger.Warn("business error in CreateCase", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.Success(created))
}

func (h *Handler) UpdateCase(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid case id"))
		return
	}

	var req dto.UpdateCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	updated, err := h.svc.Case.Update(id, req)
	if err != nil {
		logging.Logger.Warn("business error in UpdateCase", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(updated))
}

func (h *Handler) DeleteCase(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid case id"))
		return
	}

	if err := h.svc.Case.Delete(id); err != nil {
		logging.Logger.Warn("business error in DeleteCase", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(nil))
}

// ListProjectCases returns cases belonging to a project.
func (h *Handler) ListProjectCases(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}

	cases, _, err := h.svc.Case.List(dto.CaseListRequest{
		ProjectID: &projectID,
	})
	if err != nil {
		logging.Logger.Error("failed in ListProjectCases", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(cases))
}

// CreateProjectCase creates a case bound to a project.
func (h *Handler) CreateProjectCase(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}

	var caseModel model.Case
	if err := c.ShouldBindJSON(&caseModel); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}
	caseModel.ID = 0 // ensure DB auto-increment is used
	caseModel.ProjectID = &projectID

	created, err := h.svc.Case.Create(&caseModel)
	if err != nil {
		logging.Logger.Warn("business error in CreateProjectCase", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.Success(created))
}

// UpdateProjectCase updates a case under a project.
func (h *Handler) UpdateProjectCase(c *gin.Context) {
	caseID, err := parseIDParam(c, "cid")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid case id"))
		return
	}

	var req dto.UpdateCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	updated, err := h.svc.Case.Update(caseID, req)
	if err != nil {
		logging.Logger.Warn("business error in UpdateProjectCase", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(updated))
}

// DeleteProjectCase hard-deletes a case under a project.
func (h *Handler) DeleteProjectCase(c *gin.Context) {
	caseID, err := parseIDParam(c, "cid")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid case id"))
		return
	}

	if err := h.svc.Case.Delete(caseID); err != nil {
		logging.Logger.Warn("business error in DeleteProjectCase", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(nil))
}
