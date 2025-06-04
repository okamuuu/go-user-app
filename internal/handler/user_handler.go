package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/okamuuu/go-user-app/internal/domain"
	"github.com/okamuuu/go-user-app/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/users/:id", h.GetUsers)
	r.POST("/users", h.CreateUser)
	r.GET("/users/:id", h.GetUser)
	r.PUT("/users/:id", h.UpdateUser)
	r.DELETE("/users/:id", h.DeleteUser)
}

// @Summary ユーザー一覧取得
// @Description 登録されている全ユーザーを取得
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.User
// @Failure 401 {object} handler.ErrorResponse
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, err := h.service.GetUsers(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// CreateUser godoc
// @Summary      ユーザーの新規作成
// @Description  ユーザー情報を登録します。
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      domain.User  true  "ユーザー情報"
// @Success      201   {string}  string       "Created"
// @Failure      400   {object}  handler.ErrorResponse        "invalid request"
// @Failure      500   {object}  handler.ErrorResponse        "internal server error"
// @Router       /users [post]
// @Security     BearerAuth
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req domain.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request"})
		return
	}
	if err := h.service.CreateUser(&req); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

// GetUser godoc
// @Summary      ユーザーの取得
// @Description  指定されたIDのユーザー情報を取得します。
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ユーザーID"
// @Success      200  {object}  domain.User
// @Failure      400  {object}  handler.ErrorResponse  "invalid ID"
// @Failure      404  {object}  handler.ErrorResponse  "user not found"
// @Router       /users/{id} [get]
// @Security     BearerAuth
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid ID"})
		return
	}
	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Summary      ユーザー情報の更新
// @Description  指定されたIDのユーザー情報を更新します。
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ユーザーID"
// @Param        user body      domain.User true "更新するユーザー情報"
// @Success      200  {object}  domain.User
// @Failure      400  {object}  ErrorResponse  "invalid request or ID"
// @Failure      403  {object}  ErrorResponse "unauthorized"
// @Failure      404  {object}  ErrorResponse  "user not found"
// @Router       /users/{id} [put]
// @Security     BearerAuth
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// URL パラメータからユーザーIDを取得
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	// 認証情報から userID を取得（本人チェック用）
	authUserID, exists := c.Get("userID")
	if !exists || authUserID.(uint) != uint(id) {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "You can update only your own profile"})
		return
	}

	// リクエストボディをパース
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password,omitempty" binding:"omitempty,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ERROR] Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	// 更新用のUser構造体を作成
	user := &domain.User{
		ID:       uint(id),
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	log.Printf("[DEBUG] req: %+v", req)
	log.Printf("[DEBUG] user: %+v", user)

	// パスワードは任意更新のため、あればセット（暗号化はサービス層で）
	if req.Password != "" {
		user.Password = req.Password
	}

	if err := h.service.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update user"})
		return
	}

	c.Status(http.StatusOK)
}

// @Summary      ユーザーの削除
// @Description  指定されたIDのユーザーを削除します。
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ユーザーID"
// @Success      204  {string}  string  "No Content"
// @Failure      400  {object}  ErrorResponse  "invalid ID"
// @Failure      403  {object}  ErrorResponse  "unauthorized"
// @Failure      404  {object}  ErrorResponse  "user not found"
// @Router       /users/{id} [delete]
// @Security     BearerAuth
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid ID"})
		return
	}
	if err := h.service.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *UserHandler) Me(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not found in context"})
		return
	}

	user, err := h.service.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to get user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

type ErrorResponse struct {
	Error string `json:"error"`
}
