package apihttp

import (
	"majoo-case1-rest-api/internal/comment"
	httpx "majoo-case1-rest-api/internal/http"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type commentHandler struct{ uc *comment.Usecase }

func RegisterCommentRoutes(rg *gin.RouterGroup, uc *comment.Usecase) {
	h := &commentHandler{uc: uc}
	rg.GET("/posts/:id/comments", h.listByPost)
	rg.GET("/comments/:id", h.get)
	rg.POST("/posts/:id/comments", h.create)
	rg.PUT("/comments/:id", h.update)
	rg.DELETE("/comments/:id", h.delete)
}

func (h *commentHandler) listByPost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httpx.RespondWithError(c, http.StatusBadRequest, "Invalid post ID")
		return
	}
	items, err := h.uc.ListByPost(postID)
	if err != nil {
		httpx.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch comments")
		return
	}
	httpx.RespondWithSuccess(c, http.StatusOK, gin.H{"post_id": postID, "comments": items})
}

func (h *commentHandler) get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httpx.RespondWithError(c, http.StatusBadRequest, "Invalid comment ID")
		return
	}
	cm, err := h.uc.Get(id)
	if err != nil {
		httpx.RespondWithError(c, http.StatusNotFound, "Comment not found")
		return
	}
	httpx.RespondWithSuccess(c, http.StatusOK, cm)
}

func (h *commentHandler) create(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httpx.RespondWithError(c, http.StatusBadRequest, "Invalid post ID")
		return
	}
	var req comment.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	userID := c.MustGet("userID").(int)
	cm, err := h.uc.Create(postID, userID, req.Content)
	if err != nil {
		switch err {
		case comment.ErrNotFound:
			httpx.RespondWithError(c, http.StatusNotFound, "Post not found")
		default:
			httpx.RespondWithError(c, http.StatusInternalServerError, "Failed to create comment")
		}
		return
	}
	httpx.RespondWithSuccess(c, http.StatusCreated, cm)
}

func (h *commentHandler) update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httpx.RespondWithError(c, http.StatusBadRequest, "Invalid comment ID")
		return
	}
	var req comment.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	userID := c.MustGet("userID").(int)
	cm, err := h.uc.Update(userID, id, req.Content)
	if err != nil {
		if err == comment.ErrForbidden {
			httpx.RespondWithError(c, http.StatusForbidden, "Forbidden")
			return
		}
		httpx.RespondWithError(c, http.StatusInternalServerError, "Failed to update comment")
		return
	}
	httpx.RespondWithSuccess(c, http.StatusOK, cm)
}

func (h *commentHandler) delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httpx.RespondWithError(c, http.StatusBadRequest, "Invalid comment ID")
		return
	}
	userID := c.MustGet("userID").(int)
	if err := h.uc.Delete(userID, id); err != nil {
		if err == comment.ErrForbidden {
			httpx.RespondWithError(c, http.StatusForbidden, "Forbidden")
			return
		}
		httpx.RespondWithError(c, http.StatusInternalServerError, "Failed to delete comment")
		return
	}
	httpx.RespondWithMessage(c, http.StatusOK, "Comment deleted successfully")
}
