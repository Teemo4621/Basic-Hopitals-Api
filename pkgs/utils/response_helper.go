package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OkResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data":    data,
		"message": "success",
		"status":  "ok",
	})
}

func NotFoundResponse(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{
		"message": message,
		"status":  "error",
	})
}

func ErrorResponse(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"message": message,
		"status":  "error",
	})
}

func BadRequestResponse(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"message": message,
		"status":  "error",
	})
}

func UnauthorizedResponse(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"message": message,
		"status":  "error",
	})
}
