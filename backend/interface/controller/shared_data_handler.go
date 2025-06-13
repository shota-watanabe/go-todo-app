package controller

import (
	"net/http"
	"strconv"
	"go-todo-app/backend/usecase"

	"github.com/labstack/echo/v4"
)

// SharedDataHandler は共有データとプロジェクト関連のHTTPハンドラ
type SharedDataHandler struct {
	sharedDataUseCase usecase.SharedDataUseCase
	// TODO: Authミドルウェアなど、他の依存関係
}

// NewSharedDataHandler は SharedDataHandler を生成します
func NewSharedDataHandler(s usecase.SharedDataUseCase) *SharedDataHandler {
	return &SharedDataHandler{sharedDataUseCase: s}
}

// GetViewGroupCompanies は全てのグループ会社を取得するHTTPハンドラ
func (h *SharedDataHandler) GetViewGroupCompanies(c echo.Context) error {
	companies, err := h.sharedDataUseCase.GetViewGroupCompanies()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, companies)
}

// GetSharedProducts は共有商品を取得するHTTPハンドラ
func (h *SharedDataHandler) GetSharedProducts(c echo.Context) error {
	// TODO: 実際のオーガナイゼーションIDは認証情報から取得する
	// 今は仮でクエリパラメータから取得
	orgIDStr := c.QueryParam("organization_id")
	if orgIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "organization_id is required"})
	}
	orgID, err := strconv.ParseUint(orgIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid organization_id"})
	}

	products, err := h.sharedDataUseCase.GetSharedProducts(uint(orgID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, products)
}

// CreateProjectRequest はプロジェクト作成リクエストのDTO
type CreateProjectRequest struct {
	Name string `json:"name" validate:"required"`
	// TODO: organization_idは認証情報から取得するか、パスパラメータで渡す
}

// CreateProject はプロジェクトを作成するHTTPハンドラ
func (h *SharedDataHandler) CreateProject(c echo.Context) error {
	req := new(CreateProjectRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	// TODO: 認証情報からOrganizationIDを取得
	// 例: userID := c.Get("userID").(*jwt.Token).Claims.(*jwtCustomClaims).UserID
	// user, err := h.userUseCase.GetUserByID(userID)
	// if err != nil { return err }
	// orgID := user.OrganizationID
	// 仮のOrganizationID
	orgID := uint(1) // FIXME: 実際には認証されたユーザーのOrganizationIDを使用

	project, err := h.sharedDataUseCase.CreateProject(req.Name, orgID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, project)
}

// GetProjectProducts はプロジェクトの商品を取得するHTTPハンドラ
func (h *SharedDataHandler) GetProjectProducts(c echo.Context) error {
	projectIDStr := c.Param("projectID")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid project ID"})
	}

	// TODO: リクエストしているユーザーがこのプロジェクトにアクセスする権限があるかチェック
	// h.authUsecase.CheckProjectAccess(c.Get("userID"), projectID)

	products, err := h.sharedDataUseCase.GetProjectProducts(uint(projectID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, products)
}

// CopySharedProductToProjectRequest は共有商品をプロジェクトにコピーするリクエストのDTO
type CopySharedProductToProjectRequest struct {
	SharedProductID uint `json:"shared_product_id" validate:"required"`
}

// CopySharedProductToProject は共有商品をプロジェクトにコピーするHTTPハンドラ
func (h *SharedDataHandler) CopySharedProductToProject(c echo.Context) error {
	projectIDStr := c.Param("projectID")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid project ID"})
	}

	req := new(CopySharedProductToProjectRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// TODO: ユーザーがこのプロジェクトへのアクセス権限、およびコピー操作の権限があるかチェック
	// h.authUsecase.CheckProjectAccess(c.Get("userID"), projectID)

	projectProduct, err := h.sharedDataUseCase.CopySharedProductToProject(req.SharedProductID, uint(projectID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, projectProduct)
}

// PromoteProjectProductToShared はプロジェクトの商品を共有商品に格上げするHTTPハンドラ
// TODO: このAPIは管理者権限のみ許可されるべき
func (h *SharedDataHandler) PromoteProjectProductToShared(c echo.Context) error {
	projectProductIDStr := c.Param("projectProductID")
	projectProductID, err := strconv.ParseUint(projectProductIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid project product ID"})
	}

	// TODO: ユーザーが親オーガナイゼーションの管理者権限を持つかチェック
	// h.authUsecase.CheckAdminPermission(c.Get("userID"))

	sharedProduct, err := h.sharedDataUseCase.PromoteProjectProductToShared(uint(projectProductID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, sharedProduct)
}