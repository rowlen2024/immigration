package handler

import (
	"errors"
	"net/http"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/logging"
	"mygo-immigration/backend/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *Handler) AdminListProjectCategories(c *gin.Context) {
	req, ok := bindProjectCategoryList(c)
	if !ok {
		return
	}
	items, total, err := h.svc.ProjectCategory.List(req)
	if err != nil {
		respondProjectCategoryError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.SuccessPaginated(items, req.Page, req.PerPage, total))
}

func (h *Handler) AdminListProjectCategoryOptions(c *gin.Context) {
	req, ok := bindProjectCategoryList(c)
	if !ok {
		return
	}
	items, total, err := h.svc.ProjectCategory.Options(req)
	if err != nil {
		respondProjectCategoryError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.SuccessPaginated(items, req.Page, req.PerPage, total))
}

func bindProjectCategoryList(c *gin.Context) (dto.ProjectCategoryListRequest, bool) {
	var req dto.ProjectCategoryListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "无效的类目查询参数"))
		return req, false
	}
	if req.Page == 0 {
		req.Page = defaultPage
	}
	if req.PerPage == 0 {
		req.PerPage = defaultPerPage
	}
	return req, true
}

func (h *Handler) CreateProjectCategory(c *gin.Context) {
	var req dto.SaveProjectCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "请检查类目名称、slug、SEO 文案和状态"))
		return
	}
	item, err := h.svc.ProjectCategory.Create(req)
	if err != nil {
		respondProjectCategoryError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.Success(item))
}

func (h *Handler) UpdateProjectCategory(c *gin.Context) {
	id, ok := parseProjectCategoryID(c)
	if !ok {
		return
	}
	var req dto.SaveProjectCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "请检查类目名称、slug、SEO 文案和状态"))
		return
	}
	item, err := h.svc.ProjectCategory.Update(id, req)
	if err != nil {
		respondProjectCategoryError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.Success(item))
}

func (h *Handler) DeleteProjectCategory(c *gin.Context) {
	id, ok := parseProjectCategoryID(c)
	if !ok {
		return
	}
	if err := h.svc.ProjectCategory.Delete(id); err != nil {
		respondProjectCategoryError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.Success(nil))
}

func (h *Handler) AdminListProjectCategoryProjects(c *gin.Context) {
	id, ok := parseProjectCategoryID(c)
	if !ok {
		return
	}
	req, ok := bindProjectCategoryList(c)
	if !ok {
		return
	}
	items, total, err := h.svc.ProjectCategory.Projects(id, req.Page, req.PerPage)
	if err != nil {
		respondProjectCategoryError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.SuccessPaginated(items, req.Page, req.PerPage, total))
}

func (h *Handler) UpdateProjectCategoryProjectOrder(c *gin.Context) {
	id, ok := parseProjectCategoryID(c)
	if !ok {
		return
	}
	projectID, err := parseIDParam(c, "projectId")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "无效的项目 ID"))
		return
	}
	var req dto.UpdateProjectCategoryProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "请输入有效的项目排序"))
		return
	}
	if err := h.svc.ProjectCategory.UpdateProjectOrder(id, projectID, *req.SortOrder); err != nil {
		respondProjectCategoryError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.Success(nil))
}

func (h *Handler) ListPublishedProjectCategories(c *gin.Context) {
	items, err := h.svc.ProjectCategory.PublicList()
	if err != nil {
		respondProjectCategoryError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.Success(items))
}

func (h *Handler) GetPublishedProjectCategory(c *gin.Context) {
	item, err := h.svc.ProjectCategory.PublicDetail(c.Param("slug"))
	if err != nil {
		respondProjectCategoryError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.Success(item))
}

func parseProjectCategoryID(c *gin.Context) (uint64, bool) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "无效的类目 ID"))
		return 0, false
	}
	return id, true
}

func respondProjectCategoryError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, dto.Error(404, "类目或类目项目不存在"))
	case errors.Is(err, service.ErrProjectCategoryName), errors.Is(err, service.ErrProjectCategorySlug),
		errors.Is(err, service.ErrProjectCategorySlugFrozen), errors.Is(err, service.ErrProjectCategorySlugTaken),
		errors.Is(err, service.ErrProjectCategoryInUse):
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
	default:
		logging.Logger.Error("project category request failed", "error", err)
		c.JSON(http.StatusInternalServerError, dto.Error(500, "internal server error"))
	}
}
