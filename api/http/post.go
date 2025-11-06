package apihttp

import (
    httpx "majoo-case1-rest-api/internal/http"
    "majoo-case1-rest-api/internal/post"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

type postHandler struct{ uc *post.Usecase }

func RegisterPostRoutes(rg *gin.RouterGroup, uc *post.Usecase) {
    h := &postHandler{uc: uc}
    rg.GET("/posts", h.list)
    rg.GET("/posts/:id", h.get)
    rg.POST("/posts", h.create)
    rg.PUT("/posts/:id", h.update)
    rg.DELETE("/posts/:id", h.delete)
}

func (h *postHandler) list(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    posts, err := h.uc.List(page, limit)
    if err != nil { httpx.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch posts"); return }
    httpx.RespondWithSuccess(c, http.StatusOK, gin.H{"posts": posts, "page": page, "limit": limit})
}

func (h *postHandler) get(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil { httpx.RespondWithError(c, http.StatusBadRequest, "Invalid post ID"); return }
    p, err := h.uc.Get(id)
    if err != nil { httpx.RespondWithError(c, http.StatusNotFound, "Post not found"); return }
    httpx.RespondWithSuccess(c, http.StatusOK, p)
}

func (h *postHandler) create(c *gin.Context) {
    var req post.CreatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil { httpx.RespondWithError(c, http.StatusBadRequest, err.Error()); return }
    userID := c.MustGet("userID").(int)
    p, err := h.uc.Create(userID, req)
    if err != nil { httpx.RespondWithError(c, http.StatusInternalServerError, "Failed to create post"); return }
    httpx.RespondWithSuccess(c, http.StatusCreated, p)
}

func (h *postHandler) update(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil { httpx.RespondWithError(c, http.StatusBadRequest, "Invalid post ID"); return }
    var req post.UpdatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil { httpx.RespondWithError(c, http.StatusBadRequest, err.Error()); return }
    userID := c.MustGet("userID").(int)
    p, err := h.uc.Update(userID, id, req)
    if err != nil {
        if err == post.ErrForbidden { httpx.RespondWithError(c, http.StatusForbidden, "Forbidden"); return }
        httpx.RespondWithError(c, http.StatusInternalServerError, "Failed to update post"); return
    }
    httpx.RespondWithSuccess(c, http.StatusOK, p)
}

func (h *postHandler) delete(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil { httpx.RespondWithError(c, http.StatusBadRequest, "Invalid post ID"); return }
    userID := c.MustGet("userID").(int)
    if err := h.uc.Delete(userID, id); err != nil {
        if err == post.ErrForbidden { httpx.RespondWithError(c, http.StatusForbidden, "Forbidden"); return }
        httpx.RespondWithError(c, http.StatusInternalServerError, "Failed to delete post"); return
    }
    httpx.RespondWithMessage(c, http.StatusOK, "Post deleted successfully")
}


