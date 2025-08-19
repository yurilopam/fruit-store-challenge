package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSON(c *gin.Context, code int, v any)       { c.JSON(code, v) }
func Error(c *gin.Context, code int, msg string) { c.JSON(code, gin.H{"error": msg}) }
func BadRequest(c *gin.Context, msg string)      { Error(c, http.StatusBadRequest, msg) }
func Unauthorized(c *gin.Context)                { Error(c, http.StatusUnauthorized, "unauthorized") }
func NotFound(c *gin.Context)                    { Error(c, http.StatusNotFound, "not found") }
