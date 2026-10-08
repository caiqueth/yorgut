package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// notImplemented answers routes whose handler has not been written yet.
func notImplemented(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
