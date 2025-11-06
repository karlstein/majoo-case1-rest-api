package httpx

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message,omitempty"`
}

func RespondWithError(c *gin.Context, statusCode int, message string) {
    c.JSON(statusCode, ErrorResponse{Error: http.StatusText(statusCode), Message: message})
}

func RespondWithSuccess(c *gin.Context, statusCode int, data interface{}) {
    c.JSON(statusCode, data)
}

func RespondWithMessage(c *gin.Context, statusCode int, message string) {
    c.JSON(statusCode, gin.H{"message": message})
}


