package handler

import (
	"net/http"
	"strconv"
	"strings"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/logging"
	"mygo-immigration/backend/internal/model"

	"github.com/gin-gonic/gin"
)

const defaultPage = 1
const defaultPerPage = 10

func (h *Handler) ListProjects(c *gin.Context) {
	if c.Query("all") == "true" {
		projects, err := h.svc.Project.ListAll("", "")
		if err != nil {
			logging.Logger.Error("failed in ListProjects", "error", err)
			c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
			return
		}
		c.JSON(http.StatusOK, dto.Success(projects))
		return
	}

	page, perPage := parsePagination(c)

	projects, total, err := h.svc.Project.List(page, perPage, "", "")
	if err != nil {
		logging.Logger.Error("failed in ListProjects", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessPaginated(projects, page, perPage, total))
}

func (h *Handler) GetProject(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, dto.Error(400, "slug is required"))
		return
	}

	project, err := h.svc.Project.GetBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Error(404, "project not found"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(project))
}

func (h *Handler) CompareProjects(c *gin.Context) {
	slugsParam := c.Query("slugs")
	if slugsParam == "" {
		c.JSON(http.StatusBadRequest, dto.Error(400, "slugs query param is required"))
		return
	}

	slugs := strings.Split(slugsParam, ",")

	var fields []string
	if fieldsParam := c.Query("fields"); fieldsParam != "" {
		fields = strings.Split(fieldsParam, ",")
	}

	result, err := h.svc.Project.CompareRows(slugs, fields)
	if err != nil {
		logging.Logger.Warn("business error in CompareProjects", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *Handler) AdminListProjects(c *gin.Context) {
	search := c.Query("search")
	status := c.Query("status")

	if c.Query("all") == "true" {
		projects, err := h.svc.Project.ListAll(search, status)
		if err != nil {
			logging.Logger.Error("failed in AdminListProjects", "error", err)
			c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
			return
		}
		c.JSON(http.StatusOK, dto.Success(projects))
		return
	}

	page, perPage := parsePagination(c)

	projects, total, err := h.svc.Project.AdminList(page, perPage, search, status)
	if err != nil {
		logging.Logger.Error("failed in AdminListProjects", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessPaginated(projects, page, perPage, total))
}

func (h *Handler) CreateProject(c *gin.Context) {
	var project model.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	project.ID = 0 // ensure DB auto-increment
	created, err := h.svc.Project.Create(&project)
	if err != nil {
		logging.Logger.Warn("business error in CreateProject", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.Success(created))
}

func (h *Handler) UpdateProject(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}

	var req dto.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	updated, err := h.svc.Project.Update(id, req)
	if err != nil {
		logging.Logger.Warn("business error in UpdateProject", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(updated))
}

func (h *Handler) DeleteProject(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}

	if err := h.svc.Project.Delete(id); err != nil {
		logging.Logger.Warn("business error in DeleteProject", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(nil))
}

func parsePagination(c *gin.Context) (int, int) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		return defaultPage, defaultPerPage
	}

	page := req.Page
	if page < 1 {
		page = defaultPage
	}

	perPage := req.PerPage
	if perPage < 1 {
		perPage = defaultPerPage
	}

	return page, perPage
}

type addNewsRequest struct {
	PageIDs []uint64 `json:"page_ids"`
}

// ListProjectNews returns news pages linked to a project.
func (h *Handler) ListProjectNews(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}

	news, err := h.svc.Project.ListNews(projectID)
	if err != nil {
		logging.Logger.Error("failed in ListProjectNews", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(news))
}

// AddProjectNews links news pages to a project.
func (h *Handler) AddProjectNews(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}

	var req addNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	if err := h.svc.Project.AddNews(projectID, req.PageIDs); err != nil {
		logging.Logger.Warn("business error in AddProjectNews", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(nil))
}

// RemoveProjectNews unlinks a news page from a project.
func (h *Handler) RemoveProjectNews(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}

	pageID, err := parseIDParam(c, "page_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid page id"))
		return
	}

	if err := h.svc.Project.RemoveNews(projectID, pageID); err != nil {
		logging.Logger.Warn("business error in RemoveProjectNews", "error", err)
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(nil))
}

func parseIDParam(c *gin.Context, param string) (uint64, error) {
	return strconv.ParseUint(c.Param(param), 10, 64)
}
