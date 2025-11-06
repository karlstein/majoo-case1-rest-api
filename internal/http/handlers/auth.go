package handlers

import (
    "majoo-case1-rest-api/config"
    httpx "majoo-case1-rest-api/internal/http"
    "majoo-case1-rest-api/internal/user"
    "net/http"

    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    usecase *user.Usecase
    cfg     config.Config
}

func NewAuthHandler(u *user.Usecase, cfg config.Config) *AuthHandler { return &AuthHandler{usecase: u, cfg: cfg} }

func (h *AuthHandler) Register(c *gin.Context) {
    var req user.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        httpx.RespondWithError(c, http.StatusBadRequest, err.Error())
        return
    }
    u, token, err := h.usecase.Register(req.Username, req.Email, req.Password)
    if err != nil {
        switch err {
        case user.ErrConflict:
            httpx.RespondWithError(c, http.StatusConflict, "User exists")
        default:
            httpx.RespondWithError(c, http.StatusInternalServerError, "Failed to register")
        }
        return
    }
    c.SetSameSite(http.SameSiteLaxMode)
    isSecure := c.Request.TLS != nil
    c.SetCookie("token", token, 24*60*60, "/", "", isSecure, true)
    httpx.RespondWithSuccess(c, http.StatusCreated, user.LoginResponse{Token: token, User: u})
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req user.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        httpx.RespondWithError(c, http.StatusBadRequest, err.Error())
        return
    }
    u, token, err := h.usecase.Login(req.Email, req.Password)
    if err != nil {
        switch err {
        case user.ErrUnauthorized:
            httpx.RespondWithError(c, http.StatusUnauthorized, "Invalid credentials")
        default:
            httpx.RespondWithError(c, http.StatusInternalServerError, "Login failed")
        }
        return
    }
    c.SetSameSite(http.SameSiteLaxMode)
    isSecure := c.Request.TLS != nil
    c.SetCookie("token", token, 24*60*60, "/", "", isSecure, true)
    httpx.RespondWithSuccess(c, http.StatusOK, user.LoginResponse{Token: token, User: u})
}


